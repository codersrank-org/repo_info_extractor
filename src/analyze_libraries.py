import pathlib
import os
import hashlib as md5
import tempfile
import shutil
import git
import uuid
from ui.progress import progress
from language.loader import load as load_language
from language.detect_language import supported_languages
from datetime import datetime



class AnalyzeLibraries:
    def __init__(self, commit_list, authors, basedir):
        self.commit_list = commit_list
        self.authors = authors
        self.basedir = basedir

    # Return a dict of commit -> language -> list of libraries
    def get_libraries(self):
        res = {}
        commits = _filter_commits_by_authors(self.commit_list, self.authors)
        # Before we do anything, copy the repo to a temporary location so that we don't mess with the original repo
        tmp_repo_path = _get_temp_repo_path()
        
        now = datetime.now()
        print("[%s] Copying the repository to a temporary location, this can take a while..." % now.strftime("%d/%m/%Y %H:%M:%S"))	

        shutil.copytree(self.basedir, tmp_repo_path, symlinks=True)
        now = datetime.now()
        print("[%s] Finished copying the repository" % now.strftime("%d/%m/%Y %H:%M:%S"))	

        # Initialise the next tmp directory as a repo and hard reset, just in case
        repo = git.Repo(tmp_repo_path)
        repo.git.clean('-fd')
        repo.git.checkout('master')
        repo.git.reset('--hard')

        prog = 0
        total = len(commits)

        for commit in commits:
            libs_in_commit = {}
            files = [os.path.join(tmp_repo_path, x.file_name)
                     for x in commit.changed_files]
            for lang, extensions in supported_languages.items():
                # we have extensions now, filter the list to only files with those extensions
                lang_files = list(filter(lambda x: pathlib.Path(
                    x).suffix[1:].lower() in extensions, files))
                if lang_files:
                    # if we go to this point, there were files modified in the language we support
                    # check out the commit in our temporary branch
                    repo.git.checkout(commit.hash)
                    # now we need to run regex for imports for every single of such file
                    # Load the language plugin that is responsible for parsing those files for libraries used
                    parser = load_language(lang)
                    # Only parse libraries if we support the current language
                    if parser:
                        if lang not in libs_in_commit.keys():
                            libs_in_commit[lang] = []

                        libs_in_commit[lang].extend(
                            parser.extract_libraries(lang_files))

            prog += 1
            progress(prog, total, 'Analyzing libraries')
            if libs_in_commit:
                res[commit.hash] = libs_in_commit

        shutil.rmtree(tmp_repo_path)
        return res


# Return only commits authored by provided obfuscated_author_emails
def _filter_commits_by_authors(commit_list, authors):
    return list(filter(lambda x: (x.author_name, x.author_email) in authors, commit_list))


def _get_temp_repo_path():
    return os.path.join(tempfile.gettempdir(), str(uuid.uuid4()))
