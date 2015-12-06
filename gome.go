package main

import (
	"encoding/json"
	"github.com/bamarni/gome/vbb"
	forecast "github.com/mlbright/forecast/v2"
	"log"
	"net/http"
	"os"
	"time"
)

type AppConfig struct {
	Location    string
	Lat         string
	Long        string
	ForecastKey string
	VbbKey      string
}

type AppHandler struct {
	Config *AppConfig
	Handle func(e *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error)
}

type Weather struct {
	Timezone    string
	Summary     string
	Icon        string
	Temperature int
}

type Departure struct {
	Time  string
	Delay int
}

// rounds a float (cf. https://github.com/golang/go/issues/4594)
func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}

	return int(val + 0.5)
}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	data, err := h.Handle(h.Config, w, r)
	if err == nil {
		err = json.NewEncoder(w).Encode(data)
	}
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
}

func weatherHandler(config *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	f, err := forecast.Get(config.ForecastKey, config.Lat, config.Long, "now", forecast.CA)
	if err != nil {
		return nil, err
	}

	weather := Weather{
		f.Timezone,
		f.Currently.Summary,
		f.Currently.Icon,
		Round(f.Currently.Temperature),
	}

	return weather, nil
}

func vbbHandler(config *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	depBoard, err := vbb.Get(config.VbbKey)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation(config.Location)
	timeLayout := "2006-01-02 15:04:05"

	var deps []Departure
	for _, dep := range depBoard.Departures[0:3] {
		Time, _ := time.ParseInLocation(timeLayout, dep.Date+" "+dep.Time, loc)
		realTime, _ := time.ParseInLocation(timeLayout, dep.RtDate+" "+dep.RtTime, loc)
		delay := realTime.Sub(Time).Minutes()
		dep := Departure{Time.Format("15:04"), int(delay)}
		deps = append(deps, dep)
	}

	return deps, nil
}

func main() {
	appConfig := &AppConfig{
		Location:    "Europe/Berlin",
		Lat:         "52.5167",
		Long:        "13.4",
		ForecastKey: os.Getenv("FORECAST_API_KEY"),
		VbbKey:      os.Getenv("VBB_API_KEY"),
	}

	if appConfig.ForecastKey == "" || appConfig.VbbKey == "" {
		log.Fatal("$FORECAST_API_KEY and $VBB_API_KEY must be set")
	}

	mux := http.NewServeMux()

	mux.Handle("/weather.json", AppHandler{appConfig, weatherHandler})
	mux.Handle("/vbb.json", AppHandler{appConfig, vbbHandler})
	mux.Handle("/", http.FileServer(http.Dir("/web")))

	log.Fatal(http.ListenAndServe(":80", mux))
}
