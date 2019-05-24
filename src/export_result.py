import json

class ExportResult:
    def __init__(self, data):
        self.data = data
    
    def export_to_json(self, file_name):
        f = open(file_name, 'w+')
        f.write(json.dumps(self.data.json_ready(), indent=4))
        f.close()