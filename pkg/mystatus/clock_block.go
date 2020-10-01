package mystatus

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"go.i3wm.org/i3/v4"
)

type clockBlock struct {
	format string
}

func (cb *clockBlock) Render() barBlockData {
	return barBlockData{
		Block:    cb,
		Name:     "clock",
		Instance: "local",
		Markup:   "pango",

		FullText: time.Now().Format(cb.format) + " ",
	}
}

func (cb *clockBlock) HandleEvent(ie *InputEvent) {

	// open pavucontrol
	if ie.Button == 1 {

		cmd := exec.Command("yad", "--calendar") // YAD yet another dialog
		cmd.Start()
		go func() {
			var win *i3.Node
			for {
				time.Sleep(time.Millisecond * 1)
				t, err := i3.GetTree()
				if err != nil {
					log.Printf("failed to get tree: %f", err)
				}
				win = t.Root.FindChild(func(n *i3.Node) bool {

					return n.WindowProperties.Title == "YAD"
				})
				if win != nil {
					break
				}
			}
			i3.RunCommand(fmt.Sprintf("[con_id=\"%v\"] focus", win.ID))
			i3.RunCommand("floating enable")
			i3.RunCommand("move position cursor")
			i3.RunCommand("move up")
			i3.RunCommand("move up")
		}()
	}

}
