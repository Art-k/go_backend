package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Board struct {
	gorm.Model
	Mac         string
	Name        string
	Description string
}

func main() {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&Board{})

	// Create
	db.Create(&Board{Mac: "00:00:00:00:01", Name: "test board"})

	// Read
	var board Board
	db.First(&board, 1)                           // find product with id 1
	db.First(&board, "Mac = ?", "00:00:00:00:01") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&board).Update("Description", "test description")

	// Delete - delete product
	db.Delete(&board)
}
