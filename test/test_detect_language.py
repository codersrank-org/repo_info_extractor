from language import detect_language


def test_ext_lang_map_built():
    # Just some extensions to ensure the ext to lang map was populated
    extensions = ['js', 'go', 'py']
    for ext in extensions:
        assert ext in detect_language._ext_lang


def test_languages_recognised():
    assert detect_language.detect_language("/tmp/some_file.js") == "JavaScript"
    assert detect_language.detect_language(
        "/tmp/some_file.jsx") == "JavaScript"
    assert detect_language.detect_language("/tmp/some_file.bat") == "Batchfile"
    assert detect_language.detect_language("/tmp/some_file.c") == "C"
    assert detect_language.detect_language("/tmp/some_file.cpp") == "C++"
    assert detect_language.detect_language("/tmp/some_file.cxx") == "C++"
    assert detect_language.detect_language("/tmp/some_file.cs") == "C#"
    assert detect_language.detect_language("/tmp/some_file.clj") == "Clojure"
    assert detect_language.detect_language(
        "/tmp/some_file.coffee") == "CoffeeScript"
    assert detect_language.detect_language("/tmp/some_file.ex") == "Elixir"
    assert detect_language.detect_language("/tmp/some_file.exs") == "Elixir"
    assert detect_language.detect_language("/tmp/some_file.go") == "Go"
    assert detect_language.detect_language("/tmp/some_file.lhs") == "Haskell"
    assert detect_language.detect_language("/tmp/some_file.html") == "HTML"
    assert detect_language.detect_language("/tmp/some_file.htm") == "HTML"
    assert detect_language.detect_language("/tmp/some_file.json") == "JSON"
    assert detect_language.detect_language("/tmp/some_file.java") == "Java"
    assert detect_language.detect_language(
        "/tmp/some_file.ipynb") == "Jupyter Notebook"
    assert detect_language.detect_language("/tmp/some_file.liquid") == "Liquid"
    assert detect_language.detect_language("/tmp/some_file.lua") == "Lua"
    assert detect_language.detect_language("/tmp/some_file.m") == "MATLAB"
    assert detect_language.detect_language(
        "/tmp/some_file.mm") == "Objective-C"
    assert detect_language.detect_language("/tmp/some_file.pl") == "Perl"
    assert detect_language.detect_language("/tmp/some_file.php") == "PHP"
    assert detect_language.detect_language(
        "/tmp/some_file.proto") == "Protocol Buffer"
    assert detect_language.detect_language("/tmp/some_file.py") == "Python"
    assert detect_language.detect_language("/tmp/some_file.r") == "R"
    assert detect_language.detect_language("/tmp/some_file.rb") == "Ruby"
    assert detect_language.detect_language("/tmp/some_file.rs") == "Rust"
    assert detect_language.detect_language("/tmp/some_file.scala") == "Scala"
    assert detect_language.detect_language("/tmp/some_file.scss") == "SCSS"
    assert detect_language.detect_language("/tmp/some_file.sh") == "Shell"
    assert detect_language.detect_language("/tmp/some_file.swift") == "Swift"
    assert detect_language.detect_language("/tmp/some_file.ts") == "TypeScript"
    assert detect_language.detect_language(
        "/tmp/some_file.tsx") == "TypeScript"
    assert detect_language.detect_language("/tmp/some_file.vue") == "Vue"
