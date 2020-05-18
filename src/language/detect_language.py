import os
from .detect_language_from_file import detect_language_from_file 

supported_languages = {
    '1C Enterprise': ['bsl', 'os'],
    'Apex': ['cls'],
    'Assembly': ['asm'],
    'Batchfile': ['bat'],
    'C': ['c'],
    'C++': ['cpp', 'cxx'],
    'C#': ['cs'],
    'CSS': ['css'],
    'Clojure': ['clj'],
    'COBOL': ['cbl', 'cob', 'cpy'],
    'CoffeeScript': ['coffee'],
    'Crystal': ['cr'],
    'Dart': ['dart'],
    'Groovy': ['groovy', 'gvy', 'gy', 'gsh'],
    'HTML+Razor': ['cshtml'],
    'EJS': ['ejs'],
    'Elixir': ['ex', 'exs'],
    'Elm': ['elm'],
    'ERB': ['erb'],
    'Erlang': ['erl', 'hrl'],
    'F#': ['fs', 'fsi', 'fsx', 'fsscript'],
    'Fortran': ['f90', 'f95', 'f03', 'f08', 'for'],
    'Go': ['go'],
    'Haskell': ['hs', 'lhs'],
    'HTML': ['html', 'htm'],
    'JSON': ['json'],
    'Java': ['java'],
    'JavaScript': ['js', 'jsx'],
    'Jupyter Notebook': ['ipynb'],
    'Kotlin': ['kt', 'kts'],
    'Less': ['less'], 
    'Liquid': ['liquid'],
    'Lua': ['lua'],
    'MATLAB': ['m'],
    'Nix': ['nix'],
    'Objective-C': ['mm'],
    'OpenEdge ABL': ['p', 'ab', 'w', 'i', 'x'],
    'Perl': ['pl'],
    'PHP': ['php'],
    'PLSQL': ['pks', 'pkb'],
    'Protocol Buffer': ['proto'],
    'Python': ['py'],
    'R': ['r'],
    'Ruby': ['rb'],
    'Rust': ['rs'],
    'Scala': ['scala'],
    'SASS': ['sass'],
    'SCSS': ['scss'],
    'Shell': ['sh'],
    'Smalltalk': ['st'],
    'Stylus': ['styl'],
    'Svelte': ['svelte'],
    'Swift': ['swift'],
    'TypeScript': ['ts', 'tsx'],
    'Vue': ['vue'],
    'Xtend': ['xtend'],
    'Xtext': ['xtext'],
}

_ext_lang = {}

def _build_ext_lang_map():
    """
    For optimisation purposes, build ext -> language map. Supposed to run once and cache
    """
    if not _ext_lang:
        for lang, extensions in supported_languages.items():
            for ext in extensions:
                _ext_lang[ext] = lang

    return _ext_lang


def detect_language(file_path):
    parts = file_path.split(os.sep)
    file_name = parts[-1]

    if file_name == 'Dockerfile':
        return 'Dockerfile'
    if file_name == 'Makefile':
        return 'Makefile'

    ext = file_name.split('.')[-1].lower()
    if ext == 'm':
        lang = detect_language_from_file(file_path)
        if lang is not None:
            return lang
        return ''    

    if ext in _ext_lang:
        return _ext_lang[ext]

    return ''


# This ensures the ext to lang map is build upon the module import
if not _ext_lang:
    _build_ext_lang_map()
