package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct{
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct{
	Name string `json:"name"`
	Main struct{
		Kelvin float64 `json:"temp"`	//temperature is in kelvin
	}`json:"main"`
}

func load(filename string) (apiConfigData, error){
	bytes, err := ioutil.ReadFile(filename)
	if err != nil{
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil{
		return apiConfigData{}, err
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello from go\n"))
}

func query(city string) (weatherData, error){
	apiConfig, err := load(".apiConfig")
	if err != nil{
		return weatherData{}, err
	}

	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID="+ apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil{
		return weatherData{}, err
	}

	defer res.Body.Close()
	var d weatherData
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil{
		return weatherData{}, err
	}

	return d, nil
}

func main(){
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}


//run using go run main.go command
