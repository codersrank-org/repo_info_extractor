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
    # regex for imports like this: import _ "github.com/user/repo/..."
    regex1 = re.compile(r'import[\t ]*(?:[_\.].*)?[\t ]?\(?"(.+)"\)?;?\s', re.IGNORECASE)
    # regex for imports with alias. Like this: import alias1 "github.com/user/repo/..."
    regex2 = re.compile(r'import[\t ]*[a-z].+[\t ]?\(?"(.+)"\)?;?\s', re.IGNORECASE)
    # regex for multiline imports
    # regex3 = re.compile(r'import\s*\(\s*(.*?)\s*\)', re.MULTILINE)
    regex3 = re.compile(r'import\s*\(\s*(.*?)\s*\)', re.DOTALL|re.IGNORECASE)
    # Find libraries in a multi line import
    regex4 = re.compile(r'\"(.*?)\"', re.IGNORECASE)
    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex1.findall(contents)
            matches.extend(regex2.findall(contents))

            multiline_matches = (regex3.findall(contents))
            for multiline_match in multiline_matches:
                matches.extend(regex4.findall(multiline_match))

        if matches:
            res.extend(matches)
    return {"Go": res}
