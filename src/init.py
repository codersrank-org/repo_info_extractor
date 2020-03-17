import hashlib as md5
import git
import logging

from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from analyze_libraries import AnalyzeLibraries
from ui.questions import Questions
from obfuscator import obfuscate
from timeout.timeout import timeout
from threading import Timer
import shutil

from identity_matching.matcher import match_emails


def initialize(directory, skip_obfuscation, output, parse_libraries, email, skip_upload, debug_mode, skip,
               commit_size_limit, file_size_limit):

    # Initialize logger
    logger = logging.getLogger("main")
    if debug_mode:
        logger.setLevel(logging.DEBUG)
        fh = logging.FileHandler('extractor_debug_info.log')
        fh.setLevel(logging.DEBUG)
        formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        fh.setFormatter(formatter)
        logger.addHandler(fh)
    else:
        logger.setLevel(logging.WARNING)

    logger.debug("Initialized main logger.")

    repo = git.Repo(directory)
    ar = AnalyzeRepo(repo)
    q = Questions()

    print('Analyzing repo under %s ...' % (directory))

    try:
        # Stop parsing if there are no branches
        if not repo.branches:
            print('No branches detected, will ignore this repo')
            return
    
        for branch in repo.branches:
            ar.create_commits_entity_from_branch(branch.name)
        ar.flag_duplicated_commits()
        ar.get_commit_stats()
        r = ar.create_repo_entity(directory)
    
        # Stop parsing if there are no remotes
        if not r.original_remotes:
            print('No remotes detected, will ignore this repo')
            return
    
        # Ask the user if we cannot find remote URL
        if r.primary_remote_url == '':
            answer = q.ask_primary_remote_url(r)
    
        if not r.contributors.items():
            print('No authors detected, will ignore this repo')
            return
    
        authors = [(c['name'], c['email']) for _, c in r.contributors.items()]
        identities = {}
        identities['user_identity'] = []
    
        # Stop parsing if there are no authors
        if len(authors) == 0:
            print('No authors detected, will ignore this repo')
            return
    
        identities_err = None
        identities = q.ask_user_identity(authors, identities_err, email)
        MAX_LIMIT = 50
        while len(identities['user_identity']) == 0 or len(identities['user_identity']) > MAX_LIMIT:
            if len(identities['user_identity']) == 0:
                identities_err = 'Please select at least one author'
            if len(identities['user_identity']) > MAX_LIMIT:
                identities_err = 'You cannot select more than', MAX_LIMIT
            identities = q.ask_user_identity(authors, identities_err)
        r.local_usernames = identities['user_identity']
        
        if parse_libraries:
            # build authors from the selection
            # extract email from name -> email list
            author_emails = [i.split(' -> ', 1)[1] for i in r.local_usernames]

            if author_emails:
                al = AnalyzeLibraries(r.commits, author_emails, repo.working_tree_dir,
                                      skip, commit_size_limit, file_size_limit)
                libs = al.get_libraries()
                # combine repo stats with libs used
                for i in range(len(r.commits)):
                    c = r.commits[i]
                    if c.hash in libs.keys():
                        r.commits[i].libraries = libs[c.hash]
            
        if not skip_obfuscation:
            r = obfuscate(r)
        er = ExportResult(r)
        er.export_to_json_interactive(output, skip_upload)
    except KeyboardInterrupt:
        print("Cancelled by user")
        return

# user_commit - consider only these user commits for extracting the repo information
# emails - merge these emails with these emails extracted from the repo
# reponame - name of the repo


def init_headless(directory, skip_obfuscation, output, parse_libraries, emails, debug_mode, user_commits, reponame,
                  skip, commit_size_limit, file_size_limit, seed, timeout_seconds=600):
    # Initialize logger
    logger = logging.getLogger("main")
    if debug_mode:
        logger.setLevel(logging.DEBUG)
        fh = logging.FileHandler('extractor_debug_info.log')
        fh.setLevel(logging.DEBUG)
        formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        fh.setFormatter(formatter)
        logger.addHandler(fh)
    else:
        logger.setLevel(logging.WARNING)

    repo = git.Repo(directory)
    ar = AnalyzeRepo(repo)
    q = Questions()
    timer = Timer(timeout_seconds, timeout)
    timer.start()
    # Use a context manager with signal to measure seconds, and timeout
    try:
        print('Initialization...')
        for branch in repo.branches:
            ar.create_commits_entity_from_branch(branch.name)
        ar.flag_duplicated_commits()
        ar.get_commit_stats()
        print('Analysing the master branch..')
        ar.analyse_master_user_commits(user_commits)
        print('Creating the repo entity..')
        r = ar.create_repo_entity(directory)

        r.local_usernames = list(set(r.local_usernames + emails))
        MAX_EMAIL_LIMIT = 50
        if len(r.local_usernames) > MAX_EMAIL_LIMIT:
            print("Email count (" + str(len(r.local_usernames)) + ") for this repo exceeds the limit of " + str(MAX_EMAIL_LIMIT) + " emails.")
            r.local_usernames = r.local_usernames[0:MAX_EMAIL_LIMIT]
        print('Setting the local user_names ::',r.local_usernames)
        r.repo_name = reponame

        if parse_libraries and len(ar.commit_list) > 0:
            # build authors from the the email list provided
            # we are provided only emails in the headless mode
            # TODO! Support both name -> email and email formats
            author_emails = []
            for email in r.local_usernames:
                author_emails.append(email)

            if author_emails:
                al = AnalyzeLibraries(r.commits, author_emails, repo.working_tree_dir,
                                      skip, commit_size_limit, file_size_limit)
                libs = al.get_libraries()

                # combine repo stats with libs used
                for i in range(len(r.commits)):
                    c = r.commits[i]
                    if c.hash in libs.keys():
                        r.commits[i].libraries = libs[c.hash]

            # new email detection
            try:
                emails_v2 = match_emails(directory, seed)
                r.emails_v2 = emails_v2["emails"]
            except:
                r.emails_v2 = list()

        if not skip_obfuscation:
            r = obfuscate(r)

        er = ExportResult(r)
        er.export_to_json_headless(output)
        print('Successfully analysed the repo ==>'+reponame)
    except KeyboardInterrupt:
        print("{} timeouted after {} seconds.".format(repo.working_dir, timeout_seconds))
        print("Deleting", repo.working_dir)
        try:
            shutil.rmtree(repo.working_dir)
        except (PermissionError, NotADirectoryError, Exception) as e:
            print("Error when deleting {}. Message: {}".format(repo.working_dir, str(e)))
    finally:
        timer.cancel()

