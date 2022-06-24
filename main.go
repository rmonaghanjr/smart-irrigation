package main

import (
	"database/sql"
	"fmt"
	"os"
	"smart-irrigation/m/v2/server"
	"strconv"
	"time"

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

	pinChannel := make(chan string)

	pin := rpio.Pin(18)
	pin.Output()

	go server.Start(db, &pin, pinChannel, false)

	pinOn := false
	var timeStart int64

	for value := range pinChannel {
		if value[0:3] == "out" {
			fmt.Println(value[4:])
		} else if value[0:3] == "pin" {
			// time running detection
			if value[4:] == "on" {
				pinOn = true
				timeStart = time.Now().UnixMilli()
			} else if value[4:] == "off" && pinOn {
				pinOn = false
				diff := (time.Now().UnixMilli() - timeStart)

				fmt.Println("ran for " + strconv.Itoa(int(diff)))
			}
		}
	}

}
