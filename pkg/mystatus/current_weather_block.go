package mystatus

import (
	"encoding/json"
	"fmt"
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
				cwb.LastCheckError = fmt.Errorf("failed to parser JSON response: %w", decodeErr)
			} else {
				cwb.LastCheckError = nil
				cwb.LastResponse = data
			}
		}
	}
	var text string
	var color string
	if cwb.LastCheckError == nil {
		color = "#ffffff"
	} else {
		color = "#aaaaaa"
	}
	if cwb.LastResponse != nil {
		text = cwb.LastResponse.CurrentCondition[0].TempC + "°C"
	}
	return barBlockData{
		Name:     "weather",
		Instance: "weather_" + cwb.Location,
		Color:    color,
		FullText: text,
	}
}
