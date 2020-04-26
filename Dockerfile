FROM golang:1.14.2-alpine3.11
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
