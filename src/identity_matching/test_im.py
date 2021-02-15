from src.matching import IdentityMatcher
from src.model import DistanceModel
from src.preprocessor import DistancePreprocessor
from src.transform_shortlog import process_shortlog_line

import re
import json


def test_identity_matcher(seed_path, shortlog_path):
    """
    Test the identity matcher with extended logging. You have to provide the seed as a JSON
    file and the shortlog as a .txt file.
    """

    # setup preprocessor
    with open("./resources/domain_blacklist.txt", "r",
              encoding="utf-8") as f:
        domain_blacklist = list()
        for l in f.readlines():
            domain_blacklist.append(l.strip())

    preprocessor = DistancePreprocessor(domain_blacklist=domain_blacklist)
    model = DistanceModel(
        vectorizer_path="./resources/vectorizer.p")

    im = IdentityMatcher(preprocessor=preprocessor, model=model, threshold=0.85, debug=True)

    short_log = list()

    with open(shortlog_path, "r", encoding="utf-8", errors="replace") as f:
        for line in f.readlines():
            try:
                count, n, e = process_shortlog_line(line)
                test_data = (count, n, re.sub(r"[<>]", "", e))
                short_log.append(test_data)
            except ValueError:
                print("error ", line)

    with open(seed_path, mode="r", encoding="utf-8") as f:
        seed = json.load(f)

    emails_v2 = im.get_emails(seed, short_log)

    return short_log, seed, emails_v2


if __name__ == "__main__":

    short_log_file = "./test_data/shortlog.txt"
    seed_path = "./test_data/test_seed.json"

    sl, s, v2 = test_identity_matcher(seed_path, short_log_file)
    print("**********************************")
    print("Emails found for the current user:")
    print(v2.get("emails"))
