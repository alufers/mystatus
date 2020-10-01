package mystatus

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/alufers/mystatus/pkg/mousekeys"
	"go.i3wm.org/i3/v4"
)

type volumeBlock struct {
}

type volumeStatus struct {
	volumePercent int
	unmuted       bool
}

func (vb *volumeBlock) GetVolumeStatus() (*volumeStatus, error) {
	cmd := exec.Command("amixer", "sget", "Master")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get volume status: %w", err)
	}
	commandOutput := string(stdout)
	outputRegexp := regexp.MustCompile(`\n*.Left:.+?\[([0-9]*)%\].+?\[([a-zA-z]*)\]`)
	matches := outputRegexp.FindStringSubmatch(commandOutput)[1:]
	volumeNum, _ := strconv.Atoi(matches[0])
	return &volumeStatus{
		volumePercent: volumeNum,
		unmuted:       matches[1] == "on",
	}, nil
}

func (vb *volumeBlock) Render() barBlockData {

	status, err := vb.GetVolumeStatus()
	var volStr string
	if err != nil {
		volStr = err.Error()
	} else {
		if status.unmuted {
			volStr += "🔊 "
		} else {
			volStr += "🔇 "
		}
		volStr += strconv.Itoa(status.volumePercent) + "%"
	}

	return barBlockData{
		Block:    vb,
		Name:     "volume",
		Instance: "master",
		Markup:   "pango",
		FullText: volStr,
	}
}

func (vb *volumeBlock) HandleEvent(ie *InputEvent) {
	switch ie.Button {
	case mousekeys.ScrollUp: // louden
		cmd := exec.Command("bash", "-c", "amixer -q sset Master 3%+")
		_, _ = cmd.Output()
	case mousekeys.ScrollDown: // lower volume
		cmd := exec.Command("bash", "-c", "amixer -q sset Master 3%-")
		_, _ = cmd.Output()
	case mousekeys.Middle: // mute/unmute
		cmd := exec.Command("bash", "-c", "amixer -q sset Master toggle")
		_, _ = cmd.Output()
	case mousekeys.Left: // open pavucontrol
		if ie.Button == 1 {
			cmd := exec.Command("pavucontrol")
			cmd.Start()
			go func() {
				recv := i3.Subscribe(i3.WindowEventType)
				defer recv.Close()
				for recv.Next() {
					ev := recv.Event().(*i3.WindowEvent)
					if ev.Change == "focus" && ev.Container.WindowProperties.Class != "Pavucontrol" {
						cmd.Process.Kill()
						return
					}

				}

			}()
		}
	}

}
