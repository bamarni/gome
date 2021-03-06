package main

import (
	"encoding/json"
	"github.com/bamarni/vbb"
	forecast "github.com/mlbright/forecast/v2"
	gojira "github.com/plouc/go-jira-client"
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
	Vbb         *vbb.Vbb
	Jira        *gojira.Jira
	HttpDir     string
	HttpAddr    string
}

type AppHandler struct {
	Config *AppConfig
	Handle func(c *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error)
}

type Weather struct {
	Timezone    string
	Summary     string
	Icon        string
	Temperature int
	Unit        string
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

func weatherHandler(c *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	f, err := forecast.Get(c.ForecastKey, c.Lat, c.Long, "now", forecast.CA)
	if err != nil {
		return nil, err
	}

	weather := Weather{
		Timezone:    f.Timezone,
		Summary:     f.Daily.Data[0].Summary,
		Icon:        f.Daily.Data[0].Icon,
		Temperature: Round(f.Daily.Data[0].TemperatureMax),
		Unit:        "°C",
	}

	return weather, nil
}

func vbbHandler(config *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	depRequest := &vbb.Departure{
		StopId:    "A=1@O=Bersarinplatz (Berlin)@X=13452563@Y=52519341@U=86@L=009120516@B=1@V=3.9,@p=1449160158@",
		Direction: "9120004",
	}

	depBoard, err := config.Vbb.GetDepartureBoard(depRequest)

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

func jiraHandler(config *AppConfig, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	issues := config.Jira.IssuesAssignedTo(config.Jira.Auth.Login, 3, 0)

	return issues, nil
}

func main() {
	conf := &AppConfig{
		Location:    "Europe/Berlin",
		Lat:         "52.5167",
		Long:        "13.4",
		ForecastKey: os.Getenv("FORECAST_API_KEY"),
		Vbb: &vbb.Vbb{
			os.Getenv("VBB_API_KEY"),
		},
		Jira: gojira.NewJira(
			os.Getenv("JIRA_BASE_URL"),
			os.Getenv("JIRA_API_PATH"),
			"",
			&gojira.Auth{
				os.Getenv("JIRA_LOGIN"),
				os.Getenv("JIRA_PASSWORD"),
			},
		),
		HttpDir:  "/web",
		HttpAddr: ":80",
	}

	http.Handle("/weather.json", AppHandler{conf, weatherHandler})
	http.Handle("/vbb.json", AppHandler{conf, vbbHandler})
	http.Handle("/jira.json", AppHandler{conf, jiraHandler})
	http.Handle("/", http.FileServer(http.Dir(conf.HttpDir)))

	log.Printf("Listening at %s...", conf.HttpAddr)
	log.Fatal(http.ListenAndServe(conf.HttpAddr, nil))
}
