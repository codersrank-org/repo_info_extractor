import os
from language.Go import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Go.go']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert 'library1' in libs["Go"]
    assert 'gitlab.com/username/reponame/library2' in libs["Go"]
    assert 'gitlab.com/username/library3' in libs["Go"]
    assert 'gitlab.com/username/reponame/library4' in libs["Go"]
    assert 'gitlab.com/username/library5' in libs["Go"]
    assert 'gitlab.com/username/reponame/library6' in libs["Go"]
    assert 'gitlab.com/username/library7' in libs["Go"]
    assert 'gitlab.com/username/reponame/library8' in libs["Go"]
    assert 'gitlab.com/username/library9' in libs["Go"]
    assert 'library10' in libs["Go"]
    assert 'gitlab.com/username/reponame/library11' in libs["Go"]
    assert 'gitlab.com/username/library12' in libs["Go"]
    assert 'gitlab.com/username/reponame/library13' in libs["Go"]
    assert 'gitlab.com/username/library14' in libs["Go"]
    assert 'gitlab.com/username/reponame/library15' in libs["Go"]
    assert 'gitlab.com/username/library16' in libs["Go"]
    assert 'gitlab.com/username/reponame/library17' in libs["Go"]
    assert 'gitlab.com/username/library18' in libs["Go"]