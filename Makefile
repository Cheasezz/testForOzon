.PHONY: db-up
db-up:
	docker run --name=test-db -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d --rm postgres

.PHONY: gqlgen
gqlgen:
	go run github.com/99designs/gqlgen generate