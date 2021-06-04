FROM golang:1.15-alpine
WORKDIR /usr/src/repo_info_extractor
COPY . .
RUN apk add make g++ git
RUN make build
RUN make install
ENTRYPOINT /usr/bin/repo_info_extractor