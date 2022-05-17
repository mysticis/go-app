package middleware

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"

	"github.com/mysticis/go-dcktst-demo/demo"
)

func initDB() (*demo.Queries, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBUrl := os.Getenv("DB_URL")
	db, err := sql.Open("pgx", DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	GlobalDB := demo.New(db)

	return GlobalDB, nil
}
