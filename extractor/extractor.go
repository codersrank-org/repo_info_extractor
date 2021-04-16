package extractor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/codersrank-org/repo_info_extractor/commit"
	"github.com/codersrank-org/repo_info_extractor/emailsimilarity"
	"github.com/codersrank-org/repo_info_extractor/languagedetection"
	"github.com/codersrank-org/repo_info_extractor/librarydetection"
	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
	"github.com/codersrank-org/repo_info_extractor/obfuscation"
	"github.com/codersrank-org/repo_info_extractor/ui"
	"github.com/mholt/archiver"
)

// TODO handle async errors correctly

// RepoExtractor is responsible for all parts of repo extraction process
// Including cloning the repo, processing the commits and uploading the results
type RepoExtractor struct {
	RepoPath            string
	OutputPath          string
	GitPath             string
	Headless            bool
	Obfuscate           bool
	ShowProgressBar     bool // If it is false there is no progress bar.
	UserEmails          []string
	OverwrittenRepoName string // If set this will be used instead of the
	Seed                []string
	repo                *repo
	userCommits         []*commit.Commit // Commits which are belong to user (from selected emails)
}

// Extract a single repo in the path
func (r *RepoExtractor) Extract() error {

	r.initGit()

	err := r.initRepo()
	if err != nil {
		fmt.Println("Cannot init repo_info_extractor. Error: ", err.Error())
		return err
	}

	// For library detection
	r.initAnalyzers()

	err = r.analyseCommits()
	if err != nil {
		return err
	}

	err = r.analyseLibraries()
	if err != nil {
		return err
	}

	if r.Obfuscate {
		r.obfuscate()
	}

	err = r.export()
	if err != nil {
		return err
	}

	// Only when user running this script locally
	if !r.Headless {
		err = r.upload()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoExtractor) initGit() {
	// Git path already provided by user
	if r.GitPath != "" {
		return
	}

	cmd := exec.Command("where", "git")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Try default git path
		r.GitPath = "/usr/bin/git"
		return
	}

	gitPath := string(out)
	gitPath = strings.TrimRight(gitPath, "\r\n")
	gitPath = strings.TrimRight(gitPath, "\n")

	r.GitPath = gitPath
}

// Creates Repo struct
func (r *RepoExtractor) initRepo() error {
	fmt.Println("Initializing repository")

	cmd := exec.Command(r.GitPath,
		"config",
		"--get",
		"remote.origin.url",
	)
	cmd.Dir = r.RepoPath

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Cannot get remote.origin.url. Use directory path to get repo name.")
	}

	repoName := ""
	remoteOrigin := string(out)
	remoteOrigin = strings.TrimRight(remoteOrigin, "\r\n")
	remoteOrigin = strings.TrimRight(remoteOrigin, "\n")

	repoName = r.GetRepoName(remoteOrigin)

	r.repo = &repo{
		RepoName:        repoName,
		Emails:          []string{},
		SuggestedEmails: []string{}, // TODO implement
	}
	return nil
}

// GetRepoName gets the repo name in the following format:
// in case of headless: "owner_name/repo_name"
// in case of interactive mode: "repo_name"
func (r *RepoExtractor) GetRepoName(remoteOrigin string) string {
	// If remoteOrigin is empty fall back to the repos path. It can happen in interactive mode
	if remoteOrigin == "" {
		parts := strings.Split(r.RepoPath, "/")
		return parts[len(parts)-1]
	}
	repoName := ""
	remoteOrigin = strings.TrimSuffix(remoteOrigin, ".git")
	if strings.Contains(remoteOrigin, "http") {
		// Cloned using http
		parts := strings.Split(remoteOrigin, "/")
		if r.Headless {
			repoName = parts[len(parts)-2] + "/" + parts[len(parts)-1]
		} else {
			// If it's a private repo, we only need last part of the name
			repoName = parts[len(parts)-1]
		}
	} else {
		// Cloned using ssh
		parts := strings.Split(remoteOrigin, ":")
		repoName = parts[len(parts)-1]
		parts = strings.Split(repoName, "/")
		if r.Headless {
			repoName = parts[len(parts)-2] + "/" + parts[len(parts)-1]
		} else {
			repoName = parts[len(parts)-1]
		}
	}

	return repoName
}

func (r *RepoExtractor) initAnalyzers() {
	librarydetection.AddAnalyzer("Go", languages.NewGoAnalyzer())
	librarydetection.AddAnalyzer("C", languages.NewCAnalyzer())
	librarydetection.AddAnalyzer("C++", languages.NewCppAnalyzer())
	librarydetection.AddAnalyzer("C#", languages.NewCSharpAnalyzer())
	librarydetection.AddAnalyzer("Java", languages.NewJavaAnalyzer())
	librarydetection.AddAnalyzer("JavaScript", languages.NewJavaScriptAnalyzer())
	librarydetection.AddAnalyzer("Kotlin", languages.NewKotlinAnalyzer())
	librarydetection.AddAnalyzer("TypeScript", languages.NewTypeScriptAnalyzer())
	librarydetection.AddAnalyzer("Perl", languages.NewPerlAnalyzer())
	librarydetection.AddAnalyzer("PHP", languages.NewPHPAnalyzer())
	librarydetection.AddAnalyzer("Python", languages.NewPythonScriptAnalyzer())
	librarydetection.AddAnalyzer("Swift", languages.NewSwiftAnalyzer())
}

// Creates commits
func (r *RepoExtractor) analyseCommits() error {
	fmt.Println("Analysing commits")

	var commits []*commit.Commit
	userCommits := make([]*commit.Commit, 0, len(commits))
	commits, err := r.getCommits()
	if len(commits) == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	allEmails := getAllEmails(commits)
	selectedEmails := make(map[string]bool)

	// If seed is provided use it in headless mode
	if len(r.Seed) > 0 && r.Headless {
		similarEmails := emailsimilarity.FindSimilarEmails(r.Seed, allEmails)
		similarEmailsWithoutNames, selectedSimilarEmailsMap := getEmailsWithoutNames(similarEmails)
		r.repo.SuggestedEmails = similarEmailsWithoutNames
		for mail := range selectedSimilarEmailsMap {
			selectedEmails[mail] = true
		}
	}

	if len(r.UserEmails) == 0 && !r.Headless {
		selectedEmailsWithNames := ui.SelectEmail(allEmails)
		emails, emailsMap := getEmailsWithoutNames(selectedEmailsWithNames)
		r.repo.Emails = append(r.repo.Emails, emails...)
		for mail := range emailsMap {
			selectedEmails[mail] = true
		}
	} else {
		r.repo.Emails = append(r.repo.Emails, r.UserEmails...)
		for _, email := range r.UserEmails {
			selectedEmails[email] = true
		}
	}

	// Only consider commits for user
	for _, v := range commits {
		if _, ok := selectedEmails[v.AuthorEmail]; ok {
			userCommits = append(userCommits, v)
		}
	}

	r.userCommits = userCommits
	return nil
}

func (r *RepoExtractor) getCommits() ([]*commit.Commit, error) {
	jobs := make(chan *req)
	results := make(chan []*commit.Commit)
	noMoreChan := make(chan bool)
	for w := 0; w < runtime.NumCPU(); w++ {
		go func() {
			err := r.commitWorker(w, jobs, results, noMoreChan)
			if err != nil {
				fmt.Println("Error during getting commits. Error: " + err.Error())
			}
		}()
	}

	// launch initial jobs
	lastOffset := 0
	step := 1000
	for x := 0; x < runtime.NumCPU(); x++ {
		jobs <- &req{
			Limit:  step,
			Offset: x * step,
		}
		lastOffset = step * x
	}

	var commits []*commit.Commit
	workersReturnedNoMore := 0

	var pb ui.ProgressBar
	numberOfCommits := r.getNumberOfCommits()
	if r.ShowProgressBar && numberOfCommits > 0 {
		pb = ui.NewProgressBar(numberOfCommits)
	} else {
		pb = ui.NilProgressBar()
	}

	func() {
		for {
			select {
			case res := <-results:
				lastOffset += step
				jobs <- &req{
					Limit:  step,
					Offset: lastOffset,
				}
				commits = append(commits, res...)
				pb.SetCurrent(len(commits))
			case <-noMoreChan:
				workersReturnedNoMore++
				if workersReturnedNoMore == runtime.NumCPU() {
					close(jobs)
					return
				}
			}
		}
	}()
	pb.Finish()

	return commits, nil
}

func getAllEmails(commits []*commit.Commit) []string {
	allEmails := make([]string, 0, len(commits))
	emails := make(map[string]bool) // To prevent duplicates
	for _, v := range commits {
		if _, ok := emails[v.AuthorEmail]; !ok {
			emails[v.AuthorEmail] = true
			allEmails = append(allEmails, fmt.Sprintf("%s -> %s", v.AuthorName, v.AuthorEmail))
		}
	}
	return allEmails
}

func getEmailsWithoutNames(emails []string) ([]string, map[string]bool) {
	emailsWithoutNames := make(map[string]bool, len(emails))
	emailsWithoutNamesArray := make([]string, len(emails))
	for i, selectedEmail := range emails {
		fields := strings.Split(selectedEmail, " -> ")
		// TODO handle authorName being empty
		if len(fields) > 0 {
			emailsWithoutNames[fields[1]] = true
			emailsWithoutNamesArray[i] = fields[1]
		}
	}
	return emailsWithoutNamesArray, emailsWithoutNames
}

func (r *RepoExtractor) getNumberOfCommits() int {
	cmd := exec.Command(r.GitPath,
		"--no-pager",
		"log",
		"--all",
		"--no-merges",
		"--pretty=oneline",
	)
	cmd.Dir = r.RepoPath
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Cannot get number of commits. Cannot show progress bar. Error: " + err.Error())
		return 0
	}
	return strings.Count(string(stdout), "\n")
}

// commitWorker get commits from git
func (r *RepoExtractor) commitWorker(w int, jobs <-chan *req, results chan<- []*commit.Commit, noMoreChan chan<- bool) error {
	for v := range jobs {
		var commits []*commit.Commit

		cmd := exec.Command(r.GitPath,
			"log",
			"--numstat",
			"--all",
			fmt.Sprintf("--skip=%d", v.Offset),
			fmt.Sprintf("--max-count=%d", v.Limit),
			"--pretty=format:|||BEGIN|||%H|||SEP|||%an|||SEP|||%ae|||SEP|||%ad",
			"--no-merges",
		)
		cmd.Dir = r.RepoPath
		stdout, err := cmd.StdoutPipe()
		if nil != err {
			fmt.Println("Cannot create pipe.")
			return err
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("Error during execution of Git command.")
			return err
		}

		// parse the output into stats
		scanner := bufio.NewScanner(stdout)
		currentLine := 0
		var currectCommit *commit.Commit
		for scanner.Scan() {
			m := scanner.Text()
			currentLine++
			if m == "" {
				continue
			}
			if strings.HasPrefix(m, "|||BEGIN|||") {
				// we reached a new commit
				// save the existing
				if currectCommit != nil {
					commits = append(commits, currectCommit)
				}

				// and add new one commit
				m = strings.Replace(m, "|||BEGIN|||", "", 1)
				bits := strings.Split(m, "|||SEP|||")
				changedFiles := []*commit.ChangedFile{}
				dateStr := ""
				t, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", bits[3])
				if err == nil {
					dateStr = t.Format("2006-01-02 15:04:05 -0700")
				} else {
					fmt.Println("Cannot convert date. Expected date format: Mon Jan 2 15:04:05 2006 -0700. Got: " + bits[3])
				}
				currectCommit = &commit.Commit{
					Hash:         bits[0],
					AuthorName:   bits[1],
					AuthorEmail:  bits[2],
					Date:         dateStr,
					ChangedFiles: changedFiles,
				}
				continue
			}

			bits := strings.Fields(m)

			insertionsString := bits[0]
			if insertionsString == "-" {
				insertionsString = "0"
			}
			insertions, err := strconv.Atoi(insertionsString)
			if err != nil {
				fmt.Println("Cannot convert the following into integer: " + insertionsString)
				return err
			}

			deletionsString := bits[1]
			if deletionsString == "-" {
				deletionsString = "0"
			}
			deletions, err := strconv.Atoi(deletionsString)
			if err != nil {
				fmt.Println("Cannot convert the following into integer: " + deletionsString)
				return err
			}

			fileName := bits[2]
			// it is a rename, skip
			if strings.Contains("=>", fileName) {
				continue
			}

			changedFile := &commit.ChangedFile{
				Path:       bits[2],
				Insertions: insertions,
				Deletions:  deletions,
			}

			if currectCommit == nil {
				// TODO maybe skip? does this break anything?
				return errors.New("did not expect current commit to be null")
			}

			if currectCommit.ChangedFiles == nil {
				// TODO maybe skip? does this break anything?
				return errors.New("did not expect current commit changed files to be null")
			}

			currectCommit.ChangedFiles = append(currectCommit.ChangedFiles, changedFile)
		}

		// last commit will not get appended otherwise
		// because scanner is not returning anything
		if currectCommit != nil {
			commits = append(commits, currectCommit)
		}

		if len(commits) == 0 {
			noMoreChan <- true
			return nil
		}
		results <- commits
	}
	return nil
}

// TODO This is not ready yet (can't find libraries based on language -> look at libraryWorker)
func (r *RepoExtractor) analyseLibraries() error {
	fmt.Println("Analysing libraries")

	jobs := make(chan *commit.Commit, len(r.userCommits))
	results := make(chan bool, len(r.userCommits))
	// Analyse libraries for every commit
	for w := 1; w <= runtime.NumCPU(); w++ {
		go r.libraryWorker(jobs, results)
	}
	for _, v := range r.userCommits {
		jobs <- v
	}
	close(jobs)
	var pb ui.ProgressBar
	if r.ShowProgressBar {
		pb = ui.NewProgressBar(len(r.userCommits))
	} else {
		pb = ui.NilProgressBar()
	}
	for a := 1; a <= len(r.userCommits); a++ {
		<-results
		pb.Inc()
	}
	pb.Finish()
	return nil
}

func (r *RepoExtractor) libraryWorker(commits <-chan *commit.Commit, results chan<- bool) error {
	languageAnalyzer := languagedetection.NewLanguageAnalyzer()
	for commit := range commits {
		libraries := map[string][]string{}
		for n, fileChange := range commit.ChangedFiles {

			lang := ""
			fileContents := make([]byte, 0)

			extension := filepath.Ext(fileChange.Path)
			if extension == "" {
				continue
			}
			// remove the trailing dot
			extension = extension[1:]

			cmd := exec.Command(r.GitPath,
				"--no-pager",
				"show",
				fmt.Sprintf("%s:%s", commit.Hash, fileChange.Path),
			)
			cmd.Dir = r.RepoPath
			var err error
			fileContents, err = cmd.CombinedOutput()
			if err != nil {
				searchString1 := fmt.Sprintf("Path '%s' does not exist in '%s'", fileChange.Path, commit.Hash)
				searchString2 := fmt.Sprintf("Path '%s' exists on disk, but not in '%s'", fileChange.Path, commit.Hash)
				// Ignore case is needed because on windows error message starts with lowercase letter, in other systems it starts with uppercase letter
				stringSearcher := search.New(language.English, search.IgnoreCase)
				// means the file was deleted, skip
				start, end := stringSearcher.IndexString(string(fileContents), searchString1)
				if start != -1 && end != -1 {
					continue
				}
				start, end = stringSearcher.IndexString(string(fileContents), searchString2)
				if start != -1 && end != -1 {
					continue
				}
				return err
			}
			if languageAnalyzer.ShouldUseFile(extension) {
				lang = languageAnalyzer.DetectLanguageFromFile(fileChange.Path, fileContents)
			} else {
				lang = languageAnalyzer.DetectLanguageFromExtension(extension)
			}

			// We don't know extension, nothing to do
			if lang == "" {
				continue
			}

			commit.ChangedFiles[n].Language = lang
			analyzer, err := librarydetection.GetAnalyzer(lang)
			if err != nil {
				continue
			}

			fileLibraries, err := analyzer.ExtractLibraries(string(fileContents))
			if err != nil {
				fmt.Printf("error extracting libraries for %s: %s \n", lang, err.Error())
			}
			if libraries[lang] == nil {
				libraries[lang] = make([]string, 0)
			}
			libraries[lang] = append(libraries[lang], fileLibraries...)
		}
		commit.Libraries = libraries
		results <- true
	}
	return nil
}

// Obfuscate the result
func (r *RepoExtractor) obfuscate() {
	for _, commit := range r.userCommits {
		commit = obfuscation.Obfuscate(commit)
	}
}

// Writes result to the file
func (r *RepoExtractor) export() error {
	fmt.Println("Creating output file")

	repoDataPath := r.OutputPath + "_v2.json"
	zipPath := r.OutputPath + "_v2.json.zip"
	// Remove old files
	os.Remove(repoDataPath)
	os.Remove(zipPath)

	// Only do this when not using default value
	if r.OutputPath != "./repo_data_v2" {
		err := os.MkdirAll(r.OutputPath, 0755)
		if err != nil {
			log.Println("Cannot create directory. Error:", err.Error())
		}
	}

	file, err := os.Create(repoDataPath)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	r.repo.RepoName = r.OverwrittenRepoName
	repoMetaData, err := json.Marshal(r.repo)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(repoMetaData))

	for _, commit := range r.userCommits {
		commitData, err := json.Marshal(commit)
		if err != nil {
			fmt.Printf("Couldn't write commit to file. CommitHash: %s Error: %s", commit.Hash, err.Error())
			continue
		}
		fmt.Fprintln(w, string(commitData))
	}
	w.Flush() // important
	file.Close()

	err = archiver.Archive([]string{repoDataPath}, zipPath)
	if err != nil {
		return err
	}

	// We don't need this because we already have zip file
	os.Remove(repoDataPath)
	return nil
}

// This is for repo_info_extractor used locally and for user to
// upload his/her results automatically to the codersrank
func (r *RepoExtractor) upload() error {
	fmt.Println("Uploading result to CodersRank")
	url, err := Upload(r.OutputPath+"_v2.json.zip", r.repo.RepoName)
	if err != nil {
		return err
	}
	fmt.Printf("Go to this link in the browser => %s", url)
	return nil
}

type repo struct {
	RepoName        string   `json:"repo"`
	Emails          []string `json:"emails"`
	SuggestedEmails []string `json:"suggestedEmails"`
}

type req struct {
	Limit  int
	Offset int
}
