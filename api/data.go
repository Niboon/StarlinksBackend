package api

import (
	"fmt"
	"github.com/go-pg/pg"
)

func insertUser(user *User, database *pg.DB) {
	err := database.Insert(user)
	if err != nil {
		fmt.Printf("Failed to Insert user with ID: %v : %v\n", user.ID, err)
	}
}

func insertStar(star *Star, database *pg.DB) {
	err := database.Insert(&star)
	if err != nil {
		fmt.Printf("Failed to Insert star with ID: %v : %v\n", star.ID, err)
	}
}

func selectAllStars(database *pg.DB) []Star {
	var stars []Star
	err := database.Model(&stars).Select()
	if err != nil {
		fmt.Printf("Failed to Select all users : %v\n", err)
	}
	return stars
}

func selectUserByID(id int, database *pg.DB) User {
	user := User{ID:id}
	err := database.Select(&user)
	if err != nil {
		fmt.Printf("Failed to Select user with ID: %v : %v\n", id, err)
	}
	return user
}

func selectOrInsertUserByToken(token string, database *pg.DB) User {
	user := User{Token:token}
	_, err := database.Model(&user).Where("token=?",token).SelectOrInsert()
	if err != nil {
		fmt.Printf("Failed to Select user with Token: %v : %v\n", token, err)
	}
	return user
}

func selectStarsByUser(userID int, database *pg.DB) []Star {
	var stars []Star
	err := database.Model(&stars).
		Where("user_id=?", userID).
		Select()
	if err != nil {
		fmt.Printf("Failed to Select stars with userID: %v : %v\n", userID, err)
	}
	return stars
}

func selectStarByID(id int, database *pg.DB) Star {
	star := Star{ID: id}
	err := database.Select(&star)
	if err != nil {
		fmt.Printf("Failed to Select star with ID: %v : %v\n", id, err)
	}
	return star
}

func updateUser(user User, database *pg.DB) User {
	model := database.Model(&user)
	if user.Token != "" {
		model.Column("token")
	}
	_, err := model.WherePK().Update()
	if err != nil {
		fmt.Printf("Failed to Update user with ID: %v : %v\n", user.ID, err)
	}
	return user
}


func updateStar(star Star, database *pg.DB) Star {
	model := database.Model(&star)
	if star.Name != "" {
		model.Column("name")
	}
	if star.UserID != 0 {
		model.Column("user_id")
	}
	if star.Img != "" {
		model.Column("img")
	}
	if star.Link != "" {
		model.Column("link")
	}
	_, err := model.WherePK().Update()
	if err != nil {
		fmt.Printf("Failed to Update star with ID: %v : %v\n", star.ID, err)
	}
	return star
}


func deleteUser(user User, database *pg.DB) User {
	err := database.Delete(&user)
	if err != nil {
		fmt.Printf("Failed to Delete user with ID: %v : %v\n", user.ID, err)
	}
	return user
}


func deleteStar(star Star, database *pg.DB) Star {
	err := database.Delete(&star)
	if err != nil {
		fmt.Printf("Failed to Delete star with ID: %v : %v\n", star.ID, err)
	}
	return star
}

