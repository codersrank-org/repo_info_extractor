docker build -t codersrank_extractor:latest .
docker run -it --mount type=bind,source="%1"/,target="/repo" --mount type=bind,source="%CD%"/,target="/app" codersrank_extractor:latest
