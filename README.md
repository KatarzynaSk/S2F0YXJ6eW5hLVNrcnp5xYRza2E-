
## Run project

### Set up postgres locally 

Run docker image:

```
docker run --name local-db -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres
``` 

### Create database schema
```
cat db_init.sql | docker exec -i local-db psql -U postgres -d postgres
```

### Run server

* Install go

```
go build .
./gocrud
```
