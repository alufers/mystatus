package mystatus

type InputEvent struct {
	Button    int      `json:"button"`
	Modifiers []string `json:"modifiers"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	RelativeX int      `json:"relative x"`
	RelativeY int      `json:"relative y"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	Name      string   `json:"name"`
	Instance  string   `json:"instance"`
}

type EventHandlingBlock interface {
	HandleEvent(ie *InputEvent)
}
