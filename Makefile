.PHONY: test
test:
	nose2
build-docker:
	docker build -t codersrank/repo_info_extractor:latest .