import pathlib
import os
import hashlib as md5
import tempfile
import shutil
import git
import stat
import uuid
from ui.progress import progress
from language.loader import load as load_language
from language.detect_language import supported_languages
from datetime import datetime
import logging
import time

module_logger = logging.getLogger("main.analyze_libraries")


class AnalyzeLibraries:
    def __init__(self, commit_list, author_emails, basedir, skip, commit_size_limit, file_size_limit):
        self.commit_list = commit_list
        self.author_emails = author_emails
        self.basedir = basedir
        self.skip = skip
        self.commit_size_limit = commit_size_limit
        self.file_size_limit = file_size_limit

    # Return a dict of commit -> language -> list of libraries
    def get_libraries(self):
        res = {}
        commits = _filter_commits_by_author_emails(self.commit_list, self.author_emails)
        if not commits:
            _log_info("No commmits found for the authored by selected users")
            return res

        # Before we do anything, copy the repo to a temporary location so that we don't mess with the original repo
        tmp_repo_path = _get_temp_repo_path()

        _log_info("Copying the repository to a temporary location, this can take a while...")
        try:
            shutil.copytree(self.basedir, tmp_repo_path, symlinks=True)
        except shutil.Error as e:
            module_logger.debug("Shutil error messages: {}.".format(str(e)))
        _log_info("Finished copying the repository to", tmp_repo_path)

        # Initialise the next tmp directory as a repo and hard reset, just in case
        repo = git.Repo(tmp_repo_path)
        repo.git.clean('-fd')
        try:
            repo.git.checkout('master')
        except git.exc.GitCommandError as err:
            _log_info("Cannot checkout master on repository: ", err)

        try:
            repo.git.reset('--hard')
        except git.exc.GitCommandError as err:
            _log_info("Cannot reset repository: ", err)

        prog = 0
        total = len(commits)

        if not self.skip:
            _log_info("Skipping is set to False. All commits and files will be evaluated. This may take time.")
        else:
            _log_info("Commit size limit is {} MB and file size limit is {} MB.".format(
                self.commit_size_limit, self.file_size_limit))

        try:
            for commit in commits:
                start = time.time()
                module_logger.debug("Current commit hash is {}.".format(commit.hash))
                libs_in_commit = {}
                files = [os.path.join(tmp_repo_path, x.file_name)
                         for x in commit.changed_files]

                # if skip is not set to false in args, we may skip certain commits
                # Estimate the summed size of the changed files in the commit. If changed files sum more than 10 MB
                # or there are no changed files we recognize, we skip the commit (don't check out)
                est_size = _estimate_changed_file_size(files)
                module_logger.debug("Changed file list is {} MBs.".format(est_size))
                module_logger.debug("Skip is set to {}.".format(self.skip))
                if not self.skip or ((est_size < self.commit_size_limit) and _should_we_check_out(files)):

                    module_logger.debug("Checking out and analyzing commit.")
                    co_start = time.time()
                    try:
                        repo.git.checkout(commit.hash, force=True)
                    except Exception:
                        continue
                    co_end = time.time()
                    module_logger.debug("Checking out took {0:.6f} seconds.".format(co_end - co_start))

                else:
                    module_logger.debug("Skipping commit.")
                    prog += 1
                    progress(prog, total, 'Analyzing libraries')
                    continue

                for lang, extensions in supported_languages.items():
                    # we have extensions now, filter the list to only files with those extensions
                    lang_files = list(filter(lambda x: (pathlib.Path(
                        x).suffix[1:].lower() in extensions), files))
                    if lang_files:
                        module_logger.debug("Current language is {}, and extensions are{}".format(lang, extensions))
                        # if we go to this point, there were files modified in the language we support
                        # check out the commit in our temporary branch

                        # we need to filter again for files, that got deleted during the checkout
                        # we also filter out tiles, which are larger than 2 MB to speed up the process
                        if self.skip:
                            filter_func = (lambda x: os.path.isfile(x) and os.stat(x).st_size
                                           < self.file_size_limit * (1024**2))
                        else:
                            filter_func = (lambda x: os.path.isfile(x))

                        lang_files_filtered = list(filter(filter_func
                                                          , lang_files))

                        total_size = sum(os.stat(f).st_size for f in lang_files_filtered)
                        module_logger.debug("The number of files in lang_files_filtered"
                                            " is {0}, the total size is {1:.2f} MB".
                                            format(
                                                 len(lang_files_filtered), total_size / (1024 ** 2)
                                             ))
                        # now we need to run regex for imports for every single of such file
                        # Load the language plugin that is responsible for parsing those files for libraries used
                        parser = load_language(lang)
                        # Only parse libraries if we support the current language
                        if parser:
                            mapped_libs = parser.extract_libraries(lang_files_filtered).items()
                            for lang, libraries in mapped_libs:
                                if len(libraries) == 0:
                                    continue
                                if lang not in libs_in_commit.keys():
                                    libs_in_commit[lang] = []
                                libs_in_commit[lang].extend(libraries)

                prog += 1
                end = time.time()
                module_logger.debug("Time spent processing commit {0} was {1:.4f} seconds.".format(
                    commit.hash, end-start))

                progress(prog, total, 'Analyzing libraries')

                if libs_in_commit:
                    res[commit.hash] = libs_in_commit

        except (Exception, KeyboardInterrupt) as err:
            # make sure to clean up the tmp folder before dying
            _cleanup(tmp_repo_path)
            raise err

        _cleanup(tmp_repo_path)
        return res


def _should_we_check_out(file_list):

    for lang, extensions in supported_languages.items():
        lang_files = list(filter(lambda x: (pathlib.Path(x).suffix[1:].lower() in extensions), file_list))
        if lang_files:
            return True
    return False


def _estimate_changed_file_size(file_list):
    total_size = 0
    for file in file_list:
        try:
            total_size += os.stat(file).st_size / (1024**2)
        except (FileNotFoundError, Exception) as e:
            module_logger.debug("Error when getting file size {}".format(str(e)))
            continue
    return total_size


def _remove_readonly(func, path, _):

    """Clear the readonly bit and reattempt the removal"""
    os.chmod(path, stat.S_IWRITE)
    func(path)


def _cleanup(tmp_repo_path):
    _log_info("Deleting", tmp_repo_path)
    try:
        shutil.rmtree(tmp_repo_path, onerror=_remove_readonly)
    except (PermissionError, NotADirectoryError, Exception) as e:
        _log_info("Error when deleting {}".format(str(e)))


# Return only commits authored by provided obfuscated_author_emails
def _filter_commits_by_author_emails(commit_list, author_emails):
    _log_info("Filtering commits by emails: ", author_emails)
    return list(filter(lambda x: x.author_email in author_emails, commit_list))


def _get_temp_repo_path():
    return os.path.join(tempfile.gettempdir(), str(uuid.uuid4()))


def _log_info(message, *argv):
    timed_message = "[%s] %s" % (datetime.now().strftime("%d/%m/%Y %H:%M:%S"), message)
    print(timed_message, *argv)
