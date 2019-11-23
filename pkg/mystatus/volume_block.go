package mystatus

import "os/exec"

import "strings"

type volumeBlock struct {
}

func (cb *volumeBlock) Render() barBlockData {

	cmd := exec.Command("bash", "-c", "amixer sget Master | grep 'Right:' | awk -F'[][]' '{ print $2 }'")
	stdout, err := cmd.Output()
	var volStr string
	if err != nil {
		volStr = "error"
	} else {
		volStr = (string(stdout))
	}
	volStr = strings.ReplaceAll(volStr, "\n", "")
	return barBlockData{
		Name:     "volume",
		Instance: "master",
		Markup:   "pango",
		FullText: `🔊 ` + volStr + ``,
	}
}
