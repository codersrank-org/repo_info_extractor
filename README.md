## What is it?
This script is used to extract data from your private repo. The data is used to calculate your score on https://codersrank.io

CodersRank by default only considers public repositories, however, most developers have their code in private repositories. We want to give the chance to these developers to improve their scores too by adding their private repositories.

We can understand private repos are private because of a reason. This script extracts only the most important information from the repos:
- Number of inserted lines in each commit
- Number of deleted lines in each commit

Other information such as remote URLs, file names, emails, names are hashed. So we can know if two commits belong to the same file but we won't know the file name.

Moreover, the output is saved to your machine and you can check what data is extracted and you can decide whether you want to share it with us or not.

## How does it work?
When a repository is analyzed two tools are used: this and [libraries](https://github.com/codersrank-org/libraries) repository. 
This repository is responsible to recognize the languages and export the imported libraries.
The [libraries](https://github.com/codersrank-org/libraries) contains a list of supported libraries, imports and technologies they belong to. 

### In short
- Language recognition: [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor/).
- Library recognition: [libraries](https://github.com/codersrank-org/libraries)

## How to use it
The repo_info_extractor is written in Go, so you can either clone the repo and compile the program or just download the binary and start using it.
```
git clone --depth 1 https://github.com/codersrank-org/repo_info_extractor.git
cd repo_info_extractor
go run . local --repo_path ./path/to/repo
```

### Binary approach (easiest)
If using this approach, download the binary from [releases](https://github.com/codersrank-org/repo_info_extractor/releases) and run it.

```
wget https://github.com/codersrank-org/repo_info_extractor/releases/download/vx.x.x/repo_info_extractor_osx # replace with the latest version
chmod +x repo_info_extractor_osx                                                                            # in case of Linux, OSX first make it executable
./repo_info_extractor_osx --repo_path ./path_to_repo
```
You can find a short video about the usage

[![How to use repo_info_extractor](https://img.youtube.com/vi/9IqgmYl8l2Y/0.jpg)](https://www.youtube.com/watch?v=9IqgmYl8l2Y)

### Available commands
You can see the available commands and flags with the `--help` flag. For example:
```
./repo_info_extractor_osx --help
...
./repo_info_extractor_osx bitbucket --help
```
Commands:
-  `bitbucket` Extract repository from BitBucket
-  `help` Help about any command
-  `local` Extract local repository by path
-  `version` Print the version number

The commands might have flags. For example `local` has:
`--repo_name` You can overwrite the default repo name. This name will be shown on the profile page.
`--repo_path` Path of the repo

## BitBucket
Right now only BitBucket Cloud is supported. For authentication your have to use your username
and create an app password. You can create it here: https://bitbucket.org/account/settings/app-passwords/.
The app password and username must be set via the `--password` and `--username` flags. Example usage:
```
./repo_info_extractor_osx bitbucket --username="peti2001" --password=xxxxxx --visibility=private --emails=karakas.peter@gmail.com
```
When you create the a new `app password` make sure you select all the necessary scopes.
![repo_scope](https://raw.githubusercontent.com/peti2001/multi_repo_extractor/master/docs/bitbucket-scope.png)
The safest way if you create an `app password` and use it instead of your user's password.

## Run UnitTests 
In the root directory of the repo, run the following command:

```
go test ./...
```
