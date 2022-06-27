package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/stianeikeland/go-rpio/v4"
)

type Config struct {
	StationName string `json:"station_name"`
	Version     string `json:"version"`
	Configured  bool   `json:"configured"`
}

type Router struct {
	DB      *sql.DB
	Pin     *rpio.Pin
	Channel chan string
	IsOn    bool
}

type WateringData struct {
	Logs []WaterLogRow `json:"logs"`
}

type WaterLogRow struct {
	Id     int     `json:"id"`
	Date   int     `json:"date"`
	Amount float64 `json:"amount"`
	Time   float64 `json:"time"`
}

func NewRouter(db *sql.DB, pin *rpio.Pin, channel chan string) *Router {
	return &Router{
		DB:      db,
		Pin:     pin,
		Channel: channel,
		IsOn:    false,
	}
}

func Start(db *sql.DB, pin *rpio.Pin, channel chan string, production bool) {
	router := NewRouter(db, pin, channel)

	http.HandleFunc("/data", router.GetWateringData)
	http.HandleFunc("/toggle", router.TogglePumpPower)
	http.HandleFunc("/stats", router.StationStatistics)
	http.HandleFunc("/configure", router.ConfigureServer)

	var addr string = ":8080"

	if production {
		addr = ":80"
	}

	router.Channel <- "out:opened http server"
	http.ListenAndServe(addr, nil)
}

/*
/api/data

Reports the amount of water today as well as all the times that it has dispensed water previously
*/

func (router *Router) GetWateringData(w http.ResponseWriter, req *http.Request) {
	statement, err := router.DB.Prepare("select * from water_log")
	if err != nil {
		fmt.Println(err)
	}

	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		fmt.Print(err)
	}

	res := WateringData{
		Logs: make([]WaterLogRow, 0),
	}

	for rows.Next() {
		log := WaterLogRow{}
		err := rows.Scan(&log.Id, &log.Date, &log.Amount, &log.Time)
		if err != nil {
			fmt.Println(err)
		}

		res.Logs = append(res.Logs, log)
	}

	jsonResult, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(jsonResult)
}

func (router *Router) TogglePumpPower(w http.ResponseWriter, req *http.Request) {
	router.IsOn = !router.IsOn

	if router.IsOn {
		router.Channel <- "pin:on"
		w.Write([]byte("pin:on"))
	} else {
		router.Channel <- "pin:off"
		w.Write([]byte("pin:off"))
	}
}

func (router *Router) StationStatistics(w http.ResponseWriter, req *http.Request) {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	value, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(value)
}

func (router *Router) ConfigureServer(w http.ResponseWriter, req *http.Request) {
	config := Config{}

	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	value, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(value, &config)

	interval, iErr := strconv.Atoi(req.URL.Query().Get("interval"))
	if iErr != nil {
		fmt.Println(iErr)
	}
	smartWater, sErr := strconv.Atoi(req.URL.Query().Get("smart_water"))
	if sErr != nil {
		fmt.Println(sErr)
	}

	if !config.Configured {
		statement, err := router.DB.Prepare("insert into timing (interval, smart_water) values (?, ?)")
		if err != nil {
			fmt.Println(err)
		}

		defer statement.Close()

		statement.Exec(interval, smartWater)
	} else {
		statement, err := router.DB.Prepare("update timing set interval=?, smart_water=? where id=1")
		if err != nil {
			fmt.Println(err)
		}

		defer statement.Close()

		statement.Exec(interval, smartWater)
	}

	config.Configured = true

	raw, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile("./config.json", raw, 0777)

	router.Channel <- "START"

	w.Write([]byte("true"))

}
