import re
import logging
import time
import os

"""
Extract a list of JS libraries used from a fully qualified paths of files
"""

js_logger = logging.getLogger("main.analyze_libraries.javascript")


def extract_libraries(files):
    res = []
    # regex for
    # require('abc') as well as const lib = require('abc') and others
    regex1 = re.compile(r'require\(["\'](.+)["\']\);?\s', re.IGNORECASE)
    # ES6 imports
    regex2 = re.compile(r'import\s*(?:.+ from)?\s?\(?[\'"](.+)[\'"]\)?;?\s', re.IGNORECASE)
    for f in files:
        size = os.stat(f).st_size / (1024**2)
        js_logger.info("Opening file {0}. The size of the file is {1:.4f} MB.".format(f, size))
        start1 = time.time()
        with open(file=f, mode='r', errors='ignore') as fr:
            end1 = time.time()
            lines = fr.readlines()
            js_logger.info("# of line in {} is {}.".format(f, len(lines)))
            start2 = time.time()
            contents = ' '.join(lines)
            matches = regex1.findall(contents)
            matches.extend(regex2.findall(contents))
            end2 = time.time()
            js_logger.info("Time spent on open for {0} is {1:.6f} seconds.".format(f, end1 - start1))
            js_logger.info("Time spent processing {0} is {1:.6f} seconds.".format(f, end2 - start2))
        if matches:
            res.extend(matches)
    return res
