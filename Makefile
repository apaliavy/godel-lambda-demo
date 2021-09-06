#database information
DB_USERNAME ?= postgres
DB_PASSWORD ?= mysecretpassword
DB_NAME ?= demo
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5438
DB_OPTIONS ?= sslmode=disable

goose-status:
	goose postgres "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_OPTIONS)" status

goose-migrate:
	goose -dir tools/goose postgres "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_OPTIONS)" up

goose-degrade:
	goose -dir tools/goose postgres "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_OPTIONS)" down

goose-create:
	goose -dir tools/goose create $(FILENAME) sql

database-up:
	docker-compose up -d

database-down:
	docker-compose down --remove-orphans

binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/csv cmd/functions/csv-handler/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/xls cmd/functions/xls-handler/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/json cmd/functions/json-handler/main.go

zip: binary
	zip -j build/csv.zip bin/csv
	zip -j build/xls.zip bin/xls
	zip -j build/json.zip bin/json