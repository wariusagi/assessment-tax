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
		return fmt.Errorf("open connection: %v", err)
	}

	// check connection
	log.Println("Database pinging...")
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping: %v", err)
	}

	log.Println("Database connected!")

	err = initMasterData()
	if err != nil {
		return fmt.Errorf("init master data: %v", err)
	}

	log.Println("Initialized master data!")
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database closed!")
	}
}

func initMasterData() error {
	err := CreateMasterTaxDeduction()
	if err != nil {
		return err
	}
	err = InsertMasterTaxDeduction()
	if err != nil {
		return err
	}
	return nil
}
