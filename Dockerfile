FROM golang:1.14.2-alpine3.11
WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go get -u github.com/cosmtrek/air
ENTRYPOINT air -d
