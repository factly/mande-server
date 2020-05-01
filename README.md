## Api documentation
 Install swag by `go get -u github.com/swaggo/swag/cmd/swag`
 
 For generating docs, run `swag init`  it will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).

 ## Development environment ( sample )

```
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=data_portal
DB_HOST=localhost
```

##  Run

To start  `go run main.go`

With docker run `docker-compose up`

Swagger-ui :- http://localhost:3000/swagger/index.html