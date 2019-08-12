# What is it? (Beta Test)
This script is used to extract data from your private repo. The data is used to calculate your score on https://codersrank.io

CodersRank by default only considers public repositories, however, most developers have their code in private repositories. We want to give the chance to these developers to improve their scores too by adding their private repositories.

We can understand private repos are private because of a reason. This script extracts only the most important information from the repos:
- Number of inserted lines in each commit
- Number of deleted lines in each commit

Other information such as remote URLs, file names, emails, names are hashed. So we can know if two commits belong to the same file but we won't know the file name.

Moreover, the output is saved to your machine and you can check what data is extracted and you can decide whether you want to share it with us or not.

# How to use it
## OSX
First install pip if you don't have it
```
$ curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
$ python get-pip.py
```
```
$ git clone https://github.com/codersrankOrg/repo_info_extractor.git
$ cd repo_info_extractor
$ ./install.sh
$ ./run.sh path/to/repository
$ ls -al ./repo_data.json
```
## Linux
```
$ git clone https://github.com/codersrankOrg/repo_info_extractor.git
$ cd repo_info_extractor
$ ./install.sh
$ ./run.sh path/to/repository
$ ls -al ./repo_data.json
```
## Windows (in development)
```
*** First, be sure you have Python installed. You can download Python from https://www.python.org/downloads/. ***
git clone https://github.com/codersrankOrg/repo_info_extractor.git
cd repo_info_extractor
install.bat
python src\main.py path\to\repo
dir
```

# Roadmap
1. v0.3.0: Improve language recognition. The current dummy solution only checks the file extensions. 
1. v0.4.0: Recognize external libraries. The current script only considers the programming languages. 
