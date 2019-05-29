import os
import hashlib as md5
import datetime as dt

from .file_change import FileChange

def detect_language(file_path):
    parts = file_path.split(os.sep)
    file_name = parts[-1]

    if file_name == 'Dockerfile':
        return 'Dockerfile'
    if file_name == 'Makefile':
        return 'Makefile'

    extension = file_name.split('.')[-1].lower()

    if extension == 'c':
        return 'C'
    if extension == 'cpp' or extension == 'cxx':
        return 'C++'
    if extension == 'go':
        return 'go'
    if extension == 'json':
        return 'JSON'
    if extension == 'html' or extension == 'htm':
        return 'HTML'
    if extension == 'css':
        return 'CSS'
    if extension == 'java':
        return 'Java'
    if extension == 'js':
        return 'JavaScript'
    if extension == 'vue':
        return 'Vue'
    if extension == 'proto':
        return 'Protocol Buffer'
    
    return ''

class Commit:
    def __init__(self, author_name, author_email, created_at, hash, parents, branch):
        self.original_author_name = author_name
        self.original_author_email = author_email
        self.author_name = author_name
        self.author_email = author_email
        self.created_at = created_at.strftime("%Y-%m-%d %H:%M:%S")
        self.hash = hash
        self.parents = []
        for parent in parents:
            self.parents.append(parent.hexsha)
        self.is_merge = len(self.parents) >= 2
        self.branch = branch
        self.changed_files = []
        self.is_duplicated = False
        detect_language

        self.obfuscate()

    def set_commit_stats(self, stats):
        for f in stats:
            self.changed_files.append(FileChange(f, stats[f]['deletions'], stats[f]['insertions'], detect_language(f)))

    def obfuscate(self):
        md5_hash = md5.md5()
        md5_hash.update(self.author_name.encode('utf-8'))
        self.author_name = md5_hash.hexdigest()
        md5_hash = md5.md5()
        md5_hash.update(self.author_email.encode('utf-8'))
        self.author_email = md5_hash.hexdigest()

    def json_ready(self):
        changed_files = []
        for f in self.changed_files:
            changed_files.append(f.json_ready())
        data = {
            "authorName": self.author_name,
            "authorEmail": self.author_email,
  	        "createdAt": self.created_at,
  	        "commitHash": self.hash,
  	        "isMerge": self.is_merge,
  	        "parents": self.parents,
            "changedFiles": changed_files,
            "isDuplicated": self.is_duplicated
        }

        return data