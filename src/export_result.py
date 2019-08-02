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
        with ZipFile(file_name + '.zip','w', compression=8) as zip: 
            zip.write(file_name)
            zip.close()

        print('Result has has been saved in: ' + file_name + '.zip')    
        q = Questions()

        result = q.query_yes_no('Would you like to upload it?')
        if result:
            response = uploadRepo( file_name + '.zip')
            if response is not None:
                reponame = os.path.splitext(file_name)[0]
                reponame = reponame.replace(".","")
                reponame = reponame.replace("/","")
                url = 'https://profile.codersrank.io/repo?token='+response['token']+'&reponame='+reponame
                # print("Open this url to upload the private repo:")
                # print(url)
                webbrowser.open(url)

        os.remove(file_name)