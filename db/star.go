package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"fmt"
)

/*
	name: star
 */
type Star struct {
	tableName 	struct{} 	`sql:"stars"`
	ID			int 		`sql:"id,pk"`
	Name		string		`sql:"name"`
	UserID 		int 		`sql:"user_id"`
	Img 		string 		`sql:"img"`
	Link 		string 		`sql:"link"`
}

func CreateStarTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createErr := db.CreateTable(&Star{}, opts)
	if createErr != nil {
		fmt.Printf("Error while Creating table stars, %v\n", createErr)
		return createErr
	}
	fmt.Println("Successfully Created Table stars")
	return nil
}