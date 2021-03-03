# Data Portal Server

**Releasability:** [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=alert_status)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Reliability:** [![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=bugs)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Security:** [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=security_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Maintainability:** [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=sqale_index)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=code_smells)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  
**Other:** [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=ncloc)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=factly_data-portal-server) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=factly_data-portal-server&metric=coverage)](https://sonarcloud.io/dashboard?id=factly_data-portal-server)  

## Api documentation
 Install swag by `go get -u github.com/swaggo/swag/cmd/swag`
 
 For generating docs, run `swag init`  it will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).

 ## Development environment config vars ( sample )

```
DATABASE_HOST=postgres 
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=mande 
DATABASE_PORT=5432 
DATABASE_SSL_MODE=disable
MODE=development

MEILI_URL=http://meilisearch:7700
MEILI_KEY=password

RAZORPAY_KEY=<razorpay access key>
RAZORPAY_SECRET=<razorpay secret key>

KETO_URL=http://keto:4466
KAVACH_URL=http://kavach-server:8000
KRATOS_PUBLIC_URL=http://kratos:4433
OATHKEEPER_HOST=oathkeeper:4455

SUPER_ORGANISATION_TITLE=<Super organisation title>
DEFAULT_USER_EMAIL=<user email>
DEFAULT_USER_PASSWORD=<user password>
```

##  Run
To start  `go run main.go`  
With docker `docker build -t data-portal-server .`

Swagger UI (admin): http://localhost:7721/swagger/index.html

## Run Tests
`go test ./test/... -coverpkg ./action/... -coverprofile=cov.out && go tool cover -html=cov.out`