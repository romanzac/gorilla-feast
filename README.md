# Gorilla Feast - API example in Go

To clone the repository, use the following command:

```sh
git clone https://github.com/romanzac/gorilla-feast
```

Run Gorilla Feast with docker-compose:

```sh
cd gorilla-feast
make docker
docker-compose up
```

Docker will start two services:
- gorilla-feast-db with Postgres
- gorilla-feast with API listening at port 4439

Send a GET request to https://localhost:4439/ping to test API is ready

SignUp first user with POST request to https://localhost:4439/user

