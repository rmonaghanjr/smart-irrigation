package api

import (
	"database/sql"
	"fmt"
	"time"
)

func TimingScheduler(db *sql.DB, output chan string) {

	ticker := time.NewTicker(time.Hour)
	done := make(chan int)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			rows, err := db.Query("select * from timing limit 1")
			if err != nil {
				fmt.Println(err)
			}

			var interval int
			var smartWater int

			for rows.Next() {
				rows.Scan(&interval, &smartWater)
			}

			output <- "pin:on"
			time.Sleep(time.Duration(interval) * time.Minute)
			output <- "pin:off"
		}
	}
}