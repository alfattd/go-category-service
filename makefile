db_url=postgres://$(user):$(password)@localhost:5432/$(db)?sslmode=disable

migrate-create:
	migrate create -ext sql -dir postgres/migrations -seq create_$(name)_table

migrate-up:
	migrate -path postgres/migrations -database "$(db_url)" up

migrate-down:
	migrate -path postgres/migrations -database "$(db_url)" down 1
