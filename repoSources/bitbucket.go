package repoSource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/codersrank-org/repo_info_extractor/entities"
	"github.com/pkg/errors"
)

// BitbucketProvider Bitbucket provider used for handling Bitbucket API operations
type BitbucketProvider struct {
	Scheme        string
	BaseURL       string
	Path          string
	Username      string
	Password      string
	Visibility    string
	GitExecutable string
	dirsToCleanUp []string
}

// NewBitbucketProvider constructor
func NewBitbucketProvider(username, password, repoVisibility, gitPath string) RepoSource {
	return &BitbucketProvider{
		Scheme:        "https",
		BaseURL:       "api.bitbucket.org",
		Path:          "2.0/repositories",
		Username:      username,
		Password:      password,
		Visibility:    repoVisibility,
		GitExecutable: gitPath,
		dirsToCleanUp: make([]string, 0),
	}
}

func (p *BitbucketProvider) CleanUp() {
	for _, path := range p.dirsToCleanUp {
		os.RemoveAll(path)
	}
}

// Clone repository from given url to given path
func (p *BitbucketProvider) cloneRepository(url string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "clone_dir_")
	if err != nil {
		return "", fmt.Errorf("couldn't create temp dir. Error: %s", err.Error())
	}

	cmd := exec.Command(p.GitExecutable,
		"clone",
		url,
		tmpDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "could not clone repo. Error: "+string(output))
	}

	return tmpDir, nil
}

func (p *BitbucketProvider) Clone(repo *entities.Repository) (string, error) {
	path, err := p.cloneRepository(repo.CloneURL)

	return path, err
}

// GetRepos returns list of repositories with given token and visibility from provider
func (p *BitbucketProvider) GetRepos() []*entities.Repository {
	requestURL := url.URL{
		Scheme: p.Scheme,
		Host:   p.BaseURL,
		Path:   p.Path,
	}

	query := requestURL.Query()
	// role is required otherwise we will get all bitbucket repos.
	query.Set("role", "contributor")

	if p.Visibility == "public" {
		// By default Bitbucket API returns all repositories
		query.Set("q", "is_private = false")
	}

	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		log.Fatal("Couldn't create request. Error: " + err.Error())
	}

	request.SetBasicAuth(p.Username, p.Password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Couldn't make HTTP request to the BitBucket API. Error: " + err.Error())
	}
	if response.StatusCode != 200 {
		log.Fatalf("Couldn't make HTTP request to the BitBucket API. Status: %s (Code: %d)", response.Status, response.StatusCode)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Couldn't read response body. Error: " + err.Error())
	}

	var bitbucketRepos *bitbucketRepository
	err = json.Unmarshal(body, &bitbucketRepos)
	if err != nil {
		log.Fatal("Couldn't parse JSON. Error: " + err.Error() + "\nJSON content: " + string(body))
	}

	repos := make([]*entities.Repository, len(bitbucketRepos.Values))
	for index, repo := range bitbucketRepos.Values {
		cloneUrl := ""
		for _, u := range repo.Links.Clone {
			if u.Name == "https" {
				cloneUrl = u.Href
			}
		}
		repos[index] = &entities.Repository{
			ID:       repo.UUID,
			FullName: repo.FullName,
			Name:     repo.Name,
			CloneURL: cloneUrl,
		}
	}

	return repos
}

// bitbucketRepository response from Bitbucket API
type bitbucketRepository struct {
	Values []struct {
		Links struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
		} `json:"links"`
		UUID     string `json:"uuid"`
		FullName string `json:"full_name"`
		Name     string `json:"name"`
	} `json:"values"`
}
