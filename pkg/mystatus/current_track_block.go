package mystatus

import (
	"bufio"
	"fmt"
	"html"
	"log"
	"os/exec"
	"regexp"
	"sync"

	"github.com/alufers/mystatus/pkg/mousekeys"
	"go.i3wm.org/i3/v4"
)

type currentTrackBlock struct {
	currentTrackInfoMutex sync.RWMutex
	currentTrackInfo      map[string]string
}

func (ctb *currentTrackBlock) Init() {
	go ctb.StatusProcessRoutine()
	go ctb.MetadataProcessRoutine()
}

func (ctb *currentTrackBlock) StatusProcessRoutine() {
	cmd := exec.Command("playerctl", "status", "-F")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print("failed to run status process routine", err)
		return
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		ctb.currentTrackInfoMutex.Lock()
		ctb.currentTrackInfo["status"] = line
		ctb.currentTrackInfoMutex.Unlock()
		forceRerender <- nil
	}
}

var metadataRegex = regexp.MustCompile(`^[A-Za-z0-9_]*\ [A-za-z0-9_]*:([A-Za-z0-9_]*)\ + (.*)$`)

func (ctb *currentTrackBlock) MetadataProcessRoutine() {
	cmd := exec.Command("playerctl", "metadata", "-F")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print("failed to run metadata process routine", err)
		return
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		ctb.currentTrackInfoMutex.Lock()
		matches := metadataRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			ctb.currentTrackInfo[matches[1]] = matches[2]
		}
		ctb.currentTrackInfoMutex.Unlock()
		forceRerender <- nil
	}
}

func (ctb *currentTrackBlock) Render() barBlockData {
	ctb.currentTrackInfoMutex.RLock()
	defer ctb.currentTrackInfoMutex.RUnlock()

	statusIcon := ""
	color := "lightgray"
	switch ctb.currentTrackInfo["status"] {
	case "Playing":
		statusIcon = "⏸️"
		color = "white"
	case "Paused":
		statusIcon = "▶️"
		color = "lightgray"
	}
	return barBlockData{
		Block:    ctb,
		Name:     "current_track",
		Instance: "master",
		Markup:   "pango",
		FullText: fmt.Sprintf(`<span foreground="%s">%s %s - %s</span>`, color, statusIcon, html.EscapeString(ctb.currentTrackInfo["artist"]), html.EscapeString(ctb.currentTrackInfo["title"])),
	}
}

func (ctb *currentTrackBlock) HandleEvent(ie *InputEvent) {
	switch ie.Button {
	case mousekeys.Left:
		cmd := exec.Command("playerctl", "play-pause")
		cmd.Run()
	case mousekeys.Right:
		i3.RunCommand("[class=Spotify] focus")
	}
}
