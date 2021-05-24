package entities

import "strings"

// Repository is the internal representation of repository information
type Repository struct {
	// ID of the external repository. e.g. GitHub ID
	ID string
	// FullName the full name of the repo including the vendor. e.g. microsoft/vscode
	FullName string
	// Name this name will be used to save the artifact e.g. vscode
	Name string
	// CloneURL this URL needs to be used to clone the repo
	CloneURL string
}

// GetSafeFullName returns with a string that can be used as file name.
func (r *Repository) GetSafeFullName() string {
	return strings.Replace(r.FullName, "/", "_", -1)
}
