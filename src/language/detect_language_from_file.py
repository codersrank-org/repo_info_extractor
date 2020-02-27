import os
import sys
import subprocess
import json
from pygments.lexers import guess_lexer

def detect_language_from_file(file_path):
    if os.path.exists(file_path):
        try:
            with open(file_path, 'r') as file:
                code = file.read()
                language = guess_lexer(code)
                if language.name == 'Objective-C':
                    return language.name
                else:
                    return "MATLAB"
        except:
            return None