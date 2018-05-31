package mystatus

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type resourceChartBlock struct {
	history         []float64
	dangerTreshhold float64
	chartType       string
	baseColor       string
	dangerColor     string
}

func (ccb *resourceChartBlock) Render() barBlockData {

	var resourcePercent float64
	var resourceError error
	if ccb.chartType == "cpu" {
		cpuData, err := cpu.Percent(0, false)
		resourcePercent = cpuData[0]
		resourceError = err
	} else if ccb.chartType == "mem" {
		virtStat, err := mem.VirtualMemory()
		resourcePercent = virtStat.UsedPercent
		resourceError = err
	} else {
		panic(errors.New("Invalid chartType"))
	}

	var text string
	if resourceError != nil {
		text = "CPU error: " + resourceError.Error()
	} else {

		if len(ccb.history) < 10 {
			for i := 0; i < 10; i++ {
				ccb.history = append(ccb.history, 0)
			}
		}

		ccb.history = append(ccb.history[1:], resourcePercent)

		text += `<span rise="5000" size="5000">` + strconv.Itoa(int(resourcePercent)) + `</span>`

		for i := 0; i < len(ccb.history); i++ {
			resultingColor := ccb.baseColor
			if ccb.history[i] > ccb.dangerTreshhold {
				resultingColor = ccb.dangerColor
			}
			text += fmt.Sprintf(`<span rise="%s"  size="5000" color="%s">•</span>`, strconv.Itoa(int(ccb.history[i]*100)), resultingColor)

		}
	}

	return barBlockData{
		FullText: text,
		Markup:   "pango",
	}
}
