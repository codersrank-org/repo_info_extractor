# What is it? (Beta Test)
This script is used to extract data from your private repo. This data is used to calculate your score on https://codersrank.io

CodersRank by defult only considers public repositoris however, most of the developers have their code in private repositories. We want to give the change to these developers too to improve their scores by adding their private repositories.

We can understand private repos are private beacuse of a reason. This script extract only the most important information from the repos:
- Number of inserted lines in each commit
- Number of deleted lines in each commit

Other information such as remote URL, file names, emails, names are hashed. So we can know if two commits belong to the same file but we don't know the file name.

Moreover, the output is saved to your machine and you can check what data is extracted and you can decided whether you want to share it with us or not. 

# How to use it
## Linux/Unix
```
$ git clone https://github.com/codersrankOrg/repo_info_extractor.git
$ cd repo_info_extractor
$ ./install.sh
$ ./run.sh path/to/repository
$ ls -al ./repo_data.json
```
## Windows (in development)
```
git clone https://github.com/codersrankOrg/repo_info_extractor.git
cd repo_info_extractor
install.bat
python src\main.py path\to\repo
dir
```

# Roadmap
1. v0.2.0: Add auto upload output file option. Now the output is saved as a ZIP file and you have to upload manauly, which is not very user friendly. 
1. v0.3.0: Improve language recognition. The current dumy solution only check the file extansions. 
1. v0.4.0: Recognize external libraries. The current script only considers the programming languages. 
