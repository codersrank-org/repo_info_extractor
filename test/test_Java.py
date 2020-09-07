import os
from language.Java import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Java.java']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert 'org.bytedeco.opencv' in libs["Java"]
    assert 'org.bytedeco' in libs["Java"]
    assert 'java.awt.Color' not in libs["Java"]    