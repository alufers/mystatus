package mystatus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type currentWeatherBlock struct {
	Location        string        // the location to get the wiather from,
	RecheckInterval time.Duration // how oftn to check the weather
	LastCheck       *time.Time
	LastCheckError  error
	LastResponse    *WttrInResponse
}

func (cwb *currentWeatherBlock) Render() barBlockData {
	now := time.Now()
	var lastCheckTime time.Time
	if cwb.LastCheck == nil {
		lastCheckTime = time.Unix(0, 0)
	} else {
		lastCheckTime = *cwb.LastCheck
	}
	if (cwb.LastCheckError != nil && now.Sub(lastCheckTime) > time.Second*30) || now.Sub(lastCheckTime) > cwb.RecheckInterval {

		resp, err := http.Get(fmt.Sprintf("http://wttr.in/%s?format=j1", cwb.Location))
		cwb.LastCheck = &now
		if err == nil {
			decoder := json.NewDecoder(resp.Body)
			var data *WttrInResponse
			if decodeErr := decoder.Decode(&data); decodeErr != nil {
				log.Print(decodeErr)
				cwb.LastCheckError = fmt.Errorf("failed to parser JSON response: %w", decodeErr)
				cwb.LastResponse = nil
			} else {
				cwb.LastCheckError = nil
				cwb.LastResponse = data
			}
		} else {

		}
	}
	var text string
	var color string
	if cwb.LastCheckError == nil && cwb.LastResponse != nil {
		color = "#ffffff"
	} else if cwb.LastCheckError == nil {
		color = "#ffaaaa"
	} else {
		color = "#aaaaaa"
	}
	if cwb.LastResponse != nil {
		var emoji string
		if name, ok := WeatherCodesToNames[cwb.LastResponse.CurrentCondition[0].WeatherCode]; ok {
			if emojiValue, ok := WeatherNamesToEmoji[name]; ok {
				emoji = emojiValue
			}
		}
		text = fmt.Sprintf("%s %s", emoji, cwb.LastResponse.CurrentCondition[0].TempC+"°C")
	}
	if cwb.LastCheckError == nil && cwb.LastResponse == nil {
		text = "E:" + cwb.LastCheckError.Error()
	}
	return barBlockData{
		Name:     "weather",
		Instance: "weather_" + cwb.Location,
		Color:    color,
		FullText: text,
	}
}
