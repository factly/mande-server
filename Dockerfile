FROM golang:1.14.2-alpine3.11
WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
ENV DB_HOST "postgresdb"
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
