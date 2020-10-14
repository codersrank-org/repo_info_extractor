#!/usr/bin/env bash
REPO_NAME=$(basename $1)
docker run -i --mount type=bind,source="$1"/,target="/$REPO_NAME" --mount type=bind,source="$(pwd)"/,target="/app" codersrank/repo_info_extractor:latest ${@:2} /$REPO_NAME 
