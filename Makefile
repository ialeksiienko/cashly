run:
	go run cmd/main.go --config=config/config.yml

goose-path:
	export GOOSE_MIGRATION_DIR=internal/migration

goose-up: goose-path
	goose postgres "postgres://admin:admin@localhost:5432/dbname" up

goose-down: goose-path
	goose postgres "postgres://admin:admin@localhost:5432/dbname" down

docker-up:
	docker build -t cashly-bot .
	docker