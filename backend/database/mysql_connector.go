package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB is a global variable to hold the connection to the database.
// It's now of type *sqlx.DB instead of *sql.DB
var DB *sqlx.DB

func Connect() error {
	var err error

	// Get environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

	// Open database connection with sqlx
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	return nil
}
