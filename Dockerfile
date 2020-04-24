FROM golang:latest
WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go build -o main .
EXPOSE 3000
ENV DB_HOST "postgresdb"
CMD ["./main"]