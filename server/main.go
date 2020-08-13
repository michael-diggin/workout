package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-diggin/workout/pkg/database"
	"github.com/michael-diggin/workout/pkg/routing"
)

func main() {
	var port *string
	port = flag.String("p", ":8010", "port for server to listen on")
	flag.Parse()

	db, err := database.Open(os.Getenv("DBNAME"))
	if err != nil {
		log.Fatal(err)
	}
	database.EnsureTableExists(db)
	dbService := database.NewDBService(db)
	a := routing.App{}
	a.SetUp(dbService)
	fmt.Println("Running the server!")
	a.Run(*port)
}

// GetLocalIP returns the non loopback local IP of the host
// Used for debugging in docker/wsl
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
