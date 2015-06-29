package main

import (
	
	 "net/http"
	 "io/ioutil"
	 "fmt"
	 "encoding/json"
	
	)

//date":"20150604","wcode":"04","weather":"雷阵雨","tmax":22.0,"tmin":17.0,"tavg":19.6}

func main() {

	getAqiSummary()
}

type WSummary struct {
	
	Date    string  `json:"date"`
	Wcode   string  `json:"wcode"`
	Weather string  `json:"weather"`
	Tmax    float64 `json:"tmax"`
	Tmin    float64 `json:"tmin"`
	Tavg 	float64 `json:"tavg"`
}

//{"week":6,"aqi":67.9,"date":"20150605","weekDay":"周五"}
type AqiSummary struct {
	Week    string  `json:"week"`
	Aqi     float64 `json:"aqi"`
	Date    string  `json:"date"`
	WeekDay string  `json:"weekday"`
}
type Area struct{
	AreaId string
	Province string
	County string
	City string
}

func getWsummary(areas... Area )[]WSummary{
	resp, err :=http.Get("http://dev.api.mlogcn.com:8000/api/weather/v1/aqi/summary/days/7/area/101010200.json")
	if err != nil {
	// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var sums []WSummary
	json.Unmarshal(body,&sums)
	fmt.Println(len(sums))
	fmt.Println(sums[0].Wcode)
	return nil
}

func getAqiSummary(areas... Area)[]AqiSummary{
	resp, err :=http.Get("http://dev.api.mlogcn.com:8000/api/weather/v1/aqi/summary/days/7/area/101010200.json")
	if err != nil {
	// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var sums []AqiSummary
	json.Unmarshal(body,&sums)
	fmt.Println(len(sums))
	fmt.Println(sums[0].Date)
	fmt.Println(sums[0].Aqi)
	
	return nil
}