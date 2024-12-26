postgres:
	docker run -d --name postgres-amole -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -v ~/docker/postgres-amole:/var/lib/postgresql/data -p 5432:5432 postgres:17-alpine
createdb:
	docker exec -it postgres-amole createdb --username=root --owner=root amole_db
dropdb:
	docker exec -it postgres-amole dropdb amole_db
migration:
	@if [ "$(name)" = "" ]; then \
		echo "Please provide migration name. Usage: make migration name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./script/migrations -seq $(name)
migrateup:
	migrate -path script/migrations -database "postgresql://root:secret@localhost:5432/amole_db?sslmode=disable" -verbose up
migratedown:
	migrate -path script/migrations -database "postgresql://root:secret@localhost:5432/amole_db?sslmode=disable" -verbose down
migrateup-aws:
	migrate -path script/migrations -database "postgresql://root:bRSvlQ6WXYKnIijwIPnp@amole-database.c9kw2g80qxj7.ap-southeast-2.rds.amazonaws.com:5432/amole_db" -verbose up

.PHONY: postgres17 createdb dropdb migrateup migratedown