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
    # regex to find imports like import Stuff
    regex_import = re.compile(r"import ((?!func|var|let|typealias|protocol|enum|class|struct)\s?[a-zA-Z0-9.]+)", re.IGNORECASE)

    # regex to find imports like import kind module.symbol
    # here kind can be func, var, let, typealias, protocol, enum, class, struct
    # this regex will return a list of (kind, module.symbol)
    regex_declarations = re.compile(r"import (func|var|let|typealias|protocol|enum|class|struct) ([a-zA-Z0-9.]+)", re.IGNORECASE)

    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:

            contents = ' '.join(fr.readlines())
            matches = regex_import.findall(contents)
            matches.extend([d[1] for d in regex_declarations.findall(contents)])

        if matches:
            res.extend(matches)

    # remove duplicates
    res = list(set(res))
    return {"Swift": res}