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

    # regex to find imports like using System.IO;
    regex = re.compile(r"using\s?([/a-zA-z0-9.]+);", re.IGNORECASE)

    # regex to find imports like using WebApi = GitHub.DistributedTask.WebApi;
    regex_alias = re.compile(r"using\s?[/a-zA-Z0-9.]*\s?=\s?([/a-zA-Z0-9.]+);", re.IGNORECASE)

    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex.findall(contents)
            matches.extend(regex_alias.findall(contents))

        if matches:
            res.extend(matches)
    return {"C#": res}
