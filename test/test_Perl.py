import os
from language.Perl import extract_libraries

def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Perl.pm']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert libs['Perl'] == [
        "strict",
        "Benchmark",
        "Carp",
        "sigtrap",
        "Sub::Module",
        "Import::This",
        "utf8",
        "warnings",
        "Foo::Bar",
        "Module",
    ]
