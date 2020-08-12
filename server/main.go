package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-diggin/workout/pkg/database"
	"github.com/michael-diggin/workout/pkg/routing"
)

func main() {
	db, err := database.Open("./exercise.db")
	if err != nil {
		log.Fatal(err)
	}
	database.EnsureTableExists(db)
	dbService := database.NewDBService(db)
	a := routing.App{}
	a.SetUp(dbService)
	fmt.Println("Running the server!")
	a.Run(":8010")
}
