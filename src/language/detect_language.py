import os

supported_languages = {
    'JavaScript': ['js', 'jsx'],
    'Batchfile': ['bat'],
    'C': ['c'],
    'C++': ['cpp', 'cxx'],
    'C#': ['cs'],
    'Clojure': ['clj'],
    'CoffeeScript': ['coffee'],
    'HTML+Razor': ['cshtml'],
    'CSS': ['css'],
    'Elixir': ['ex', 'exs'],
    'Go': ['go'],
    'Haskell': ['lhs', 'lhs'],
    'HTML': ['html', 'htm'],
    'JSON': ['json'],
    'Java': ['java'],
    'JavaScript': ['js', 'jsx'],
    'Jupyter Notebook': ['ipynb'],
    'Liquid': ['liquid'],
    'Lua': ['lua'],
    'MATLAB': ['m'],
    'Objective-C': ['mm'],
    'Perl': ['pl'],
    'PHP': ['php'],
    'Protocol Buffer': ['proto'],
    'Python': ['py'],
    'R': ['r'],
    'Ruby': ['rb'],
    'Rust': ['rs'],
    'Scala': ['scala'],
    'SCSS': ['scss'],
    'Shell': ['sh'],
    'Swift': ['swift'],
    'TypeScript': ['ts', 'tsx'],
    'Vue': ['vue'],
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

    if ext in _ext_lang:
        return _ext_lang[ext]

    return ''


# This ensures the ext to lang map is build upon the module import
if not _ext_lang:
    _build_ext_lang_map()
