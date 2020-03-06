import os
import json
from init import init_headless


# You can test other repos on your computer by changing this values
# Default values are for repo_info_extractor repository
repo_name = "repo_info_extractor"
# Use full paths instead of relative ones
repo_path = os.getcwd() 
test_files = {"py": "Python", "txt": "", "xml": "", "bsl":"1C Enterprise"}

local_username = "alimgiray@Alims-MacBook-Pro.local"
output_path = os.getcwd() + "/test/generated/result.json"


# Anaylze repo and create results
def analyze_repo():
    # because of nose2
    os.chdir('../')
    init_headless(repo_path, True, output_path, True, [local_username], False, [], repo_name, True, 5, 2, None)


def test_repo():
    analyze_repo()

    # output generated correctly
    assert os.path.exists(output_path) == True
    assert os.path.exists(output_path + ".zip") == True

    with open(output_path) as json_file:
        data = json.load(json_file)

        # fields set correctly
        assert (local_username in data["localUsernames"]) == True
        assert data["repoName"] == repo_name

        # there are commits
        assert len(data["commits"]) > 0

        for commit in data["commits"]:
            # file extensions are correct
            for f in commit["changedFiles"]:
                parts = f["fileName"].split(os.sep)
                file_name = parts[-1]
                ext = file_name.split('.')[-1].lower()
                detected_langauge = f["language"]

                for k, v in test_files.items():
                    if ext == k:
                        assert detected_langauge == v