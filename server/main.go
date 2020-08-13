package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-diggin/workout/pkg/database"
	"github.com/michael-diggin/workout/pkg/routing"
)

func main() {
	GetLocalIP()

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

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting IP")
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
			}
		}
	}
}
