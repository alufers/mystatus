package mystatus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ^[A-Za-z]*\ +([A-Za-z\:]*)\ +(.*)

var blocks = []barBlock{
	&ipAddressBlock{
		interfaceNames: []string{"enp3s0", "wlp3s0", "enx503eaa70e8da", "enp0s20f0u1u4"},
	},
	&resourceChartBlock{
		dangerTreshhold: 70,
		chartType:       "cpu",
		baseColor:       "#3498db",
		dangerColor:     "#e74c3c",
	},
	&resourceChartBlock{
		dangerTreshhold: 70,
		chartType:       "mem",
		baseColor:       "#2ecc71",
		dangerColor:     "#e67e22",
	},
	&currentWeatherBlock{
		Location:        "Rybnik",
		RecheckInterval: time.Minute * 5,
	},
	&batteryBlock{},
	&volumeBlock{},
	&currentTrackBlock{
		currentTrackInfoMutex: sync.RWMutex{},
		currentTrackInfo:      map[string]string{},
	},
	&clockBlock{
		format: "2006-01-02 <b>15:04:05</b>",
	},
}

var lastData = []barBlockData{}
var lastDataMutex = sync.RWMutex{}
var forceRerender = make(chan interface{}, 10)

func printJSON(data interface{}) {
	jsonOut, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(jsonOut[:]))
}

func printHeader() {
	printJSON(map[string]interface{}{
		"version":      1,
		"click_events": true,
	})
	fmt.Print("\n")
}

func inputEventsScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lineWithoutTheComma := strings.TrimPrefix(line, ",")
		var event = &InputEvent{}
		err := json.Unmarshal([]byte(lineWithoutTheComma), &event)
		if err == nil {
			func() {
				lastDataMutex.RLock()
				defer lastDataMutex.RUnlock()
				for _, d := range lastData {
					if d.Name != "" && d.Instance != "" && d.Name == event.Name && d.Instance == event.Instance {
						if d.Block != nil {
							if handler, ok := d.Block.(EventHandlingBlock); ok {

								handler.HandleEvent(event)
							}
						}
					}
				}
			}()
			forceRerender <- nil
		}
	}
}

func printBar() {
	combinedBlockData := []barBlockData{}
	for _, b := range blocks {
		data := b.Render()
		combinedBlockData = append(combinedBlockData, data)
	}
	lastDataMutex.Lock()
	defer lastDataMutex.Unlock()
	lastData = combinedBlockData
	printJSON(combinedBlockData)
	fmt.Print("\n")
	fmt.Print(",")
}

/*
Run is the real entry func of the program
*/
func Run() {
	go inputEventsScanner()
	for _, blk := range blocks {
		if initable, ok := blk.(initableBarBlock); ok {
			initable.Init()
		}
	}
	printHeader()
	fmt.Println("[")
	for {
		printBar()
		select {
		case <-time.After(time.Second):
		case <-forceRerender:
		}

	}
}
