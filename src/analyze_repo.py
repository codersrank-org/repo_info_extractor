import os
import multiprocessing as mp

from entity.commit import Commit
from entity.repository import Repository
from entity.file_change import FileChange
from ui.progress import progress

# Needs to be remote because passing as a parameter would be very slow
commit_stats = {}
prog = 0
total = 0
results = []

class AnalyzeRepo:
    def __init__(self, repo, skip_obfuscation):
        self.repo = repo
        self.commit_list = {}
        self.skip_obfuscation = skip_obfuscation
    
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
                self.commit_list[commit.hexsha] = Commit(commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, branch, self.skip_obfuscation)
                commit_stats[commit.hexsha] = commit
            skip += n
            commits = list(self.repo.iter_commits(branch, max_count=n, skip=skip))

    def create_repo_entity(self, repo_dir):
        return Repository(os.path.basename(repo_dir.rstrip(os.sep)), self.repo, self.commit_list)
    
    def get_commit_stats(self):
        with mp.Pool(mp.cpu_count()) as pool:
            for hash, commit in self.commit_list.items():
                if not commit.is_duplicated:
                    r = pool.apply_async(call_set_commit_stats, args=[commit, commit_stats], callback=callback_func)
                    r.wait(10)
                    print(r.ready())
                    print(r.successful())
                    print(r.get())
            print(results)

        for result in results:
            #ret = result.get()
            self.commit_list[result['hash']].set_commit_stats(result['stats'])

    def flag_duplicated_commits(self):
        '''
        If the branch is not deleted the merge commits duplicates the changes. 
        This method detects these merge commits.
        '''
        global total
        
        total = len(self.commit_list)

        for hash in self.commit_list:
            if self.commit_list[hash].is_merge:
                count = 0
                for parent in self.commit_list[hash].parents:
                    if parent in self.commit_list:
                        count += 1
                if count > 1:
                    self.commit_list[hash].is_duplicated = True
                    total -= 1


def call_set_commit_stats(commit, stats):
    # print('Analyze commit ' + commit.hash[:8] + ' from branch ' + commit.branch + ', date: ' + commit.created_at)
    try:
        ret = {'hash': commit.hash, 'stats': stats[commit.hash].stats.files}
        print("success")
        sys.stdout.flush()
        return ret
    except:
        print("Error at call_set_commit_stats.")
        sys.stdout.flush()
        return -1
    
    

def callback_func(data):
    global results
    global prog
    try:
        results.append(data)
        prog += 1
        progress(prog, total, 'Analyzing commits')
        print("Ok for callback")
        sys.stdout.flush()
    except:
        print("Error at callback")
        sys.stdout.flush()