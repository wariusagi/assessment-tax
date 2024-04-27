package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dbUrl string) error {
	// open connection
	log.Println("Database connecting...")
	var err error
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return fmt.Errorf("open connection failed: %v", err)
	}

	// check connection
	log.Println("Database pinging...")
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	log.Println("Database connected!")

	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database closed!")
	}
}
