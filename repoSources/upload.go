package repoSource

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codersrank-org/repo_info_extractor/entities"
	"github.com/codersrank-org/repo_info_extractor/ui"
	"github.com/pkg/browser"
)

// ArtifactUploader uploads and merge results with codersrank
type ArtifactUploader interface {
	UploadRepos(repos []*entities.Repository)
}

type artifactUploader struct {
	// UploadRepoURL the API end pont where the artifacts will be uploaded
	UploadRepoURL   string
	UploadResultURL string
	ProcessURL      string
	// OutputPath contains the artifacts that need to be uploaded
	OutputPath string
}

// NewArtifactUploader constructor
func NewArtifactUploader(outputPath string) ArtifactUploader {
	return &artifactUploader{
		//UploadRepoURL:   "https://grpcgateway.codersrank.io/candidate/privaterepo/Upload",
		UploadRepoURL: "http://localhost:9900/candidate/privaterepo/Upload",
		//UploadResultURL: "https://grpcgateway.codersrank.io/multi/repo/results",
		UploadResultURL: "http://localhost:9900/multi/repo/results",
		//ProcessURL:      "https://profile.codersrank.io/repo?multiToken=",
		ProcessURL: "http://localhost:8080/repo?multiToken=",
		OutputPath: outputPath,
	}
}

func (c *artifactUploader) UploadRepos(repos []*entities.Repository) {
	uploadResults := make(map[string]string)
	done := 1
	for _, repo := range repos {
		fmt.Printf("Uploading %s results (%d/%d)\n", repo.FullName, done, len(repos))
		uploadToken, err := c.uploadRepo(repo.GetSafeFullName())
		if err != nil {
			fmt.Printf("Couldn't upload, error: %s", err.Error())
			continue
		}
		uploadResults[repo.Name] = uploadToken
		done++
	}
	resultToken := c.uploadResults(uploadResults)
	c.processResults(resultToken)
}

func (c *artifactUploader) uploadRepo(repoName string) (string, error) {

	// Read file
	filename := fmt.Sprintf("%s/%s_v2.json.zip", c.getSaveResultPath(), repoName)
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Add file as multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}
	io.Copy(part, file)
	writer.Close()

	// Create and make the request
	request, err := http.NewRequest("POST", c.UploadRepoURL, body)
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Server returned non 200 response")
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Get response and return resulting token
	var result CRUploadResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *artifactUploader) uploadResults(results map[string]string) string {

	multiUpload := MultiUpload{}
	multiUpload.Results = make([]CRUploadResultWithRepoName, len(results))

	i := 0
	for reponame, token := range results {
		multiUpload.Results[i] = CRUploadResultWithRepoName{
			Token:    token,
			Reponame: reponame,
		}
		i++
	}

	b, err := json.Marshal(multiUpload)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", c.UploadResultURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result CRUploadResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result.Token

}

func (c *artifactUploader) processResults(resultToken string) {
	browserURL := c.ProcessURL + resultToken
	ok := ui.Confirm(fmt.Sprintf("You are being navigated to '%s'. Do you wish to proceed?", browserURL))
	if ok {
		browser.OpenURL(browserURL)
	} else {
		fmt.Println("Finished")
	}
}

func (c *artifactUploader) getSaveResultPath() string {
	resultPath := c.OutputPath
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		os.Mkdir(resultPath, 0700)
	}
	return resultPath
}

// CRUploadResult is the result of single repo upload
type CRUploadResult struct {
	Token string `json:"token"`
}

// MultiUpload is the request body
type MultiUpload struct {
	Results []CRUploadResultWithRepoName `json:"results"`
}

// CRUploadResultWithRepoName token-reponame pair
type CRUploadResultWithRepoName struct {
	Token    string `json:"token"`
	Reponame string `json:"reponame"`
}
