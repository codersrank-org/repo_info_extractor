import argparse
import git

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('directory', help='Path to the repository. Example usage: run.sh path/to/directory')

args = parser.parse_args()
print 'Repo dir:' + args.directory

# repo = git.Repo.clone_from(self._small_repo_url(), os.path.join(rw_dir, 'repo'), branch='master')

repo = git.Repo(args.directory)
commits = list(repo.iter_commits('master', max_count=5))
for commit in commits:
    print commit.author.name
    print commit.author.email
    print commit.committed_datetime
    print commit.hexsha
    print commit.parents
    files = commit.stats.files
    for f in files:
        print f
        print files[f]['deletions']
        print files[f]['insertions']
    # diffs = commit.diff()
    # for diff in diffs:
    #     print diff
    print '-------'