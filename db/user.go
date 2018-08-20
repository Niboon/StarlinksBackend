package db

import (
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg"
	"fmt"
)

/*
	name: user
 */
type User struct {
	tableName 	struct{} 	`sql:"users"`
	ID			int 		`sql:"id,pk"`
	Token		string		`sql:"token,unique"`
}

func CreateUserTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createErr := db.CreateTable(&User{}, opts)
	if createErr != nil {
		fmt.Printf("Error while Creating table users, %v\n", createErr)
		return createErr
	}
	fmt.Println("Successfully Created Table users")
	return nil
}
