package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct { 
	OpenWeatherKey string `json:"OpenWeatherKey"`
}
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Сelsius float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var conf apiConfigData
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return apiConfigData{}, err
	}
	return conf, nil
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))
}
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=57&lon=-2.15&&units=metric&appid=" + apiConfig.OpenWeatherKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}
func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", nil)

}
