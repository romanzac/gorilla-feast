BINARY_NAME=gorilla-feast

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v ./...

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

start-db:
	mkdir -p pgdata
	ln -s $(shell pwd)/pgdata /tmp/pgdata
	docker run --name postgr-gf -e POSTGRES_USER=gf_test -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=gf_test -p 5432:5432 -v "/tmp/pgdata:/var/lib/postgresql/data" -d postgres:15.1
	sleep 25
	cp scripts/gf_postgres_*.sql /tmp/pgdata
	# Installing schema from gf_postgres_schema.sql...
	docker exec -it postgr-gf bash -c "psql -U gf_test -h localhost -d gf_test -f /var/lib/postgresql/data/gf_postgres_schema.sql"
	# Installing initial data from gf_postgres_data.sql...
	docker exec -it postgr-gf bash -c "psql -U gf_test -h localhost -d gf_test -f /var/lib/postgresql/data/gf_postgres_data.sql"
	rm -rf /tmp/pgdata/gf_postgres_*.sql

gen-keys:
	./scripts/gen_tls_jwt_keys.sh localhost gorilla-feast
	./scripts/gen_tls_jwt_keys.sh gorilla-feast localhost

docker:
	./scripts/build_docker.sh
