package main

import (
	"encoding/json"
	forecast "github.com/mlbright/forecast/v2"
	"log"
	"net/http"
	"os"
	"github.com/bamarni/gome/vbb"
	"time"
)

const (
	location = "Europe/Berlin"
	lat      = "52.5167"
	long     = "13.4"
)

var (
	forecastKey = os.Getenv("FORECAST_API_KEY")
	vbbKey      = os.Getenv("VBB_API_KEY")
)

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

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	f, err := forecast.Get(forecastKey, lat, long, "now", forecast.CA)

	if err != nil {
		log.Fatal(err)
	}

	weather := Weather{
		f.Timezone,
		f.Currently.Summary,
		f.Currently.Icon,
		Round(f.Currently.Temperature),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(weather); err != nil {
		panic(err)
	}
}

func vbbHandler(w http.ResponseWriter, r *http.Request) {
	depBoard, err := vbb.Get(vbbKey)

	if err != nil {
		log.Fatal(err)
	}

	loc, _     := time.LoadLocation(location)
	timeLayout := "2006-01-02 15:04:05"

	var deps []Departure
	for _, dep := range depBoard.Departures[0:3] {
		Time, _     := time.ParseInLocation(timeLayout, dep.Date + " " + dep.Time, loc)
		realTime, _ := time.ParseInLocation(timeLayout, dep.RtDate + " " + dep.RtTime, loc)
		delay       := realTime.Sub(Time).Minutes()
		dep         := Departure{Time.Format("15:04"), int(delay)}
		deps         = append(deps, dep)
	}

	if err := json.NewEncoder(w).Encode(deps); err != nil {
		panic(err)
	}
}

func init() {
	if forecastKey == "" {
		log.Fatal("$FORECAST_API_KEY is not set")
	}
	if vbbKey == "" {
		log.Fatal("$VBB_API_KEY is not set")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/weather.json", weatherHandler)

	mux.HandleFunc("/vbb.json", vbbHandler)

	mux.Handle("/", http.FileServer(http.Dir("/web")))

	log.Fatal(http.ListenAndServe(":80", mux))
}
