import re
import logging
import time
import os

js_logger = logging.getLogger("main.analyze_libraries.javascript")


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
        size = os.stat(f).st_size / (1024**2)
        js_logger.debug("Opening file {0}. The size of the file is {1:.4f} MB.".format(f, size))
        start1 = time.time()
        with open(file=f, mode='r', errors='ignore') as fr:
            end1 = time.time()
            lines = fr.readlines()
            js_logger.debug("# of line in {} is {}.".format(f, len(lines)))
            start2 = time.time()
            contents = ' '.join(lines)
            matches = regex1.findall(contents)
            matches.extend(regex2.findall(contents))

            if "// @flow" in contents:
                matches.extend(["flowjs"])

            end2 = time.time()
            js_logger.debug("Time spent on open for {0} is {1:.6f} seconds.".format(f, end1 - start1))
            js_logger.debug("Time spent processing {0} is {1:.6f} seconds.".format(f, end2 - start2))
        if matches:
            js_logger.debug("Library found in {}. The first 20 chars of matches is {}".format(f, matches[0][0:20]))
            res.extend(matches)

    # remove relative imports
    res = [x for x in res if ".." not in x and "./" not in x]

    return {"TypeScript": res}
