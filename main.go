package main

import (
	"database/sql"
	"fmt"
	"os"
	"smart-irrigation/m/v2/server"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	const ML_PER_S float64 = 36.6666666667

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("opened database")

	devMode := false
	_, err = os.Stat("/dev/mem")
	if os.IsNotExist(err) {
		devMode = true
	}

	pinChannel := make(chan string)

	var pin rpio.Pin

	if !devMode {
		pin := rpio.Pin(18)
		pin.Output()
	}

	go server.Start(db, &pin, pinChannel, false)

	pinOn := false
	var timeStart int64

	for value := range pinChannel {
		if value[0:3] == "out" {
			fmt.Println(value[4:])
		} else if value[0:3] == "pin" {
			if value[4:] == "on" {
				pinOn = true
				timeStart = time.Now().UnixMilli()
			} else if value[4:] == "off" && pinOn {
				pinOn = false
				diff := (time.Now().UnixMilli() - timeStart)

				statement, err := db.Prepare("insert into water_log (date, amount, time) values (?, ?, ?)")
				if err != nil {
					fmt.Println(err)
				}

				statement.Exec(time.Now().UnixMilli(), ML_PER_S*float64(diff)/1000.0, float64(diff)/1000.0)
				statement.Close()
			}
		}
	}

}
