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
    # regex to find imports like require('..')
    regex_require = re.compile(r'require\(["\'](.+)["\']\)', re.IGNORECASE)
    # or without parenthesis
    regex_require_without_parens = re.compile(r'require ["\'](.+)["\']', re.IGNORECASE)
    # regex to find imports like require_once('..')
    regex_require_once = re.compile(r'require_once\(["\'](.+)["\']\)', re.IGNORECASE)
    # or without parenthesis
    regex_require_once_without_parens = re.compile(r'require_once ["\'](.+)["\']', re.IGNORECASE)
    # regex to find imports like include('..')
    regex_include = re.compile(r'include\(["\'](.+)["\']\)', re.IGNORECASE)
    # or without parenthesis
    regex_include_without_parens = re.compile(r'include ["\'](.+)["\']', re.IGNORECASE)
    
    # regex to find imports with use (exclude App because it is generally used to import internal modules)
    regex_use = re.compile(r'use ((?!App\\)[a-zA-Z]+\\).*;', re.IGNORECASE)

    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex_require.findall(contents)
            matches.extend(regex_require_without_parens.findall(contents))
            matches.extend(regex_require_once.findall(contents))
            matches.extend(regex_require_without_parens.findall(contents))
            matches.extend(regex_require_once_without_parens.findall(contents))
            matches.extend(regex_include.findall(contents))
            matches.extend(regex_include_without_parens.findall(contents))
            matches.extend(regex_use.findall(contents))

        if matches:
            res.extend(matches)

    # remove duplicates
    res = list(set(res))
    return {"PHP": res}