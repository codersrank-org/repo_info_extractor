import argparse
import hashlib as md5
import git
import os
from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from analyze_libraries import AnalyzeLibraries
from ui.questions import Questions
from obfuscator import obfuscate

if __name__ == '__main__':
    import multiprocessing
    multiprocessing.set_start_method('spawn', True)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        'directory', help='Path to the repository. Example usage: run.sh path/to/directory')
    parser.add_argument('--output', default='./repo_data.json', dest='output',
                        help='Path to the JSON file that will contain the result')
    parser.add_argument('--skip_obfuscation', default=False, dest='skip_obfuscation',
                        help='If true it won\'t obfuscate the sensitive data such as emails and file names. Mostly for testing purpuse')
    parser.add_argument('--parse_libraries', default=False,
                        dest='parse_libraries', help='If true, used libraries will be parsed')

    args = parser.parse_args()

    repo = git.Repo(args.directory)
    ar = AnalyzeRepo(repo)
    q = Questions()

    print('Initialization...')
    for branch in repo.branches:
        ar.create_commits_entity_from_branch(branch.name)
    ar.flag_duplicated_commits()
    ar.get_commit_stats()
    r = ar.create_repo_entity(args.directory)

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

    if args.parse_libraries:
        # build authors from the selection
        al = AnalyzeLibraries(r.commits, authors, repo.working_tree_dir)
        libs = al.get_libraries()

        # combine repo stats with libs used
        for i in range(len(r.commits)):
            c = r.commits[i]
            if c.hash in libs.keys():
                r.commits[i].libraries = libs[c.hash]

    if not args.skip_obfuscation:
        r = obfuscate(r)

    er = ExportResult(r)
    er.export_to_json(args.output)


if __name__ == "__main__":
    main()
