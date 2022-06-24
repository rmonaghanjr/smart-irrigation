package server

import (
	"database/sql"
	"net/http"
	"smart-irrigation/m/v2/api"

	"github.com/stianeikeland/go-rpio/v4"
)

type Router struct {
	DB  *sql.DB
	Pin *rpio.Pin
}

func NewRouter(db *sql.DB, pin *rpio.Pin) *Router {
	return &Router{
		DB:  db,
		Pin: pin,
	}
}

func Start(db *sql.DB, pin *rpio.Pin, production bool) {
	router := NewRouter(db, pin)

	http.HandleFunc("/api/data", router.GetWateringData)

	var addr string = ":8080"

	if production {
		addr = ":80"
	}

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
}
