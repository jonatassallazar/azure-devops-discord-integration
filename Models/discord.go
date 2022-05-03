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
