import json
import hashlib as md5


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

    # user_commit - consider only these user commits for extracting the repo information
    # emails - merge these emails with these emails extracted from the repo
    # reponame - name of the repo
    def __init__(self, repo_name, repo, commits, user_commits):
        """
        Initialize the commits.

        Args:
            self: (todo): write your description
            repo_name: (str): write your description
            repo: (todo): write your description
            commits: (todo): write your description
            user_commits: (todo): write your description
        """
        remotes = {}
        self.emails_v2 = []
        self.original_remotes = {}
        self.contributors = {}
        self.local_usernames = []
        for remote in repo.remotes:
            for url in repo.remote(remote.name).urls:
                remotes[remote.name] = convert_remote_url(url)
                self.original_remotes[remote.name] = url

        self.repo_name = repo_name
        self.remotes = remotes
        if 'origin' in self.remotes:
            self.primary_remote_url = convert_remote_url(
                self.remotes['origin'])
        else:
            self.primary_remote_url = ''
        self.number_of_branches = len(repo.branches)
        self.number_of_tags = len(repo.tags)
        self.commits = []
        for hash in commits:
            name = ""
            email = ""
            if commits[hash].original_author_name is not None:
                name = commits[hash].original_author_name
            if commits[hash].original_author_email is not None:
                email = commits[hash].original_author_email
            self.contributors[name + email] = {
                'name': name,
                'email': email
            }
            self.commits.append(commits[hash])

        if user_commits:
            print("Extracting usernames from user_commits..")
            usernames = set()
            key_list = list( user_commits.keys() )

            for hash in key_list:
                email = ""
                if hash in commits.keys():
                    if commits[hash].original_author_email is not None:
                        email = commits[hash].original_author_email
                    usernames.add(email)
                else:
                    print("Commit hash not found ==>",hash)
                    user_commits.pop(hash, None)

            self.local_usernames = list(usernames)

        self.obfuscate()

    def obfuscate(self):
        """
        Perform the md5 md5 hash

        Args:
            self: (todo): write your description
        """
        if self.primary_remote_url != '':
            md5_hash = md5.md5()
            md5_hash.update(self.primary_remote_url.encode('utf-8'))
            self.primary_remote_url = md5_hash.hexdigest()
        for remote in self.remotes:
            md5_hash = md5.md5()
            md5_hash.update(self.remotes[remote].encode('utf-8'))
            self.remotes[remote] = md5_hash.hexdigest()

    def json_ready(self):
        """
        : return : class :.

        Args:
            self: (todo): write your description
        """
        commits = []
        for commit in self.commits:
            commits.append(commit.json_ready())
        return {
            'repoName': self.repo_name,
            'localUsernames': self.local_usernames,
            'remotes': self.remotes,
            'primaryRemoteUrl': self.primary_remote_url,
            'numberOfBranches': self.number_of_branches,
            'numberOfTags': self.number_of_tags,
            'commits': commits,
            'emails_v2': self.emails_v2
        }
