package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	com "digibank/infrastructure/functions"
	"digibank/infrastructure/migration"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func NewDatabaseConnect(dir string) *sql.DB {
	com.GenerateRandomID()
	currDir := fmt.Sprint(dir, "/.env")
	err := godotenv.Load(currDir)
	if err != nil {
		log.Fatal("(INFRASTRUCTURE:1001): ", err)
	}

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"),
		os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"))

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		log.Fatal("(INFRASTRUCTURE:1002): ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("(INFRASTRUCTURE:1003): ", err)
	}

	migration.DBMigrate(db, "up")

	fmt.Println("Database is successfully connected")

	com.PrintLog("Database is successfully connected")

	return db
}
