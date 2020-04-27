FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV DB_HOST "postgresdb"
ENV DB_NAME "postgres"
ENV DB_PASSWORD "postgres"
ENV DB_USER "postgres"
RUN go build -o main .
EXPOSE 3000
CMD ["./main"]