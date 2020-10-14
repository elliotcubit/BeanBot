package state

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var database *sql.DB

func init() {
	// Heroku gives us the URL in this form
	// Add option to connect to local PG
	var err error
	database, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to SQL database")
	} else {
		log.Println("Successfully connected to the SQL database")
	}
}
