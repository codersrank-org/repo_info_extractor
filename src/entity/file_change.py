import hashlib as md5
import os


class FileChange:
    def __init__(self, file_name, deletions, insertions, language):
        self.file_name = file_name
        self.deletions = deletions
        self.insertions = insertions
        self.language = language

    def json_ready(self):
        return {
            'fileName': self.file_name,
            'language': self.language,
            'insertions': self.insertions,
            'deletions': self.deletions
        }
