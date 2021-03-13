package extractor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/codersrank-org/repo_info_extractor/librarydetection"
	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
	"github.com/codersrank-org/repo_info_extractor/obfuscation"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mholt/archiver"
	"github.com/src-d/enry/v2"
)

// TODO handle async errors correctly

// RepoExtractor is responsible for all parts of repo extraction process
// Including cloning the repo, processing the commits and uploading the results
type RepoExtractor struct {
	RepoPath    string
	OutputPath  string
	GitPath     string
	Headless    bool
	Obfuscate   bool
	UserEmails  []string
	Seed        []string
	repo        *repo
	userCommits []*commit.Commit // Commits which are belong to user (from selected emails)
}

// Extract a single repo in the path
func (r *RepoExtractor) Extract() error {

	r.initGit()

	err := r.initRepo()
	if err != nil {
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
		return err
	}

	repoName := ""
	remoteOrigin := string(out)
	remoteOrigin = strings.TrimRight(remoteOrigin, "\r\n")
	remoteOrigin = strings.TrimRight(remoteOrigin, "\n")

	repoName = r.GetRepoName(remoteOrigin)

	r.repo = &repo{
		RepoName:         repoName,
		Emails:           []string{},
		SuggestedEmails:  []string{}, // TODO implement
		PrimaryRemoteURL: remoteOrigin,
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
		repoName = parts[1]
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
	librarydetection.AddAnalyzer("C#", languages.NewCSharpAnalyzer())
	librarydetection.AddAnalyzer("JavaScript", languages.NewJavaScriptAnalyzer())
	librarydetection.AddAnalyzer("TypeScript", languages.NewTypeScriptAnalyzer())
	librarydetection.AddAnalyzer("Python", languages.NewPythonScriptAnalyzer())
}

// Creates commits
func (r *RepoExtractor) analyseCommits() error {
	fmt.Println("Analysing commits")

	var commits []*commit.Commit
	commits, err := r.getCommits()
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

	if len(r.UserEmails) == 0 {
		// Ask user for emails
		// TODO sort by alphabetical order (or frequency?)
		selectedEmailsWithNames := []string{}
		prompt := &survey.MultiSelect{
			Message:  "Please choose your emails:",
			Options:  allEmails,
			PageSize: 50,
		}
		survey.AskOne(prompt, &selectedEmailsWithNames)

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
	userCommits := make([]*commit.Commit, 0, len(commits))
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
		go r.commitWorker(w, jobs, results, noMoreChan)
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
			case <-noMoreChan:
				workersReturnedNoMore++
				if workersReturnedNoMore == runtime.NumCPU() {
					close(jobs)
					return
				}
			}
		}
	}()

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

// commitWorker get commits from git
func (r *RepoExtractor) commitWorker(w int, jobs <-chan *req, results chan<- []*commit.Commit, noMoreChan chan<- bool) error {
	for v := range jobs {
		var commits []*commit.Commit

		cmd := exec.Command(r.GitPath,
			"log",
			"--numstat",
			fmt.Sprintf("--skip=%d", v.Offset),
			fmt.Sprintf("--max-count=%d", v.Limit),
			"--pretty=format:|||BEGIN|||%H|||SEP|||%an|||SEP|||%ae|||SEP|||%ad",
			"--no-merges",
		)
		cmd.Dir = r.RepoPath
		stdout, err := cmd.StdoutPipe()
		if nil != err {
			return err
		}
		if err := cmd.Start(); err != nil {
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
				return err
			}

			deletionsString := bits[1]
			if deletionsString == "-" {
				deletionsString = "0"
			}
			deletions, err := strconv.Atoi(deletionsString)
			if err != nil {
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
				return errors.New("did not expect currect commit to be null")
			}

			if currectCommit.ChangedFiles == nil {
				// TODO maybe skip? does this break anything?
				return errors.New("did not expect currect commit changed files to be null")
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
	for a := 1; a <= len(r.userCommits); a++ {
		<-results
	}
	return nil
}

func (r *RepoExtractor) libraryWorker(commits <-chan *commit.Commit, results chan<- bool) error {
	extensionToLanguageMap := buildExtensionToLanguageMap(fileExtensionMap)
	for commit := range commits {
		libraries := map[string][]string{}
		for n, fileChange := range commit.ChangedFiles {
			extension := filepath.Ext(fileChange.Path)
			if extension == "" {
				continue
			}
			// remove the trailing dot
			extension = extension[1:]
			lang, ok := extensionToLanguageMap[extension]
			// We don't know extension, nothing to do
			if !ok {
				continue
			}

			commit.ChangedFiles[n].Language = lang

			cmd := exec.Command(r.GitPath,
				"show",
				fmt.Sprintf("%s:%s", commit.Hash, fileChange.Path),
			)
			cmd.Dir = r.RepoPath
			out, err := cmd.CombinedOutput()
			if err != nil {
				searchString1 := fmt.Sprintf("Path '%s' does not exist in '%s'", fileChange.Path, commit.Hash)
				searchString2 := fmt.Sprintf("Path '%s' exists on disk, but not in '%s'", fileChange.Path, commit.Hash)
				// means the file was deleted, skip
				if strings.Contains(string(out), searchString1) || strings.Contains(string(out), searchString2) {
					continue
				}
				return err
			}

			// For .m (Objective-C or Matlab) and .pl extensions we can't just use filenames
			// We need a proper way to actually detect language, so we use enry
			if strings.ToLower(extension) == "m" || strings.ToLower(extension) == "pl" {
				lang, _ = enry.GetLanguageByContent(fileChange.Path, out)
				commit.ChangedFiles[n].Language = lang
			}

			analyzer, err := librarydetection.GetAnalyzer(lang)
			if err != nil {
				continue
			}

			fileLibraries := analyzer.ExtractLibraries(string(out))
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

	repoDataPath := r.OutputPath + ".json"
	zipPath := r.OutputPath + ".json.zip"
	// Remove old files
	os.Remove(repoDataPath)
	os.Remove(zipPath)

	err := os.MkdirAll(r.OutputPath, 0755)
	if err != nil {
		log.Println("Cannot create directory. Error:", err.Error())
	}
	file, err := os.Create(repoDataPath)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
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
	url, err := Upload(r.OutputPath+".json.zip", r.repo.RepoName)
	if err != nil {
		return err
	}
	fmt.Printf("Go to this link in the browser => %s", url)
	return nil
}

type repo struct {
	RepoName         string   `json:"repo"`
	Emails           []string `json:"emails"`
	SuggestedEmails  []string `json:"suggestedEmails"`
	PrimaryRemoteURL string   `json:"primaryRemoteUrl"`
}

type req struct {
	Limit  int
	Offset int
}

func buildExtensionToLanguageMap(input map[string][]string) map[string]string {
	extensionMap := map[string]string{}
	for lang, extensions := range input {
		for _, extension := range extensions {
			extensionMap[extension] = lang
		}
	}
	return extensionMap
}

var fileExtensionMap = map[string][]string{
	"1C Enterprise":    {"bsl", "os"},
	"Apex":             {"cls"},
	"Assembly":         {"asm"},
	"Batchfile":        {"bat", "cmd", "btm"},
	"C":                {"c", "h"},
	"C++":              {"cpp", "cxx", "hpp", "cc", "hh", "hxx"},
	"C#":               {"cs"},
	"CSS":              {"css"},
	"Clojure":          {"clj"},
	"COBOL":            {"cbl", "cob", "cpy"},
	"CoffeeScript":     {"coffee"},
	"Crystal":          {"cr"},
	"Dart":             {"dart"},
	"Groovy":           {"groovy", "gvy", "gy", "gsh"},
	"HTML+Razor":       {"cshtml"},
	"EJS":              {"ejs"},
	"Elixir":           {"ex", "exs"},
	"Elm":              {"elm"},
	"EPP":              {"epp"},
	"ERB":              {"erb"},
	"Erlang":           {"erl", "hrl"},
	"F#":               {"fs", "fsi", "fsx", "fsscript"},
	"Fortran":          {"f90", "f95", "f03", "f08", "for"},
	"Go":               {"go"},
	"Haskell":          {"hs", "lhs"},
	"HCL":              {"hcl", "tf", "tfvars"},
	"HTML":             {"html", "htm", "xhtml"},
	"JSON":             {"json"},
	"Java":             {"java"},
	"JavaScript":       {"js", "jsx", "mjs", "cjs"},
	"Jupyter Notebook": {"ipynb"},
	"Kivy":             {"kv"},
	"Kotlin":           {"kt", "kts"},
	"LabVIEW":          {"vi", "lvproj", "lvclass", "ctl", "ctt", "llb", "lvbit", "lvbitx", "lvlad", "lvlib", "lvmodel", "lvsc", "lvtest", "vidb"},
	"Less":             {"less"},
	"Lex":              {"l"},
	"Liquid":           {"liquid"},
	"Lua":              {"lua"},
	"MATLAB":           {"m"},
	"Nix":              {"nix"},
	"Objective-C":      {"mm"},
	"OpenEdge ABL":     {"p", "ab", "w", "i", "x"},
	"Perl":             {"pl", "pm", "t"},
	"PHP":              {"php"},
	"PLSQL":            {"pks", "pkb"},
	"Protocol Buffer":  {"proto"},
	"Puppet":           {"pp"},
	"Python":           {"py"},
	"QML":              {"qml"},
	"R":                {"r"},
	"Raku":             {"p6", "pl6", "pm6", "rk", "raku", "pod6", "rakumod", "rakudoc"},
	"Robot":            {"robot"},
	"Ruby":             {"rb"},
	"Rust":             {"rs"},
	"Scala":            {"scala"},
	"SASS":             {"sass"},
	"SCSS":             {"scss"},
	"Shell":            {"sh"},
	"Smalltalk":        {"st"},
	"Stylus":           {"styl"},
	"Svelte":           {"svelte"},
	"Swift":            {"swift"},
	"TypeScript":       {"ts", "tsx"},
	"Vue":              {"vue"},
	"Xtend":            {"xtend"},
	"Xtext":            {"xtext"},
	"Yacc":             {"y"},
}
