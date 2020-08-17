# Data Portal Server

**Releasability:** [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=alert_status)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Reliability:** [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=bugs)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Security:** [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=security_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Maintainability:** [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=sqale_index)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=code_smells)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Other:** [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=ncloc)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=coverage)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  

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

## Test

Go to the following url to get all the test details
http://localhost:8898/

We are using `goconvey` for manage out BDD
For more details visit http://goconvey.co/