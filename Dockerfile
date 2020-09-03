FROM golang:1.14.2-alpine3.11

WORKDIR /app

COPY . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go" --command="./main -dsn=postgres://postgres:postgres@postgres:5432/dataportal?sslmode=disable -meili=http://localhost:7700 -meiliKey=password"