from HTMLParser import HTMLParser

class HTMLExtractor(HTMLParser):    
    def handle_starttag(self, tag, attrs):
        if tag == "script":
            print(attrs)
    #     print "Encountered a start tag:", tag
    # def handle_endtag(self, tag):
    #     print "Encountered an end tag :", tag
    # def handle_data(self, data):
    #     print "Encountered some data  :", data

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

    return {"HTML": ['a']}
