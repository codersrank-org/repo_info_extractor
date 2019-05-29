import os
import multiprocessing as mp

from entity.commit import Commit
from entity.repository import Repository
from entity.file_change import FileChange

# Needs to be remote because passing as a parameter would be very slow
commit_stats = {}

class AnalyzeRepo:
    def __init__(self, repo):
        self.repo = repo
        self.commit_list = {}
    
    def create_commits_entity_from_branch(self, branch):
        '''
        Extract the commits from a given branch
        '''
        global commit_stats

        n = 100
        commits = list(self.repo.iter_commits(branch, max_count=n)) 
        skip = 0
        while len(commits) > 0:
            for commit in commits:
                if commit.hexsha in self.commit_list:
                    break
                self.commit_list[commit.hexsha] = Commit(commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, branch)
                commit_stats[commit.hexsha] = commit
            skip += n
            commits = list(self.repo.iter_commits(branch, max_count=n, skip=skip))

    def create_repo_entity(self, repo_dir):
        return Repository(os.path.basename(repo_dir.rstrip(os.sep)), self.repo, self.commit_list)
    
    def get_commit_stats(self):
        results = []

        pool = mp.Pool(mp.cpu_count())
        for hash in self.commit_list:
            if not self.commit_list[hash].is_duplicated:
                results.append(pool.apply_async(call_set_commit_stats, args=(self.commit_list[hash],)))

        pool.close()
        pool.join()
        for result in results:
            ret = result.get()
            self.commit_list[ret['hash']].set_commit_stats(ret['stats'])

    def flag_duplicated_commits(self):
        '''
        If the branch is not deleted the merge commits duplicates the changes. 
        This method detects the these merge commits.
        '''

        for hash in self.commit_list:
            if self.commit_list[hash].is_merge:
                count = 0
                for parent in self.commit_list[hash].parents:
                    if parent in self.commit_list:
                        count += 1
                if count > 1:
                    self.commit_list[hash].is_duplicated = True

def call_set_commit_stats(commit):
    global commit_stats

    # print('Analyze commit ' + commit.hash[:8] + ' from branch ' + commit.branch + ', date: ' + commit.created_at)
    return {'hash': commit.hash, 'stats': commit_stats[commit.hash].stats.files}
