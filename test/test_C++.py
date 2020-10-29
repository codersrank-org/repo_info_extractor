import os
from language.C import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Cpp.cpp']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert libs['C++'] == [
        "common.h",
        "Accident.h",
        "Ped.h",
        "Pools.h",
        "World.h",
    ]