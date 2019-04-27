package mystatus

import (
	"encoding/json"
	"fmt"
	"time"
)

var blocks = []barBlock{
	&ipAddressBlock{
		interfaceNames: []string{"enp3s0", "wlp3s0", "enx503eaa70e8da"},
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
	&batteryBlock{},
	&clockBlock{
		format: "2006-01-02 15:04:05",
	},
}

func printJSON(data interface{}) {
	jsonOut, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(jsonOut[:]))
}

func printHeader() {
	printJSON(map[string]interface{}{
		"version": 1,
	})
	fmt.Print("\n")
}

func printBar() {
	combinedBlockData := []barBlockData{}
	for _, b := range blocks {
		data := b.Render()
		combinedBlockData = append(combinedBlockData, data)
	}
	printJSON(combinedBlockData)
	fmt.Print("\n")
	fmt.Print(",")
}

/*
Run is the real entry func of the program
*/
func Run() {
	printHeader()
	fmt.Println("[")
	for {
		printBar()
		time.Sleep(time.Second)
	}
}
