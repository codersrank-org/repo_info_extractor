generated: build

test:
	@go test ./...

build: test
	export CGO_ENABLED=0 export GOOS=linux && go build -a -tags netgo -ldflags '-w' -o repo_info_extractor_linux
	export CGO_ENABLED=0 export GOOS=darwin && go build -a -tags netgo -ldflags '-w' -o repo_info_extractor_osx
	export CGO_ENABLED=0 export GOOS=windows && go build -a -tags netgo -ldflags '-w' -o repo_info_extractor_windows.exe
	export GOOS=$GOOS_OLD