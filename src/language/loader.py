import importlib

_cache = {}

def load(language):
    if language not in _cache.items():
        try:
            mod  = ".%s" % language
            # TODO! I really do not like this hardcoded
            _cache[language] = importlib.import_module(mod, 'language')
        except ImportError:
            print("Could not load a parser for %s" % language)
            exit(1)
    return _cache[language]
