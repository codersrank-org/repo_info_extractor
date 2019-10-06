import hashlib as md5
import git
import os
from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from analyze_libraries import AnalyzeLibraries
from ui.questions import Questions
from obfuscator import obfuscate


def initialize(directory, skip_obfuscation, output, parse_libraries):
    repo = git.Repo(directory)
    ar = AnalyzeRepo(repo)
    q = Questions()

    print('Initialization...')
    for branch in repo.branches:
        ar.create_commits_entity_from_branch(branch.name)
    ar.flag_duplicated_commits()
    ar.get_commit_stats()
    r = ar.create_repo_entity(directory)

    # Ask the user if we cannot find remote URL
    if r.primary_remote_url == '':
        answer = q.ask_primary_remote_url(r)

    authors = [(c['name'], c['email']) for _, c in r.contributors.items()]

    identities_err = None
    identities = q.ask_user_identity(authors, identities_err)
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
        al = AnalyzeLibraries(r.commits, authors, repo.working_tree_dir)
        libs = al.get_libraries()

        # combine repo stats with libs used
        for i in range(len(r.commits)):
            c = r.commits[i]
            if c.hash in libs.keys():
                r.commits[i].libraries = libs[c.hash]

    if not skip_obfuscation:
        r = obfuscate(r)

    er = ExportResult(r)
    er.export_to_json(output)

def init_headless(directory, skip_obfuscation, output, parse_libraries, emails, user_commits, reponame):
    repo = git.Repo(directory)
    ar = AnalyzeRepo(repo)
    q = Questions()

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
    print('Setting the local user_names ::',r.local_usernames)
    r.repo_name = reponame

    if parse_libraries:
        # build authors from the selection
        al = AnalyzeLibraries(r.commits, authors, repo.working_tree_dir)
        libs = al.get_libraries()

        # combine repo stats with libs used
        for i in range(len(r.commits)):
            c = r.commits[i]
            if c.hash in libs.keys():
                r.commits[i].libraries = libs[c.hash]

    if not skip_obfuscation:
        r = obfuscate(r)


    er = ExportResult(r)
    er.export_to_json_headless(output)
    print('Successfully analysed the repo ==>'+reponame)
