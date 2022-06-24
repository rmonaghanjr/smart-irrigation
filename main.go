package main

import (
	"database/sql"
	"fmt"
	"os"
	"smart-irrigation/m/v2/server"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("opened database")

	err = rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("opened gpio connection")

	defer rpio.Close()

	pin := rpio.Pin(18)

	server.Start(db, &pin, false)

}