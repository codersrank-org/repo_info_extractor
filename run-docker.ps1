$current = (Get-Location).Path;
$target = $args[0];
$name = Split-Path $target -Leaf;

docker run -it --mount type=bind,source="$target",target="/$name" --mount type=bind,source="$current",target="/app" codersrank/repo_info_extractor:latest /$name
