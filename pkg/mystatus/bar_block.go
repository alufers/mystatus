package mystatus

type barBlockData struct {
	Name     string `json:"name,omitempty"`
	Instance string `json:"instance,omitempty"`

	FullText  string `json:"full_text"`
	ShortText string `json:"short_text,omitempty"`

	Color      string `json:"color,omitempty"`
	Background string `json:"background,omitempty"`
	Border     string `json:"border,omitempty"`

	MinWidth  uint   `json:"min_width,omitempty"`
	Align     string `json:"align,omitempty"`
	Separator bool   `json:"separator,omitempty"`
	Urgent    bool   `json:"urgent,omitempty"`

	Markup string `json:"markup,omitempty"`
}

type barBlock interface {
	Render() barBlockData
}
