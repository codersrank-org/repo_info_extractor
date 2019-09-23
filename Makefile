.PHONY: test
test:
	nose2
docker:
	docker build -t codersrank/repo_info_extractor:latest .