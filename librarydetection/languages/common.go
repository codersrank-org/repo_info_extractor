package languages

import "regexp"

// executeRegexes is a helper function that executes a number of regexes and assumes each of them returns 1 group only
func executeRegexes(contents string, regexes []*regexp.Regexp) []string {
	var res []string
	for _, r := range regexes {
		matches := r.FindAllStringSubmatch(contents, -1)
		for _, match := range matches {
			if len(match) > 1 {
				libs := match[1:]
				res = append(res, libs...)
			}
		}
	}
	return res
}
