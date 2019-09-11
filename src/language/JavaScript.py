import re

"""
Extract a list of JS libraries used from a fully qualified paths of files
"""
def extract_libraries(files):
    res = []
    #regex for 
    # require('abc') as well as const lib = require('abc') and others
    regex1 = re.compile(r'require\(\'(.+)\'\);?\s')
    for f in files:
        try:
            fr = open(f, 'r')
        except FileNotFoundError:
            continue
        contents = ' '.join(fr.readlines())
        match = regex1.search(contents)
        if match and match.group(1):
            res.append(match.group(1))
        fr.close()
    return res
        
