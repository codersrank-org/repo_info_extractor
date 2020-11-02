import os
from language.Swift import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Swift.swift']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1
    assert len(libs['Swift']) == 15

    assert 'Cocoa' in libs['Swift']
    assert 'os.log' in libs['Swift']
    assert 'StatsKit' in libs['Swift']
    assert 'ModuleKit' in libs['Swift']
    assert 'CPU' in libs['Swift']
    assert 'Disk' in libs['Swift']
    assert 'Memory' in libs['Swift']
    assert 'Net' in libs['Swift']
    assert 'Battery' in libs['Swift']
    assert 'Sensors' in libs['Swift']
    assert 'GPU' in libs['Swift']
    assert 'Fans' in libs['Swift']
    assert 'CoreServices.DictionaryServices' in libs['Swift']
    assert 'Pentathlon.swim' in libs['Swift']
    assert 'test.test' in libs['Swift']
