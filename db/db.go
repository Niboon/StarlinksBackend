package db

import (
	"github.com/go-pg/pg"
	"fmt"
	"os"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		Database: "starlinks",
		User: "postgres",
		Password: "2Obvious",
		Addr: "localhost:5432",
	}
	db := pg.Connect(opts)
	if db == nil {
		fmt.Println("Error: Failed to Connect to Database")
		os.Exit(100)
	}
	fmt.Println("Connection to Database Successful")
	return db
}

func Close(db *pg.DB) {
	closeErr := db.Close()
	if closeErr != nil {
		fmt.Println("Error: Failed while Closing Database")
	} else {
		fmt.Println("Successfully Closed Database")
	}
}
