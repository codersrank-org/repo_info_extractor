import argparse
import git
import os
from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from ui.questions import Questions

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('directory', help='Path to the repository. Example usage: run.sh path/to/directory')
parser.add_argument('--output', default='./repo_data.json', dest='output', help='Path to the JSON file that will contain the result. By default exports to the STDOUT.')
args = parser.parse_args()

repo = git.Repo(args.directory)
ar = AnalyzeRepo(repo)
q = Questions()

for branch in repo.branches:
    ar.create_commits_entity_from_branch(branch.name)
ar.flag_duplicated_commits()
ar.get_commit_stats()
r = ar.create_repo_entity(args.directory)

# Ask the user if we cannot find remote URL
if r.primary_remote_url == '':
    answer = q.ask_primary_remote_url(r)

identities = q.ask_user_identity(r)
while len(identities['user_identity']) == 0:
    print('Please select at least one.')
    identities = q.ask_user_identity(r)
r.local_usernames = identities['user_identity']
er = ExportResult(r)
er.export_to_json(args.output)