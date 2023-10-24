run:
	@go run ./cmd/main/main.go 
.PHONY: run

database-up:
	@cd ./sample-database && docker compose up -d
.PHONY: database-up

database-down:
	@cd ./sample-database && docker compose down
.PHONY: database-down

database-drop:
	@cd ./sample-database && go run drop.go
	.PHONY: database-drop