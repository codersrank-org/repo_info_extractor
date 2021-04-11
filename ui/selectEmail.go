package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// SelectEmail shows a CLI select interface.
// The user has a chance to select the given emails from
// a predefined list (allEmails).
// At least one option must be selected
// The returning value is the selected emails.
func SelectEmail(allEmails []string) []string {
askForEmails:
	// TODO sort by alphabetical order (or frequency?)
	selectedEmailsWithNames := []string{}
	prompt := &survey.MultiSelect{
		Message: "Please choose your emails:",
		Options: allEmails,
		Filter: func(filterValue string, optValue string, optIndex int) bool {
			return strings.Contains(optValue, filterValue)
		},
	}
	err := survey.AskOne(prompt, &selectedEmailsWithNames, survey.WithKeepFilter(true))
	if err == terminal.InterruptErr {
		os.Exit(0)
	}

	if len(selectedEmailsWithNames) == 0 {
		fmt.Println("Please choose at least one email!")
		goto askForEmails
	}

	return selectedEmailsWithNames
}
