import glob
import os
import multiprocessing as mp
import sys
import pickle
import time
import pathlib

from entity.commit import Commit
from entity.repository import Repository
from entity.file_change import FileChange
from ui.progress import progress

# Needs to be remote because passing as a parameter would be very slow
commit_stats = {}
prog = 0
total = 0
results = []
supported_library_languages = {
    'JavaScript': ['js', 'jsx', 'md'],
    'Python': ['py']
}

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
                commit.tree = None
                commit.parents = None
                commit_stats[commit.hexsha] = commit
            skip += n
            commits = list(self.repo.iter_commits(branch, max_count=n, skip=skip))

    def create_repo_entity(self, repo_dir):
        return Repository(os.path.basename(repo_dir.rstrip(os.sep)), self.repo, self.commit_list)
    
    def get_commit_stats(self):

        pool = mp.Pool(mp.cpu_count())
        for h, commit in self.commit_list.items():
            if not commit.is_duplicated:
                pool.apply_async(call_set_commit_stats, [h, commit_stats[h]], callback=callback_func)

        pool.close()
        pool.join()

        self.parse_libraries()

        # for h, commit in self.commit_list.items():
        #     if not commit.is_duplicated:
        #         files = x[]
        #         libraries[h] = self.parse_libraries(result['stats'])

        libraries = {}
        for result in results:
            self.commit_list[result['hash']].set_commit_stats(result['stats'], libraries)
        

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
    
    def parse_libraries(self):
        res = {}
        prog = 0
        for hash, commit in self.commit_list.items():
            if not commit.is_duplicated:
                prog += 1
                progress(prog, total, 'Analyzing libraries')
                res[hash] = self.parse_libraries_in_commit(commit)
        return res

    def parse_libraries_in_commit(self, commit):
        res = {}
        # Build a fully qualified paths of files modified in the commit
        files = [os.path.join(self.repo.git_dir, x) for x in commit_stats[commit.hash].stats.files]
        # Check out the given hash
        self.repo.git.checkout(commit.hash)
        [print(pathlib.Path(x).suffix[:-1]) for x in files]
        print('Checking out %s' % commit.hash)
        # print(files)
        print(pathlib.Path(files[0]).suffix)
        for lang, extensions in supported_library_languages.items():
            # we have extensions now, filter the list to only files with those extensions
            lang_files = list(filter(lambda x: pathlib.Path(x).suffix[1:] in extensions, files))
            if lang_files:
                # if we go to this point, there were files modified in the language we support
                # now we need to run regex for imports for every single of such file
                print(lang_files)

        #     files = [] # Files is a list of fully qualified paths of every single file is of a given language
        #     for ext in extensions:
        #         files.extend([f for f in glob.glob(self.repo.git_dir + "**/*." + ext, recursive=True)])
        #     print(files)

        #     res[lang] = []
        return res
        




def call_set_commit_stats(h, commit):
    # print('Analyze commit ' + commit.hash[:8] + ' from branch ' + commit.branch + ', date: ' + commit.created_at)
    ret = {'hash': h, 'stats': commit.stats.files}
    return ret


def callback_func(data):
    global results
    global prog

    results.append(data)
    prog += 1
    progress(prog, total, 'Analyzing commits')
