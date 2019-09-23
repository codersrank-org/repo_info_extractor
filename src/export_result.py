import json
import os
from zipfile import ZipFile
from ui.questions import Questions
from upload import uploadRepo
import webbrowser


class ExportResult:
    def __init__(self, data):
        self.data = data

    def export_to_json(self, file_name):
        f = open(file_name, 'w+')
        f.write(json.dumps(self.data.json_ready(), indent=4))
        f.close()
        # Zip the output
        with ZipFile(file_name + '.zip', 'w', compression=8) as zip:
            zip.write(file_name)
            zip.close()

        print('Result has has been saved in: ' + file_name + '.zip')
        q = Questions()

        result = q.query_yes_no(
            'Do you want to upload the result to your profile automatically?')
        if result:
            response = uploadRepo(file_name + '.zip')
            if response is not None:
                reponame = self.data.repo_name
                url = 'https://profile.codersrank.io/repo?token=' + \
                    response['token']+'&reponame='+reponame
                print('Go to this link in the browser => ' + url)
                webbrowser.open(url)

        os.remove(file_name)
