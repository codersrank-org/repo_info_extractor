generated: build

test:
	@go test ./...

build: test
	export CGO_ENABLED=0 export GOOS=linux && go build -a -tags netgo -ldflags '-w -X main.version=v1.1.0' -o repo_info_extractor_linux
	export CGO_ENABLED=0 export GOOS=darwin && go build -a -tags netgo -ldflags '-w -X main.version=v1.1.0' -o repo_info_extractor_osx
	export CGO_ENABLED=0 export GOOS=windows && go build -a -tags netgo -ldflags '-w -X main.version=v1.1.0' -o repo_info_extractor_windows.exe
	export GOOS=$GOOS_OLD

install:
	mkdir /var/opt/repo_info_extractor
	cp repo_info_extractor_linux /var/opt/repo_info_extractor/repo_info_extractor_linux
	ln -s /var/opt/repo_info_extractor/repo_info_extractor_linux /usr/bin/repo_info_extractor