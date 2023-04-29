migrate_up:
	./migrate -path sql/migrations -database "postgresql://docker:docker@localhost:5432/bank?sslmode=disable" -verbose up

migrate_down:
	./migrate -path sql/migrations -database "postgresql://docker:docker@localhost:5432/bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: migrate_up migrate_down