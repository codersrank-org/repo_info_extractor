FROM golang:1.15-alpine
WORKDIR /go/src/app
COPY . .
RUN apk add make g++
RUN make build
ENTRYPOINT repo_info_extractor_linux