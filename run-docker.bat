for /F "delims=" %%i in ("%1") do set "repo_name=%%~ni"
docker run -it --mount type=bind,source="%1",target="/%repo_name%" --mount type=bind,source="%CD%",target="/app" codersrank/repo_info_extractor:latest /%repo_name%
