package mystatus

import (
	"fmt"
	"github.com/distatus/battery"
)

type batteryBlock struct {
}

func (cb *batteryBlock) Render() barBlockData {
	fullText := "no battteries"
	batteries, _ := battery.GetAll()
	for _, battery := range batteries {
		state := battery.State.String()
		if state == "Unknown" {
			state = "Not charging"
		}
		fullText = fmt.Sprintf("🔋<span size=\"6000\">%s</span> %.0f%%", state, (battery.Current/battery.Full)*100.0)
	}
	return barBlockData{
		Name:     "battery",
		FullText: fullText,
		Markup:   "pango",
	}
}
