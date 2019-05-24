import argparse
import git
import os
from export_result import ExportResult
from analyze_repo import AnalyzeRepo

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('directory', help='Path to the repository. Example usage: run.sh path/to/directory')
parser.add_argument('--output', default='./repo_data.json', dest='output', help='Path to the JSON file that will contain the result. By default exports to the STDOUT.')
args = parser.parse_args()

repo = git.Repo(args.directory)
ar = AnalyzeRepo(repo)

for branch in repo.branches:
    ar.create_commits_entity_from_branch(branch.name)
ar.flag_duplicated_commits()

r = ar.create_repo_entity(args.directory)

er = ExportResult(r)
er.export_to_json(args.output)