import hashlib as md5
import os


class FileChange:
    def __init__(self, file_name, deletions, insertions, language):
        """
        Sets language.

        Args:
            self: (todo): write your description
            file_name: (str): write your description
            deletions: (todo): write your description
            insertions: (todo): write your description
            language: (str): write your description
        """
        self.file_name = file_name
        self.deletions = deletions
        self.insertions = insertions
        self.language = language

    def json_ready(self):
        """
        Return the json - ready representation of this task.

        Args:
            self: (todo): write your description
        """
        return {
            'fileName': self.file_name,
            'language': self.language,
            'insertions': self.insertions,
            'deletions': self.deletions
        }
