import argparse
import git
import os
import json
from entity.commit import Commit
from entity.repository import Repository

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('directory', help='Path to the repository. Example usage: run.sh path/to/directory')

args = parser.parse_args()
print 'Repo dir:' + args.directory

def create_commits_entity_from_branch(repo, branch, commit_list):
    """
    Extract the commites from a given branch
    """
    
    commits = list(repo.iter_commits(branch, max_count=100)) 
    skip = 0
    while len(commits) > 0:
        for commit in commits:
            if commit_list.has_key(commit.hexsha):
                break
            commit_list[commit.hexsha] = Commit(commit.author.name, commit.author.email, commit.committed_datetime, commit.hexsha, commit.parents, commit.stats.files, branch)
            print('Get commit ' + commit.hexsha + ' from branch ' + branch)
        skip += 50
        commits = list(repo.iter_commits(branch, max_count=50, skip=skip))
    
    return commit_list

def create_repo_entity(repo, repo_dir, commits):
    remotes = {}
    for remote in repo.remotes:
        for url in repo.remote(remote.name).urls:
            remotes[remote.name] = url
    cr = repo.config_reader()
    # TODO: what if there is no origin?
    return Repository(os.path.basename(repo_dir), remotes, cr.get_value('remote "origin"', 'url'), len(repo.branches), len(repo.tags), commits)

commit_list = {}
repo = git.Repo(args.directory)

for branch in repo.branches:
    create_commits_entity_from_branch(repo, branch.name, commit_list)

r = create_repo_entity(repo, args.directory, commit_list)

# print os.path.basename((args.directory.rstrip(os.sep)))
# print r.json_ready()['repoName']
print json.dumps(r.json_ready(), indent=4)