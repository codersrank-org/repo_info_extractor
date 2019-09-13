import os
from language.JavaScript import extract_libraries

def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['test_fixtures/JavaScript.js']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    print(libs)
    assert libs == ['lib1', 'lib2', 'lib3', 'lib4', './lib5', './lib6'], "Should extract all libraries used"
    