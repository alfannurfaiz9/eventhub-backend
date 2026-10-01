include ./.env

DB_URL=://$(DBUSER):$(DBPASS)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable
MIGRATION_PATH=db/migrations
SEEDER_PATH=db/seeds

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

seed-categories:
	@psql $(DB_URL) < $(SEEDER_PATH)/001_categories.sql