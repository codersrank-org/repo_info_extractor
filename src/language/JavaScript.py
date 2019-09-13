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
            print("[WARN] File %s not found" % f)
            continue
        contents = ' '.join(fr.readlines())
        matches = regex1.findall(contents)
        if matches:
            res.extend(matches)
        fr.close()
    return res
        
