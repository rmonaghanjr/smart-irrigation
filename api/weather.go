package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type WeatherResponse struct {
	Lat            float64  `json:"lat"`
	Lon            float64  `json:"lon"`
	Timezone       string   `json:"timezone"`
	TimezoneOffset int      `json:"timezone_offset"`
	Current        Current  `json:"current"`
	Minutely       []Minute `json:"minutely"`
	Alerts         []Alert  `json:"alerts"`
}

type Current struct {
	Dt         int            `json:"dt"`
	Sunrise    int            `json:"sunrise"`
	Sunset     int            `json:"sunset"`
	Temp       float64        `json:"temp"`
	FeelsLike  float64        `json:"feels_like"`
	Pressure   int            `json:"pressure"`
	Humidity   int            `json:"humidity"`
	DewPoint   float64        `json:"dew_point"`
	Uvi        int            `json:"uvi"`
	Clouds     int            `json:"clouds"`
	Visibility int            `json:"visibility"`
	WindSpeed  float64        `json:"wind_speed"`
	WindDeg    int            `json:"wind_deg"`
	Weather    []WeatherState `json:"weather"`
}

type WeatherState struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Alert struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int      `json:"start"`
	End         int      `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Minute struct {
	Dt            int `json:"dt"`
	Precipitation int `json:"precipitation"`
}

func FetchWeatherData(url string) *WeatherResponse {
	var weatherResponse WeatherResponse
	resp, err := http.Get(url)
	if err != nil {
		panic("http request error")
	}

	jsonBody, bErr := ioutil.ReadAll(resp.Body)
	if bErr != nil {
		panic("body read error")
	}

	json.Unmarshal(jsonBody, &weatherResponse)

	return &weatherResponse
}
