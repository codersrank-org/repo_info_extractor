docker build -t codersrank_extractor:latest .
for /F %%i in ("%1") do @set repo_name=%%~nxi
docker run -it --mount type=bind,source="%1"/,target="/ %repo_name%" --mount type=bind,source="%CD%"/,target="/app" codersrank_extractor:latest /%repo_name%
