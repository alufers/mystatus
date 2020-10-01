package mystatus

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alufers/mystatus/pkg/mousekeys"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type resourceChartBlock struct {
	history         []float64
	dangerTreshhold float64
	chartType       string
	baseColor       string
	dangerColor     string

	resourcePercent float64
	resourceError   error
}

func (ccb *resourceChartBlock) Tick() {

	if ccb.chartType == "cpu" {
		cpuData, err := cpu.Percent(0, false)
		ccb.resourcePercent = cpuData[0]
		ccb.resourceError = err
	} else if ccb.chartType == "mem" {
		virtStat, err := mem.VirtualMemory()
		ccb.resourcePercent = virtStat.UsedPercent
		ccb.resourceError = err
	} else {
		panic(errors.New("Invalid chartType"))
	}

	if ccb.resourceError == nil {
		if len(ccb.history) < 30 {
			for i := 0; i < 30; i++ {
				ccb.history = append(ccb.history, 0)
			}
		}
		ccb.history = append(ccb.history[1:], ccb.resourcePercent)
	}

}

func (ccb *resourceChartBlock) Render() barBlockData {

	var text string
	if ccb.resourceError != nil {
		text = "Error: " + ccb.resourceError.Error()
	} else {
		text += fmt.Sprintf(`<span rise="3000" size="10000">%02d</span>`, int(ccb.resourcePercent))

		for i := 0; i < len(ccb.history); i++ {
			resultingColor := ccb.baseColor
			if ccb.history[i] > ccb.dangerTreshhold {
				resultingColor = ccb.dangerColor
			}
			text += fmt.Sprintf(`<span rise="%s"  size="5000" color="%s">•</span>`, strconv.Itoa(int(ccb.history[i]*100)), resultingColor)

		}
	}

	return barBlockData{
		Block:    ccb,
		Name:     "resource_chart",
		Instance: ccb.chartType,
		FullText: text,
		Markup:   "pango",
	}
}

func (ccb *resourceChartBlock) HandleEvent(ie *InputEvent) {
	switch ie.Button {
	case mousekeys.Left:
		cmd := exec.Command("kitty", strings.Split("--class i3_float_mouse --override remember_window_size=false --override initial_window_height=650 --override initial_window_width=950 htop", " ")...)
		cmd.Run()
	}
}

//
