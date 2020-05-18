import os

from identity_matching.src.preprocessor import DistancePreprocessor
from identity_matching.src.model import DistanceModel
from identity_matching.src.matching import IdentityMatcher


def test_preprocessor():

    # Change path for Docker
    pwd = os.getcwd()
    if pwd == "/":
        pwd = "/app"
    with open(pwd + "/src/identity_matching/resources/domain_blacklist.txt", "r",
              encoding="utf-8") as f:
        domain_blacklist = list()
        for line in f.readlines():
            domain_blacklist.append(line.strip())

    preprocessor = DistancePreprocessor(domain_blacklist=domain_blacklist)
    assert preprocessor.transform("arrow@gmail.com") == "arrow@"
    assert preprocessor.transform("arrow.green@codersrank.io") == "arrow.green@codersrank.io"
    assert preprocessor.transform("arrow@127.0.0.1") == "arrow@"
    assert preprocessor.transform("arrow@e5c1c795-43da-0310-a71f-fac65c449510") == "arrow@"

    return


def test_matching():

    # Change path for docker
    pwd = os.getcwd()
    if pwd == "/":
        pwd = "/app"

    seed = {"user_name": "",
            "names": ["Clark Kent"],
            "emails": ["kent.clark@gotham.io"]}
    shortlog_list = [(1, "Clark E Kent", "superman@yahoo.com"),
                     (1, "Wayne Bruce", "batman@gotham.com")]

    output = {
        'user_name': '',
        'names': ['Clark Kent', 'Clark E Kent'],
        'emails': ['kent.clark@gotham.io', 'superman@yahoo.com']}
    with open(pwd+"/src/identity_matching/resources/domain_blacklist.txt", "r", encoding="utf-8") as f:
        domain_blacklist = list()
        for line in f.readlines():
            domain_blacklist.append(line.strip())

    preprocessor = DistancePreprocessor(domain_blacklist=domain_blacklist)
    model = DistanceModel(vectorizer_path=pwd+"/src/identity_matching/resources/vectorizer.p")

    im = IdentityMatcher(preprocessor=preprocessor, model=model, threshold=0.85)
    assert im.get_emails(seed, shortlog_list) == output

    return


if __name__ == "__main__":

    test_preprocessor()
    test_matching()

