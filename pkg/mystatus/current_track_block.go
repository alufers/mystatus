package mystatus

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sync"
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
	switch ctb.currentTrackInfo["status"] {
	case "Playing":
		statusIcon = "⏸️"
	case "Paused":
		statusIcon = "▶️"
	}
	return barBlockData{
		Block:    ctb,
		Name:     "current_track",
		Instance: "master",
		Markup:   "pango",
		FullText: fmt.Sprintf("%s %s - %s", statusIcon, ctb.currentTrackInfo["artist"], ctb.currentTrackInfo["title"]),
	}
}

func (ctb *currentTrackBlock) HandleEvent(ie *InputEvent) {
	if ie.Button == 1 {
		cmd := exec.Command("playerctl", "play-pause")
		cmd.Run()
	}
}
