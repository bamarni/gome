package main

import (
	"encoding/json"
	//"fmt"
	forecast "github.com/mlbright/forecast/v2"
	"log"
	"net/http"
	"os"
	//"runtime"
)

const (
	// Berlin
	lat  = "52.5167"
	long = "13.4"
)

var key = os.Getenv("FORECAST_API_KEY")

type Weather struct {
	Timezone    string
	Summary     string
	Icon        string
	Temperature int
}

// rounds a float (cf. https://github.com/golang/go/issues/4594)
func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}

	return int(val + 0.5)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	f, err := forecast.Get(key, lat, long, "now", forecast.CA)
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

func init() {
	if key == "" {
		log.Fatal("$FORECAST_API_KEY is not set")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/weather.json", weatherHandler)

	mux.Handle("/", http.FileServer(http.Dir("/web")))

	log.Fatal(http.ListenAndServe(":80", mux))
}
