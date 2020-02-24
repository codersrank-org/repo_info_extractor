import os
import sys
import subprocess
import json

def detect_language_from_file(repo_dir, file_path, test=False):

    enry_path = '/app/src/bin/enry_linux'

    # This is for mac os beacuse we need to use different binaries for different platforms
    if sys.platform == "darwin":
        enry_path = os.getcwd() + '/src/bin/enry'

    full_path = repo_dir + "/" + file_path
    if os.path.exists(full_path):
        try:
            output = subprocess.run([enry_path, "--json", full_path], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
            if output.stderr is None:
                output = json.loads(output.stdout.strip())
                return output["language"]
            else:
                return None    
        except:
            return None