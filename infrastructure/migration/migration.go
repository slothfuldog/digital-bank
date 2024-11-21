package migration

import (
	"database/sql"
	com "digibank/infrastructure/functions"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migration/*.sql
var dbMigrations embed.FS

var DbConnection *sql.DB

func DBMigrate(dbParam *sql.DB, direction string) {
	com.PrintLog("==== DBMigrate START ====")
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migration",
	}

	var migrateDirection migrate.MigrationDirection
	switch direction {
	case "up":
		migrateDirection = migrate.Up
	case "down":
		migrateDirection = migrate.Down
	default:
		com.PrintLog("Invalid migration direction. Use 'up' or 'down'.")
		return
	}

	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrateDirection)
	if errs != nil {
		com.PrintLog(fmt.Sprintf("(DBMigrate:1004) %s", errs))
	}

	DbConnection = dbParam

	com.PrintLog(fmt.Sprintf("Migration success, applied %d migrations!", n))
	com.PrintLog("==== DBMigrate END ====")
}
