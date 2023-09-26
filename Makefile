include .env
export
run:
	@go run main.go

migrate:
	@go run main.go migrate

drop-db:
	rm ./finance.db
