import subprocess
import re
import os

from identity_matching.src.matching import IdentityMatcher
from identity_matching.src.model import DistanceModel
from identity_matching.src.preprocessor import DistancePreprocessor
from identity_matching.src.transform_shortlog import process_shortlog_line


def match_emails(directory, seed):
    # setup processor
    with open(os.getcwd()+"/repo_info_extractor/src/identity_matching/resources/domain_blacklist.txt", "r", encoding="utf-8") as f:
        domain_blacklist = list()
        for l in f.readlines():
            domain_blacklist.append(l.strip())

    preprocessor = DistancePreprocessor(domain_blacklist=domain_blacklist)
    model = DistanceModel(vectorizer_path=os.getcwd()+"/repo_info_extractor/src/identity_matching/resources/vectorizer.p")

    im = IdentityMatcher(preprocessor=preprocessor, model=model, threshold=0.85)
        
    # generate shortlog
    short_log = list()
    short_log_file = directory + "/shortlog.txt"

    with open(short_log_file, "w+", encoding="latin-1", errors="replace") as outfile:
        try:
            subprocess.run(["git", "-C", directory, "shortlog", "-se"], stdout=outfile, universal_newlines=True, timeout=5)
        except subprocess.TimeoutExpired:
            print("Shortlog timeouted for ", directory)
        else:
            with open(short_log_file, "r", encoding="latin-1", errors="replace") as f:
                for line in f.readlines():
                    try:
                        count, n, e = process_shortlog_line(line)
                        test_data = (count, n, re.sub(r"[<>]", "", e))
                        short_log.append(test_data)
                    except ValueError:
                        print("error ", line)
            if len(short_log) > 1:
                short_log = short_log[1:]

    # convert seed
    seed_obj = {}
    seed_obj["user_name"] = ""
    seed_obj["names"] = list()
    seed_obj["emails"] = list()
    if seed is not None:
        seed_obj["user_name"] = seed.username
        seed_obj["names"] = list(seed.names)
        seed_obj["emails"] = list(seed.emails)

    # get results
    emails_v2 = im.get_emails(seed_obj, short_log)
    return emails_v2