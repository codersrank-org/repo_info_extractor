import os
from language.Java import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Swift.swift']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert libs['Swift'] == [
        'Cocoa',
        'os.log',
        'StatsKit',
        'ModuleKit',
        'CPU',
        'Memory',
        'Disk',
        'Net',
        'Battery',
        'Sensors',
        'GPU',
        'Fans',
        'CoreServices.DictionaryServices',
        'Pentathlon.swim',
        'test.test'
        ]
