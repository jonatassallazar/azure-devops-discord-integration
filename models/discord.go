package models

type Author struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type URL struct {
	Url string `json:"url"`
}

type Footer struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

type Embeds struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
	Color       int32   `json:"color"`
	Fields      []Field `json:"fields"`
	Thumbnail   URL     `json:"thumbnail"`
	Image       URL     `json:"image"`
	Footer      Footer  `json:"footer"`
}

type DiscordPayload struct {
	Username  string   `json:"username"`
	AvatarUrl string   `json:"avatarUrl"`
	Content   string   `json:"content"`
	Embeds    []Embeds `json:"embeds"`
}

const (
	ORANGE  int32 = 15105570 // #E67E22
	RED     int32 = 15548997 // #ED4245
	GRAY    int32 = 10070709 // #99AAB5
	YELLOW  int32 = 16705372 // #FEE75C
	BLURPLE int32 = 5793266  // #5865F2
	WHITE   int32 = 16777215 // #FFFFFF
	GREEN   int32 = 5763713  // #57F287
	BLACK   int32 = 2303786  // #23272A
)
