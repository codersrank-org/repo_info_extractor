import re


def extract_libraries(files):
    """Extracts a list of imports that were used in the files

    Parameters
    ----------
    files : []string
        Full paths to files that need to be analysed

    Returns
    -------
    dict
        imports that were used in the provided files, mapped against the language
    """

    res = []
    # regex to find imports like org.springframework (exlude standart java libraries)
    regex_import = re.compile(r'import ((?!java)[a-zA-Z0-9]*\.[a-zA-Z0-9]*)', re.IGNORECASE)
    # regex to find imports like org.springframework.boot
    regex_import_long = re.compile(r'import ((?!java)[a-zA-Z0-9]*\.[a-zA-Z0-9]*\.[a-zA-Z0-9]*)', re.IGNORECASE)


    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex_import.findall(contents)
            matches.extend(regex_import_long.findall(contents))

        if matches:
            res.extend(matches)

    # remove duplicates
    res = list(set(res))
    return {"Java": res}