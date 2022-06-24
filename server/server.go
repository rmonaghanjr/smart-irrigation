package server

import (
	"database/sql"
	"net/http"
	"smart-irrigation/m/v2/api"

	"github.com/stianeikeland/go-rpio/v4"
)

type Router struct {
	DB      *sql.DB
	Pin     *rpio.Pin
	Channel chan string
	IsOn    bool
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

	http.HandleFunc("/api/data", router.GetWateringData)
	http.HandleFunc("/toggle", router.TogglePumpPower)

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

}

func (router *Router) TogglePumpPower(w http.ResponseWriter, req *http.Request) {
	api.TogglePower(router.Pin)
	router.IsOn = !router.IsOn

	if router.IsOn {
		router.Channel <- "pin:on"
		w.Write([]byte("pin:on"))
	} else {
		router.Channel <- "pin:off"
		w.Write([]byte("pin:off"))
	}
}
