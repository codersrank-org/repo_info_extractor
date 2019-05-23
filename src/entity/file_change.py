import md5
import os

def obfuscate_filename(file_name):
        '''
        We want to obfuscate the whole path and the file name except the extension.
        Example: md5(src)/md5(components)/md5(profile).vue
        Example: md5(Makefile)
        '''
    
        parts = file_name.split(os.sep)
        obfuscated_parts = []
        if len(parts) >= 2:
            for i in range(0, len(parts) - 1):
                md5_hash = md5.new()
                md5_hash.update(parts[i])
                obfuscated_parts.append(md5_hash.hexdigest())
        
        obfuscated_file_name = ''
        extansion_index = parts[-1].rfind('.')
        if extansion_index != -1:
            md5_hash = md5.new()
            md5_hash.update(parts[-1][:extansion_index])
            obfuscated_file_name = md5_hash.hexdigest()
            obfuscated_file_name += '.' + parts[-1][extansion_index+1:]
        else:
            md5_hash = md5.new()
            md5_hash.update(parts[-1])
            obfuscated_file_name = md5_hash.hexdigest()
        obfuscated_parts.append(obfuscated_file_name)

        return os.sep.join(obfuscated_parts)

class FileChange:
    def __init__(self, file_name, deletions, insertions, language):
        self.file_name = file_name
        self.deletions = deletions
        self.insertions = insertions
        self.language = language
    
    def json_ready(self):
        return {
                'fileName': obfuscate_filename(self.file_name),
                'language': self.language,
                'insertions': self.insertions,
                'deletions': self.deletions
            }
