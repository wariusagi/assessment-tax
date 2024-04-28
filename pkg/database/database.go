package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dbUrl string) (*sql.DB, error) {
	// open connection
	log.Println("Database connecting...")
	var err error
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("open connection: %v", err)
	}

	// check connection
	log.Println("Database pinging...")
	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("ping: %v", err)
	}

	log.Println("Database connected!")

	err = initMasterData(db)
	if err != nil {
		return db, fmt.Errorf("init master data: %v", err)
	}

	log.Println("Initialized master data!")
	return db, nil
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database closed!")
	}
}

func initMasterData(db *sql.DB) error {
	err := CreateMasterTaxDeduction(db)
	if err != nil {
		return err
	}
	err = InsertMasterTaxDeduction(db)
	if err != nil {
		return err
	}
	return nil
}
