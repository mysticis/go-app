migrateup:
	migrate -path db/migration -database "postgres://postgres:secret@localhost:5432/spin?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:secret@localhost:5432/spin?sslmode=disable" -verbose down

