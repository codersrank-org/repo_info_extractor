import os
from src.language.Kotlin import extract_libraries


def test_extract_libraries():
    dir_path = os.path.dirname(os.path.realpath(__file__))
    files = ['fixtures/Kotlin.kt']
    fq_files = [os.path.join(dir_path, f) for f in files]
    libs = extract_libraries(fq_files)
    assert len(libs) == 1

    assert 'syntax.lexer.Token' in libs["Kotlin"]
    assert 'syntax.tree' in libs["Kotlin"]