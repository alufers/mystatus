package mystatus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

//lastData stores the data returned from rendering the blocks so that the inputEventHandlers knows where to dispatch the data
var lastData = []barBlockData{}
var lastDataMutex = sync.RWMutex{}
var forceRerender = make(chan interface{}, 10)
var lastTick time.Time = time.Now()

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
		log.Printf("input ev: %#v", event)
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

func performTick() {
	lastTick = time.Now()
	for _, b := range blocks {
		if tickable, ok := b.(tickableBarBlock); ok {
			tickable.Tick()
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
	os.MkdirAll(os.ExpandEnv("$HOME/.local/log"), 0777)
	f, err := os.OpenFile(os.ExpandEnv("$HOME/.local/log/mystatus.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	go inputEventsScanner()
	for _, blk := range blocks {
		if initable, ok := blk.(initableBarBlock); ok {
			initable.Init()
		}
	}
	printHeader()
	fmt.Println("[")
	performTick()
	for {
		if time.Now().Sub(lastTick) >= time.Second {
			performTick()
		}
		printBar()
		select {
		case <-time.After(time.Second):
		case <-forceRerender:
		}

	}
}
