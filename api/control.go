package api

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func TriggerForSeconds(seconds int, pin *rpio.Pin) {
	if *pin != 0 {
		pin.High()
		time.Sleep(time.Duration(seconds) * time.Second)
		pin.Low()
	}
}

func TriggerForMillis(milliseconds int, pin *rpio.Pin) {
	if *pin != 0 {
		pin.High()
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		pin.Low()
	}
}

func TriggerForMinutes(minutes int, pin *rpio.Pin) {
	if *pin != 0 {
		pin.High()
		time.Sleep(time.Duration(minutes) * time.Minute)
		pin.Low()
	}
}

func TriggerOn(pin *rpio.Pin) {
	if *pin != 0 {
		pin.High()
	}
}

func TriggerOff(pin *rpio.Pin) {
	if *pin != 0 {
		pin.Low()
	}
}

func TogglePower(pin *rpio.Pin) {
	fmt.Println(*pin)
	if *pin != 0 {
		pin.Toggle()
	}
}
