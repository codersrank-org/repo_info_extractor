import os
import multiprocessing as mp
import sys
import pickle
import time

from entity.commit import Commit
from entity.repository import Repository
from entity.file_change import FileChange
from ui.progress import progress

class AnalyzeRepo:
    def __init__(self, repo):
        """
        Initialize the git commits.

        Args:
            self: (todo): write your description
            repo: (str): write your description
        """
        self.commit_stats = {}
        self.results = []
        self.prog = 0
        self.total = 0
        self.repo = repo
        self.commit_list = {}
        self.user_commits = {}


    def create_commits_entity_from_branch(self, branch):
        '''
        Extract the commits from a given branch
        '''

        n = 100
        commits = list(self.repo.iter_commits(branch, max_count=n))
        skip = 0
        while len(commits) > 0:
            for commit in commits:
                if commit.hexsha in self.commit_list:
                    break
                # Try to solve decoding special characters problems
                try:
                    self.commit_list[commit.hexsha] = Commit(
                        commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, branch)
                except:
                    print("Could not decode commit meta")
                    continue

                commit.tree = None
                commit.parents = None
                self.commit_stats[commit.hexsha] = commit
            skip += n
            commits = list(self.repo.iter_commits(
                branch, max_count=n, skip=skip))

    def analyse_master_user_commits(self, usercommits):
        """
        Analyse all commits for the commits.

        Args:
            self: (todo): write your description
            usercommits: (todo): write your description
        """
        for commitid in usercommits:
            try:
                commit = self.repo.commit(commitid)
                self.user_commits[commit.hexsha] = Commit(commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, "master")
            except:
                print("Cannot get commit sha: ", commitid)


    def create_repo_entity(self, repo_dir):
        """
        Create a new repo entity.

        Args:
            self: (todo): write your description
            repo_dir: (str): write your description
        """
        return Repository(os.path.basename(repo_dir.rstrip(os.sep)), self.repo, self.commit_list, self.user_commits)

    def get_commit_stats(self):
        """
        Commit stats.

        Args:
            self: (todo): write your description
        """
        cpu_count = mp.cpu_count()
        with mp.get_context("spawn").Pool(cpu_count) as pool:
            for h, commit in self.commit_list.items():
                if not commit.is_duplicated:
                    pool.apply_async(call_set_commit_stats, [h, self.commit_stats[h]], callback=self.callback_func)
            pool.close()
            pool.join()

        for result in self.results:
            self.commit_list[result['hash']].set_commit_stats(result['stats'], self.repo.working_dir)

    def flag_duplicated_commits(self):
        '''
        If the branch is not deleted the merge commits duplicates the changes. 
        This method detects these merge commits.
        '''
        self.total = len(self.commit_list)
        for hash in self.commit_list:
            if self.commit_list[hash].is_merge:
                count = 0
                for parent in self.commit_list[hash].parents:
                    if parent in self.commit_list:
                        count += 1
                if count > 1:
                    self.commit_list[hash].is_duplicated = True
                    self.total -= 1

    def callback_func(self, data):
        """
        Callback to the callback

        Args:
            self: (todo): write your description
            data: (todo): write your description
        """
        # Sanitize filenames because they might have weird characters
        # Also cast dict.keys() to the list() so we don't get Runtime Errors
        keys = list(data["stats"].items())
        for k, v in keys:
            sanitized_key = sanitize_filename(k)
            if sanitized_key != k:
                data["stats"][sanitized_key] = v
                data["stats"].pop(k, None)

        self.results.append(data)
        self.prog += 1
        progress(self.prog, self.total, 'Analyzing commits')


def call_set_commit_stats(h, commit):
    """
    Call commit_stats on the commit.

    Args:
        h: (todo): write your description
        commit: (str): write your description
    """
    # print('Analyze commit ' + commit.hash[:8] + ' from branch ' + commit.branch + ', date: ' + commit.created_at)
    ret = {'hash': h, 'stats': commit.stats.files}
    return ret


def sanitize_filename(path):
    """
    Sanitize path.

    Args:
        path: (str): write your description
    """
    if len(path) > 1 and not path[-1].isalnum():
        path = path[:-1] 
    return path