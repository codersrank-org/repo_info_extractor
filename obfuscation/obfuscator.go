package obfuscation

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/codersrank-org/repo_info_extractor/v2/commit"
)

// Obfuscate private info, like filename, username and emails
func Obfuscate(c *commit.Commit) {
	c.AuthorEmail = toMD5(c.AuthorEmail)
	c.AuthorName = toMD5(c.AuthorName)
	for _, filechange := range c.ChangedFiles {
		filechange.Path = obfuscateFile(filechange.Path)
	}
}

func toMD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func obfuscateFile(path string) string {
	obfuscatedPath := ""
	dirs := strings.Split(path, "/")
	for i, dir := range dirs {
		// If file
		if i == len(dirs)-1 {
			filenameParts := strings.Split(dir, ".")
			// Obfuscate file name
			obfuscatedPath += toMD5(filenameParts[0])
			// Leave extensions as it is
			if len(filenameParts) > 1 {
				for _, extension := range filenameParts[1:] {
					obfuscatedPath += "." + extension
				}
			}
		} else {
			obfuscatedPath += toMD5(dir) + "/"
		}
	}
	return obfuscatedPath
}
