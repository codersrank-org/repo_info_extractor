import os
from language.C import extract_libraries

def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/C.c']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert libs['C'] == [
        "assert.h",
        "complex.h",
        "ctype.h",
        "stdio.h",
        "float.h",
        "string.h",
        "hey/sup/iomanip.h",
        "Hey/ssup/math.h",
        "hello/how/stdlib.h",
        "great/wchar.h",
        "stdbool.h",
        "stdint.h",
        "hello12",
        "WsSup34",
        "heyYo3-lol"
    ]