import os
import hashlib as md5
import datetime as dt
from .file_change import FileChange
from language import detect_language


class Commit:
    def __init__(self, author_name, author_email, created_at, hash, parents, branch):
        """
        Initialize the author.

        Args:
            self: (todo): write your description
            author_name: (str): write your description
            author_email: (str): write your description
            created_at: (todo): write your description
            hash: (todo): write your description
            parents: (todo): write your description
            branch: (todo): write your description
        """
        self.original_author_name = author_name
        self.original_author_email = author_email
        self.author_name = author_name
        self.author_email = author_email
        self.created_at = created_at.strftime("%Y-%m-%d %H:%M:%S %z")
        self.hash = hash
        self.parents = []
        for parent in parents:
            self.parents.append(parent.hexsha)
        self.is_merge = len(self.parents) >= 2
        self.branch = branch
        self.changed_files = []
        self.is_duplicated = False
        self.libraries = None

    def set_commit_stats(self, stats, repo_dir):
        """
        Sets the stats for a repo.

        Args:
            self: (todo): write your description
            stats: (todo): write your description
            repo_dir: (str): write your description
        """
        for f in stats:
            full_path = repo_dir + "/" + f
            self.changed_files.append(FileChange(
                f, stats[f]['deletions'], stats[f]['insertions'], detect_language.detect_language(full_path)))
            
    def json_ready(self):
        """
        Return a json ready ready to be ready.

        Args:
            self: (todo): write your description
        """
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

        if self.libraries:
            data["libraries"] = self.libraries

        return data
