SHELL = /usr/bin/env bash

YELLOW=\033[0;33m
RED=\033[0;31m
WHITE=\033[0m
GREEN=\u001B[32m
CYAN=

PYTHON3_EXECUTABLE := $(shell command -v python3 2> /dev/null)
PIP3_EXECUTABLE := $(shell command -v pip3 2> /dev/null)

.PHONY: test
test:
	nose2
docker:
	docker build -t codersrank/repo_info_extractor:latest .


install_pip:
	@if [ "$(PYTHON3_EXECUTABLE)" != "" ]; then \
		echo -e "Found $(GREEN)python3$(WHITE)" ;\
		sudo apt install -y -q python3-venv ;\
		wait ;\
		if [ "$(PIP3_EXECUTABLE)" == "" ]; then \
			if [ -f "get-pip.py" ]; then \
			echo -e "Downloading latest $(GREEN)pip$(WHITE)" ;\
			wget https://bootstrap.pypa.io/get-pip.py ; \
			fi ;\
			python3 get-pip.py --user ;\
			wait ;\
		fi ;\
		wait ;\
		rm -rf ./.pyenv ;\
		python3 -m venv ./.pyenv ;\
	else \
		echo -e "Oops it seems python3 isn't installed" ;\
		echo -e "Run $(GREEN)sudo apt install python3$(WHITE)" ;\
		echo -e " AND THEN $(GREEN)make install$(WHITE)" ;\
		exit 0 ;\
	fi

install_requirements: 
	@( \
       source ./.pyenv/bin/activate; \
       pip3 install  --upgrade pip ; \
	   pip3 install -r requirements.txt; \
	   exit 0 ; \
	)

help:
	@./run.sh -h

install: install_pip install_requirements help

 