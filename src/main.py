import argparse
import git
import os
import inquirer
from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from ui.select_primary_remote import Questions

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

r = ar.create_repo_entity(args.directory)
print(r.primary_remote_url)

if r.primary_remote_url == '':
    answer = q.ask_primary_remote_url(r)
    print(answer['remote_repo'])

er = ExportResult(r)
er.export_to_json(args.output)