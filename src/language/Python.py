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
    # regex to find imports like from foo.foo1 import bar; returns foo.foo1
    regex1 = re.compile(r"from (.+) import", re.IGNORECASE)

    # regex to find imports like import bar as foo and from foo.foo1 import bar; returns bar
    regex2 = re.compile(r"import ([a-zA-Z0-9_-]+)(?:\s| as)", re.IGNORECASE)
    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex1.findall(contents)
            matches.extend(regex2.findall(contents))

        if matches:
            res.extend(matches)
    return {"Python": res}
