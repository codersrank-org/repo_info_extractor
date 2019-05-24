import json

def convert_remote_url(remote_url):
    '''
    The remote URL can be provided by SSH or HTTPS
    This function will convert it to HTTPS format
    '''

    index = remote_url.find('@')
    if index == -1:
        return remote_url
    
    return remote_url[index+1:].replace(':', '/')

class Repository:
    def __init__(self, repo_name, repo, commits):
        remotes = {}
        for remote in repo.remotes:
            for url in repo.remote(remote.name).urls:
                remotes[remote.name] = convert_remote_url(url)
        cr = repo.config_reader()
        # TODO: what if there is no origin?
        self.repo_name = repo_name
        self.remotes = remotes
        self.primary_remote_url = convert_remote_url(cr.get_value('remote "origin"', 'url'))
        self.number_of_branches = len(repo.branches)
        self.number_of_tags = len(repo.tags)
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
            'primaryRemoteUrl': self.primary_remote_url,
            'numberOfBranches': self.number_of_branches,
            'numberOfTags': self.number_of_tags,
            'commites': commites
        }