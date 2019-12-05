import os
from language.JavaScript import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/JavaScript.js']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1
    assert libs["JavaScript"] == [
        'lib1', 
        'lib2', 
        'lib3',
        'lib4',
        'lib7', 
        'lib8', 
        'lib9',
        'lib10',
        'lib13', 
        'lib14', 
        'lib15',
        'lib16',
        'lib19', 
        'lib20', 
        'lib21',
        'lib22',
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
        'lib45',
        'lib46',
        'lib47',
        'lib48',
        'lib49',
        'lib50',
        'lib51',
        'lib52',
        'lib53',
        'flow'
        ], "Should extract all libraries used"
