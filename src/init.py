import git
from export_result import ExportResult
from analyze_repo import AnalyzeRepo
from ui.questions import Questions


def initialize(directory, skip_obfuscation, output, prompt):
    repo = git.Repo(directory)
    ar = AnalyzeRepo(repo, skip_obfuscation)
    q = Questions()

    print('Initialization...')
    for branch in repo.branches:
        ar.create_commits_entity_from_branch(branch.name)
    ar.flag_duplicated_commits()
    ar.get_commit_stats()
    r = ar.create_repo_entity(directory)

    # Ask the user if we cannot find remote URL
    if r.primary_remote_url == '' and prompt:
        answer = q.ask_primary_remote_url(r)

    if prompt:    
        identities = q.ask_user_identity(r)
        while len(identities['user_identity']) == 0:
            print('Please select at least one.')
            identities = q.ask_user_identity(r)
            r.local_usernames = identities['user_identity']
   
    er = ExportResult(r, prompt)
    er.export_to_json(output)
