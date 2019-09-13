import os
from language.JavaScript import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/JavaScript.js']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    print(libs)
    assert libs == [
        'lib1', 
        'lib2', 
        'lib3',
        'lib4',
        './lib5',
        './lib6',
        'lib7', 
        'lib8', 
        'lib9',
        'lib10',
        './lib11',
        './lib12',
        'lib13', 
        'lib14', 
        'lib15',
        'lib16',
        './lib17',
        './lib18',
        'lib19', 
        'lib20', 
        'lib21',
        'lib22',
        './lib23',
        './lib24',
        'lib25', 
        'lib26', 
        'lib27',
        'lib28',
        'lib29', 
        'lib30', 
        'lib31',
        'lib32',
        'lib33', 
        'lib34', 
        'lib35',
        'lib36',
        'lib37',
        'lib38',
        'lib39',
        'lib40',
        'lib41',
        'lib42',
        'lib43',
        'lib44',
        ], "Should extract all libraries used"