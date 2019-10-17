package main

import (
	//import from standard library
	"library1"
	//import from git repo
	"gitlab.com/username/reponame/library2"
	"gitlab.com/username/library3"

	//import with alias
	lib4 "gitlab.com/username/reponame/library4"
	lib5 "gitlab.com/username/library5"

	//blank identifier import
	_ "gitlab.com/username/reponame/library6"
	_ "gitlab.com/username/library7"

	//dot import
	. "gitlab.com/username/reponame/library8"
	. "gitlab.com/username/library9"
)

//import from standard library
import "library10"
//import from git repo
import "gitlab.com/username/reponame/library11"
import "gitlab.com/username/library12"

//import with alias
import lib4 "gitlab.com/username/reponame/library13"
import lib5 "gitlab.com/username/library14"

//blank identifier import
import _ "gitlab.com/username/reponame/library15"
import _ "gitlab.com/username/library16"

//dot import
import . "gitlab.com/username/reponame/library17"
import . "gitlab.com/username/library18"

//Don't find these libraries 
func main() {
	_ := "gitlab.com/username/library19"
	fmt.Println("import \"gitlab.com/username/library20\"")
	fmt.Println("import gitlab.com/username/library20)

}