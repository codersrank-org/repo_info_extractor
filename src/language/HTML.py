from html.parser import HTMLParser

class HTMLExtractor(HTMLParser):    
    imports = {}

    def handle_starttag(self, tag, attrs):    
        if tag == "script":
            for attr, value in attrs :
                if attr == "src" and value.endswith("js"):
                    if "JavaScript" not in self.imports:
                        self.imports["JavaScript"] = []
                    self.imports["JavaScript"].append(value.split("/")[-1])
        
        if tag == "link":
            for attr, value in attrs :
                if attr == "href" and value.endswith("css"):
                    if "CSS" not in self.imports:
                        self.imports["CSS"] = []
                    self.imports["CSS"].append(value.split("/")[-1])

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
    extractor = HTMLExtractor()

    for f in files:
        with open(file=f, mode='r', errors='ignore') as fr:
            contents = ' '.join(fr.readlines())
            extractor.feed(contents)
    
    # dedup
    for lang, imports in extractor.imports.items():
        extractor.imports[lang] = list(dict.fromkeys(imports))

    return extractor.imports
    #return extractor.imports
