import hashlib as md5
import os


def obfuscate(r):
    for i in range(len(r.commits)):
        o_user, o_email = _obfuscate_user(
            r.commits[i].author_name, r.commits[i].author_email)
        r.commits[i].author_name = o_user
        r.commits[i].author_email = o_email
        for n in range(len(r.commits[i].changed_files)):
            r.commits[i].changed_files[n].file_name = _obfuscate_filename(
                r.commits[i].changed_files[n].file_name)

    return r


def _obfuscate_user(name, email):
    '''
    Obfuscate user name and email
    '''
    if name is not None:
        name_md5_hash = md5.md5()
        name_md5_hash.update(name.encode('utf-8'))
        name = name_md5_hash.hexdigest()
    if email is not None:
        email_md5_hash = md5.md5()
        email_md5_hash.update(email.encode('utf-8'))
        email = email_md5_hash.hexdigest()
    return name, email


def _obfuscate_filename(file_name):
    '''
    We want to obfuscate the whole path and the file name except the extension.
    Example: md5(src)/md5(components)/md5(profile).vue
    Example: md5(Makefile)
    '''
    parts = file_name.split(os.sep)
    obfuscated_parts = []
    if len(parts) >= 2:
        for i in range(0, len(parts) - 1):
            md5_hash = md5.md5()
            md5_hash.update(parts[i].encode('utf-8'))
            obfuscated_parts.append(md5_hash.hexdigest())

    obfuscated_file_name = ''
    extansion_index = parts[-1].rfind('.')
    if extansion_index != -1:
        md5_hash = md5.md5()
        md5_hash.update(parts[-1][:extansion_index].encode('utf-8'))
        obfuscated_file_name = md5_hash.hexdigest()
        obfuscated_file_name += '.' + parts[-1][extansion_index+1:]
    else:
        md5_hash = md5.md5()
        md5_hash.update(parts[-1].encode('utf-8'))
        obfuscated_file_name = md5_hash.hexdigest()
    obfuscated_parts.append(obfuscated_file_name)

    return os.sep.join(obfuscated_parts)
