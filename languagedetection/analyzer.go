package languagedetection

import (
	"path/filepath"
	"strings"

	"github.com/src-d/enry/v2"
)

// LanguageAnalyzer used for detecting programming language of a file
type LanguageAnalyzer struct {
	FileExtensionMap map[string]string
}

// NewLanguageAnalyzer constructor
func NewLanguageAnalyzer() *LanguageAnalyzer {
	return &LanguageAnalyzer{
		FileExtensionMap: buildExtensionToLanguageMap(fileExtensionMap),
	}
}

// Detect returns with the detected language
// If no language could be detected it will return with an empty string.
// The filePath is the path to the file.
// For some file types it reads the content of the file.
func (l *LanguageAnalyzer) Detect(filePath string, fileContent []byte) string {
	extension := filepath.Ext(filePath)
	if extension == "" {
		fileName := filepath.Base(filePath)

		if fileName == "Dockerfile" {
			return "Dockerfile"
		}

		if fileName == "Makefile" {
			return "Makefile"
		}

		return ""
	}
	// remove the trailing dot
	extension = extension[1:]
	if l.ShouldUseFile(extension) {
		return l.DetectLanguageFromFile(filePath, fileContent)
	} else {
		return l.DetectLanguageFromExtension(extension)
	}
}

// DetectLanguageFromExtension returns programming language based on files extension
// Works for most cases, but for some cases we have to use DetectLanguageFromFile
func (l *LanguageAnalyzer) DetectLanguageFromExtension(extension string) string {
	extension = strings.ToLower(extension)
	if val, ok := l.FileExtensionMap[extension]; ok {
		return val
	}
	return ""
}

// DetectLanguageFromFile returns programming language based on file itself
// It also needs filename to increase accuracy
func (l *LanguageAnalyzer) DetectLanguageFromFile(filePath string, fileContents []byte) string {
	lang, _ := enry.GetLanguageByContent(filePath, fileContents)
	// For some reason enry is too bad at detecting Perl files
	// However it can successfully detect Prolog files
	// So, if the extension is "pl" but enry couldn't detect the language, it is probably Perl
	if filepath.Ext(filePath) == ".pl" && lang == "" {
		return "Perl"
	}
	return lang
}

// ShouldUseFile determines if it is enough to use extension, or we should try to read the file
// to determine the language
func (l *LanguageAnalyzer) ShouldUseFile(extension string) bool {
	_, ok := extensionsWithMultipleLanguages[extension]
	return ok
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

var extensionsWithMultipleLanguages = map[string]bool{
	"m":  true, // Objective-C, Matlab
	"pl": true, // Perl, Prolog
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
