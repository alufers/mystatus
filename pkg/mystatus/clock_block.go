package mystatus

import (
	"time"
)

type clockBlock struct {
	format string
}

func (cb *clockBlock) Render() barBlockData {
	return barBlockData{
		Name:     "clock",
		Instance: "local",
		FullText: time.Now().Format(cb.format),
	}
}