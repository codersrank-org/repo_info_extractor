import json

class Repository:
    def __init__(self, repo_name, remotes, clone_url, number_of_branches, number_of_tags, commits):
        self.repo_name = repo_name
        self.remotes = remotes
        self.clone_url = clone_url
        self.number_of_branches = number_of_branches
        self.number_of_tags = number_of_tags
        self.commits = []
        for hash in commits:
            self.commits.append(commits[hash])
    
    def json_ready(self):
        commites = []
        for commit in self.commits:
            commites.append(commit.json_ready())
        return {
            'repoName': self.repo_name,
            'remotes': self.remotes,
            'cloneUrl': self.clone_url,
            'numberOfBranches': self.number_of_branches,
            'numberOfTags': self.number_of_tags,
            'commites': commites
        }