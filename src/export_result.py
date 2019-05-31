import json
from zipfile import ZipFile

class ExportResult:
    def __init__(self, data):
        self.data = data
    
    def export_to_json(self, file_name):
        f = open(file_name, 'w+')
        f.write(json.dumps(self.data.json_ready(), indent=4))
        f.close()
        # Zip the output
        print('Zip output file...')
        with ZipFile(file_name + '.zip','w', compression=8) as zip: 
            zip.write(file_name)
            zip.close()
        print('Result has has been generated in: ' + file_name + '.zip')