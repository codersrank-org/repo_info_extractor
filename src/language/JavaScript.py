import re

"""
Extract a list of JS libraries used from a fully qualified paths of files
"""


def extract_libraries(files):
    res = []
    # regex for
    # require('abc') as well as const lib = require('abc') and others
    regex1 = re.compile(r'require\(["\'](.+)["\']\);?\s', re.IGNORECASE)
    # ES6 imports
    regex2 = re.compile(r'import\s*(?:.+ from)?\s?\(?[\'"](.+)[\'"]\)?;?\s', re.IGNORECASE)
    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex1.findall(contents)
            matches.extend(regex2.findall(contents))

        if matches:
            res.extend(matches)
    return res
