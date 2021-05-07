package autoupdater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type autoUpdater struct {
	version   string
	appName   string
	apiURL    string
	osPostFix string
}

func NewAutoUpdater(version string) autoUpdater {
	var osPostFix string
	switch runtime.GOOS {
	case "darwin":
		osPostFix = "_osx"
	case "linux":
		osPostFix = "_linux"
	case "windows":
		osPostFix = "_windows.exe"
	}

	return autoUpdater{
		version:   version,
		appName:   "repo_info_extractor" + osPostFix,
		apiURL:    "https://api.github.com/repos/codersrank-org/repo_info_extractor/releases/latest",
		osPostFix: osPostFix,
	}
}

// CheckUpdates checks github to see if there is a new version and if there is one, downloads it.
func (au autoUpdater) CheckUpdates() {
	fmt.Println("Checking for new versions. Current version: " + au.version)
	release, err := au.getRelease()
	if err != nil {
		fmt.Printf("Couldn't get latest release from Github, skipping update. Error: %s\n", err.Error())
		return
	}
	latestVersion, err := au.getLatestVersion(release)
	if err != nil {
		fmt.Printf("Couldn't find the latest version, skipping update. Error: %s\n", err.Error())
		return
	}
	if au.shouldUpdate(latestVersion) {
		fmt.Printf("Found new version %s, updating...\n", latestVersion)
		err := au.update(release)
		if err != nil {
			fmt.Printf("Couldn't download latest release. Error: %s\n", err.Error())
		} else {
			fmt.Println("New version downloaded. Please run the program again.")
			os.Exit(0)
		}
	} else {
		fmt.Printf("You already have latest version, skipping update\n")
	}
}

func (au autoUpdater) update(r *release) error {
	for _, asset := range r.Assets {
		// Found the correct binary
		if strings.Contains(asset.Name, au.osPostFix) {
			fmt.Printf("Downloading %s\n", asset.BrowserDownloadURL)
			return au.download(asset.BrowserDownloadURL)
		}
	}

	return nil
}

func (au autoUpdater) download(downloadURL string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	appPath, err := os.Getwd()
	if err != nil {
		return err
	}
	oldName := filepath.Join(appPath, au.appName)
	newName := filepath.Join(appPath, au.appName) + "_old"
	err = os.Rename(oldName, newName)
	if err != nil {
		fmt.Printf("Couldn't rename file from %s to %s\n", oldName, newName)
	}
	filePath := appPath + "/" + au.appName

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("New binary saved to %s\n", filePath)
	} else {
		chmodErr := os.Chmod(filePath, 0755)
		if chmodErr != nil {
			fmt.Printf("Couldn't set execute permissions for %s\n", filePath)
		}
	}
	return err
}

func (au autoUpdater) getRelease() (*release, error) {
	request, err := http.NewRequest(http.MethodGet, au.apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var r *release
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (au autoUpdater) getLatestVersion(r *release) (string, error) {

	// "v" is not part of semantic versioning
	r.Name = strings.TrimLeft(r.Name, "v")

	// Regex for finding Major, Minor and Patch versions
	// Taken from here: https://semver.org/
	regex := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	matches := regex.FindAllString(r.Name, -1)
	if len(matches) == 0 {
		return "", errors.New("Couldn't parse current version")
	}

	// Split as major, minor and patch
	matches = strings.Split(matches[0], ".")
	if len(matches) != 3 {
		return "", errors.New("Couldn't parse current version")
	}

	var err error
	major, err := strconv.Atoi(matches[0])
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(matches[2])

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), err
}

func (au autoUpdater) shouldUpdate(version string) bool {
	return version > au.version
}

type release struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		NodeID   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}
