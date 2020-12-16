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
First of all, the script needs to be cloned.

```
git clone https://github.com/codersrankOrg/repo_info_extractor.git
cd repo_info_extractor
```

### Docker approach (easiest)
If using this approach, the host machine does not need to have any tools installed apart from Docker. Generating the repository information is as easy as:

#### OSX / Linux
```
./run-docker.sh <path to the repository>
```

### Windows
```
run-docker.bat <path to the repository>
```

### Using Python on the host machine approach
First, be sure you have Python3 and pip installed. You can download Python from https://www.python.org/downloads/ or https://www.anaconda.com/distribution/ (with preinstalled packages and pip). We only support Python3,
so if using pre-isntalled python, please check the version with:
```
python -V
```
#### OSX
```
$ git clone https://github.com/codersrankOrg/repo_info_extractor.git
$ cd repo_info_extractor
$ ./install.sh
$ ./run.sh path/to/repository
$ ls -al ./repo_data.json.zip
```
#### Linux
```
$ git clone https://github.com/codersrankOrg/repo_info_extractor.git
$ cd repo_info_extractor
$ ./install.sh
$ ./run.sh path/to/repository
$ ls -al ./repo_data.json.zip
```
#### Windows
```
git clone https://github.com/codersrankOrg/repo_info_extractor.git
cd repo_info_extractor
install.bat
python src\main.py path\to\repo
dir
```

## Extracting multiple repos
In case you have multiple repos and you don't want to extract them one-by-one check out this solution: https://github.com/codersrank-org/multi_repo_extractor

## Dockerfile
The provided Dockerfile builds an image that contains the Python script as well as its dependencies. To keep the final image size low, it leverages the 
multi-stage build functionality. The first stage installs the dependencies as well as all the required build tools. The second stage, runtime,
just copies over the installed dependencies so that they can be used by the script.

In order to build a new image out of it, run `make docker` on Mac/Linux or `build-docker.bat` on Windows. It should result in 
`codersrank/repo_info_extractor:latest` image.

## Troubleshooting

```
/usr/bin/env: ‘bash\r’: No such file or directory
```

If you see the following error on a Windows machine, this is due to git converting the line endings automatically. A repository level configuration has
been added to stop this from happening, but the repo needs to be hard reset:

```
git reset --hard
```

If this for some reason does not work, just remove the repository and clone it again.

## How to contribute?

### Set up working environment
We recommend using Python virtual environments. We only support Python3, but please test your code in all major versions of Python3 starting with Python3.5.


#### Run UnitTests 
First, you have to install nose2.
```
pip3 install nose2
```

After that use the make file to run the tests
```
make test
```
