import importlib

_cache = {}


def load(language):
    """
    Loads the language from the cache.

    Args:
        language: (str): write your description
    """
    if language not in _cache.keys():
        try:
            mod = ".%s" % language
            # TODO! I really do not like this hardcoded
            _cache[language] = importlib.import_module(mod, 'language')
        except ImportError:
            _cache[language] = None
    return _cache[language]
