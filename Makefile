generated: build

test:
	@go test ./...

build: test
	export CGO_ENABLED=0 export GOOS=linux && export GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -X main.version=v2.3.2' -o repo_info_extractor_linux
	export CGO_ENABLED=0 export GOOS=darwin && export GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -X main.version=v2.3.2' -o repo_info_extractor_osx
	export CGO_ENABLED=0 export GOOS=windows && export GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -X main.version=v2.3.2' -o repo_info_extractor_windows.exe