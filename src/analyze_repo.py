import os

from entity.commit import Commit
from entity.repository import Repository

class AnalyzeRepo:
    def __init__(self, repo):
        self.repo = repo
        self.commit_list = {}

    
    def create_commits_entity_from_branch(self, branch):
        '''
        Extract the commits from a given branch
        '''
        n = 100
        commits = list(self.repo.iter_commits(branch, max_count=n)) 
        skip = 0
        while len(commits) > 0:
            for commit in commits:
                if self.commit_list.has_key(commit.hexsha):
                    break
                self.commit_list[commit.hexsha] = Commit(commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, commit.stats.files, branch)
                print('Analyze commit ' + commit.hexsha + ' from branch ' + branch)
            skip += n
            commits = list(self.repo.iter_commits(branch, max_count=n, skip=skip))

    def create_repo_entity(self, repo_dir):
        return Repository(os.path.basename(repo_dir.rstrip(os.sep)), self.repo, self.commit_list)

    def flag_duplicated_commits(self):
        '''
        If the branch is not deleted the merge commits duplicates the changes. 
        This method detects the these merge commits.
        '''
        for hash in self.commit_list:
            if self.commit_list[hash].is_merge:
                count = 0
                for parent in self.commit_list[hash].parents:
                    if self.commit_list.has_key(parent):
                        count += 1
                if count > 1:
                    self.commit_list[hash].is_duplicated = True