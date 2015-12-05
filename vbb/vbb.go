package vbb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type JourneyDetailRef struct {
	Ref string `json:"ref"`
}

type Product struct {
	Name         string `json:"name"`
	Num          string `json:"num"`
	Line         string `json:"line"`
	CatOut       string `json:"catOut"`
	CatIn        string `json:"catIn"`
	CatCode      string `json:"catCode"`
	CatOutS      string `json:"catOutS"`
	CatOutL      string `json:"catOutL"`
	OperatorCode string `json:"operatorCode"`
	Operator     string `json:"operator"`
	Admin        string `json:"admin"`
}

type Departure struct {
	JourneyDetailRef JourneyDetailRef `json:"JourneyDetailRef"`
	Product          Product          `json:"Product"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	Stop             string           `json:"stop"`
	StopId           string           `json:"stopid"`
	StopExtId        string           `json:"stopExtId"`
	PrognosisType    string           `json:"prognosisType"`
	Time             string           `json:"time"`
	Date             string           `json:"date"`
	RtTime           string           `json:"rtTime"`
	RtDate           string           `json:"rtDate"`
	Direction        string           `json:"direction"`
	TrainNumber      string           `json:"trainNumber"`
	TrainCategory    string           `json:"trainCategory"`
}

type DepartureBoard struct {
	Departures []Departure `json:"Departure"`
}

func FromJSON(json_blob []byte) (*DepartureBoard, error) {
	var d DepartureBoard

	err := json.Unmarshal(json_blob, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetResponse(key string) (*http.Response, error) {

	parameters := url.Values{}

	parameters.Add("accessId", key)
	parameters.Add("id", "A=1@O=Bersarinplatz (Berlin)@X=13452563@Y=52519341@U=86@L=009120516@B=1@V=3.9,@p=1449160158@")
	parameters.Add("direction", "9120004")
	parameters.Add("format", "json")

	Url := url.URL{
		Scheme:   "http",
		Host:     "demo.hafas.de",
		Path:     "/openapi/vbb-proxy/departureBoard",
		RawQuery: parameters.Encode(),
	}

	res, err := http.Get(Url.String())

	if err != nil {
		return res, err
	}

	return res, nil
}

func Get(key string) (*DepartureBoard, error) {
	res, err := GetResponse(key)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	dep, err := FromJSON(body)

	if err != nil {
		return nil, err
	}

	return dep, nil
}
