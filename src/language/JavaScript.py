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
        
    # remove relative imports
    res = [x for x in res if ".." not in x and "./" not in x] 

    return {"JavaScript": res}
