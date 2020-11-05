import os
from language.PHP import extract_libraries


def test_extract_libraries():
    """
    Extract libraries from the library.

    Args:
    """
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/PHP.php']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert 'lib1' in libs["PHP"]
    assert 'lib2' in libs["PHP"]
    assert 'lib3' in libs["PHP"]
    assert 'lib4' in libs["PHP"]
    assert 'lib5' in libs["PHP"]
    assert 'lib6' in libs["PHP"]
    assert 'lib7' in libs["PHP"]
    assert 'lib8' in libs["PHP"]
    assert 'lib9' in libs["PHP"]
    assert 'lib10' in libs["PHP"]
    assert 'lib11' in libs["PHP"]
    assert 'Illuminate\\' in libs["PHP"]