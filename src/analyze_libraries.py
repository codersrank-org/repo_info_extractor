from pprint import pprint
import pathlib
import os
import hashlib as md5
import tempfile
import shutil
import git
import uuid
from ui.progress import progress
import importlib


supported_library_languages = {
    'JavaScript': ['js', 'jsx'],
    #'Python': ['py'],
    #'Markdown': ['md'],
}

class AnalyzeLibraries:
    def __init__(self, commit_list, authors, basedir, skip_obfuscation):
        self.commit_list = commit_list
        self.authors = authors
        self.basedir = basedir
        self.skip_obfuscation = skip_obfuscation

    # Return a dict of commit -> language -> list of libraries
    def get_libraries(self):
        res = {}
        processed_authors = []
        if not self.skip_obfuscation:
            for email, name in self.authors:
                # This logic is in two places now...
                # Do we need to check for empty email?
                # This logic is duplicated, same as in commits.py. Move to obfuscator
                name_md5_hash = md5.md5()
                name_md5_hash.update(name.encode('utf-8'))
                email_md5_hash = md5.md5()
                email_md5_hash.update(email.encode('utf-8'))
                processed_authors.append({name_md5_hash.hexdigest(), email_md5_hash.hexdigest()})
        else:
            processed_authors = self.authors
 
        commits = _filter_commits_by_authors(self.commit_list, processed_authors)
        # Before we do anything, copy the repo to a temporary location so that we don't mess with the original repo
        tmp_repo_path = _get_temp_repo_path()
        shutil.copytree(self.basedir, tmp_repo_path)

        # Initialise the next tmp directory as a repo and hard reset, just in case
        repo = git.Repo(tmp_repo_path)
        repo.git.clean('-fd')
        repo.git.checkout('master')
        repo.git.reset('--hard')

        prog = 0
        total = len(commits)

        for commit in commits:
            libs_in_commit = {}
            files = [os.path.join(tmp_repo_path, x.file_name) for x in commit.changed_files]
            for lang, extensions in supported_library_languages.items():
                # we have extensions now, filter the list to only files with those extensions
                lang_files = list(filter(lambda x: pathlib.Path(x).suffix[1:] in extensions, files))
                if lang_files:
                    # if we go to this point, there were files modified in the language we support
                    # check out the commit in our temporary branch
                    repo.git.checkout(commit.hash)
                    print('Checking out %s' % commit.hash)
                    # now we need to run regex for imports for every single of such file
                    print(lang_files)
                    # Load the language plugin that is responsible for parsing those files for libraries used
                    # Keep the local cache of loaded language parsers
                    parser = _load_parser(lang)
                    if lang not in libs_in_commit.keys():
                        libs_in_commit[lang] = []

                    libs_in_commit[lang].extend(parser.extract_libraries(lang_files))
                    # pprint(parser_class)
                    # print(parser)

            # res[commit.hash] = list(dict.fromkeys(libs_in_commit))
            prog += 1
            progress(prog, total, 'Analyzing libraries')
            res[commit.hash] = libs_in_commit
    
        shutil.rmtree(tmp_repo_path)
        # Remove those commits without libraries
        # return {k: v for k, v in res.items() if v}
        pprint(res)
        return res


# Return only commits authored by provided obfuscated_author_emails
def _filter_commits_by_authors(commit_list, authors):
    return list(filter(lambda x: {x.author_name, x.author_email} in authors, commit_list))

def _get_temp_repo_path():
    return os.path.join(tempfile.gettempdir(), str(uuid.uuid4()))

# This could in fact be moved to a loader inside the languages package, nice and neat
def _load_parser(language):
    try:
        return importlib.import_module("language.%s" % language)
    except ImportError:
        print("Could not load a parser for %s" % language)
        exit(1)

