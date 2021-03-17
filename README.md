## What is it?
This script is used to extract data from your private repo. The data is used to calculate your score on https://codersrank.io

CodersRank by default only considers public repositories, however, most developers have their code in private repositories. We want to give the chance to these developers to improve their scores too by adding their private repositories.

We can understand private repos are private because of a reason. This script extracts only the most important information from the repos:
- Number of inserted lines in each commit
- Number of deleted lines in each commit

Other information such as remote URLs, file names, emails, names are hashed. So we can know if two commits belong to the same file but we won't know the file name.

Moreover, the output is saved to your machine and you can check what data is extracted and you can decide whether you want to share it with us or not.

## How does it work?
When a repository is analyzed two repositories are used: this and [libraries](https://github.com/codersrank-org/libraries) repository. 
This repository is responsible to recognize the languages and export the imported libraries.
The [libraries](https://github.com/codersrank-org/libraries) contains a list of supported libraries, imports and technologies they belong to. 

### In short
- Language recognition: [repo_info_extractor](https://github.com/codersrank-org/repo_info_extractor/).
- Library recognition: [libraries](https://github.com/codersrank-org/libraries)

## How to use it
We are using go for repo_info_extractor, so you can either clone the repo and compile the program or just download the binary and start using it.

```
git clone --depth 1 https://github.com/codersrankOrg/repo_info_extractor.git
cd repo_info_extractor
go run . --repo_path ./path_to_repo
```

### Binary approach (easiest)
If using this approach, download the binary from releases and run it.

```
./repo_info_extractor --repo_path ./path_to_repo
```

### Available Flags

`--repo_path` string: Mandatory. Path of the repo which will be analyzed.

`--emails` string array: Optional. By default repo_info_extractor will ask you to choose your emails from all the emails found in commits. But if you know which emails you've used, you can provide them as a comma separated list,  (e.g. "one@mail.com,two@email.com")

`--gitPath` string: Optional. By default repo_info_extractor will try to find your git, but if you see an error related to "git not found", you can manually provide your git path.

## Extracting multiple repos
In case you have multiple repos and you don't want to extract them one-by-one check out this solution: https://github.com/codersrank-org/multi_repo_extractor


## Troubleshooting

...
## How to contribute?

### Set up working environment
We recommend using latest go version.

### Run UnitTests 
In the root directory of the repo, run the following command:

```
go test ./...
```
