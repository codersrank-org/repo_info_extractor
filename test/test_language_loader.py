from language import loader


def test_load_is_caching():
    assert "JavaScript" not in loader._cache.keys()
    loader.load("JavaScript")
    assert "JavaScript" in loader._cache.keys()
