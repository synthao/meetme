# Kubernetes hello world

Build and push Docker image to https://hub.docker.com/
```shell
$ make dp
```

Run app
```shell
$ go run main.go
```

Apply all migrations
```shell
make migrate-up POSTGRES_USER=root POSTGRES_PASSWORD=123456 POSTGRES_PORT=5432 POSTGRES_DB=db_name
```

Build and push Docker image:
```shell
make docker-build-push TAG=v1.0.0
```