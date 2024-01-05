package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct{
	Location struct{
		Name string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct{
		TempC float64 `json:"temp_c"`

		Condition struct{
			Text string `json:"text"`
		} `json:"condition"`

	} `json:"current"`

	Forecast struct{
		Forecastday []struct{
			Hour []struct{
				TimeEpoch int64 `json:"time_epoch"`
				Temperature float64 `json:"temp_c"`

				Condition struct{
					Text string `json:"text"`
				} `json:"condition"`

				ChanceOfRain float64 `json:"chance_of_rain"`

			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`

}

func main() {
	var city string
	// city = "Hanoi"
	
	// if len(os.Args) >= 2 {
	// 	city = os.Args[1]
	// }


	if len(os.Args) < 2 {
		panic("Missing Location Argument!")
	} else{
		city = os.Args[1]
	}


	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=433f89aaaac44b2585183516231609&q=" + city +"&days=1&aqi=yes&alerts=yes")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API is not available")
	}

	var body, err1 = io.ReadAll(res.Body)
	if err1 != nil {
		panic(err)
	}

	var weather Weather
	err2 := json.Unmarshal(body, &weather)
	if err2 != nil {
		panic(err)
	}
	
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}

		var message = fmt.Sprintf(
			"Date: %s - Temp: %.0fC, %.0f%% rain, Condition: %s\n",
			date.Format("2006-01-02 15:04:05"),
			hour.Temperature,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 50 {
			color.Green(message)
		} else {
			color.Red(message)
		}
	}



}
	

