
import re


def process_shortlog_line(line):
    line = re.sub(r"[\ufeff]+", "", line)
    num = re.findall(r"\A(.*)\t", line)[0].strip()
    name = re.findall(r"\t(.*) <", line)[0].strip()
    email = re.findall(r"<(.*)>", line)[0].strip()
    
    return num, name, email