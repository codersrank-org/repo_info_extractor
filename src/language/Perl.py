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

    # regex to find imports
    regex = re.compile(r"(?:[^#]\s+)(?:use|require)[^\S\n]+(?:if.*,\s+)?[\"']?([a-zA-Z][a-zA-Z0-9:]*)[\"']?(?:\s+.*)?;")

    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex.findall(contents)

        if matches:
            res.extend(matches)
    return {"Perl": res}
