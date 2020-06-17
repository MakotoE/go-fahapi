# For CI pipeline
FROM golang:1.14-buster

RUN apt update && apt install -y wget curl bzip2
RUN wget https://download.foldingathome.org/releases/public/release/fahclient/debian-stable-64bit/v7.6/fahclient_7.6.13_amd64.deb
RUN dpkg -i fahclient_7.6.13_amd64.deb
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.26.0
RUN FAHClient --version
RUN golangci-lint --version