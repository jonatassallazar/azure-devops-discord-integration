package models

// Author represents the author information in a Discord embed.
// This typically appears at the top of an embed with the author's name, URL, and icon.
type Author struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

// Field represents a field in a Discord embed.
// Fields are displayed as key-value pairs and can be displayed inline (side by side)
// or as full-width blocks.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// URL represents a simple URL structure used for thumbnail and image fields in Discord embeds.
type URL struct {
	Url string `json:"url"`
}

// Footer represents the footer information in a Discord embed.
// The footer appears at the bottom of an embed and can include text and an icon.
type Footer struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

// Embeds represents a Discord embed object.
//
// Embeds are rich content blocks that can include titles, descriptions, fields,
// images, thumbnails, colors, and more. This structure follows the Discord webhook
// embed format specification.
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

// DiscordPayload represents the complete payload structure for a Discord webhook message.
//
// This structure is sent to Discord webhook URLs to post messages. It can contain
// plain text content and/or rich embeds. The username and avatar can be customized
// to override the default webhook appearance.
type DiscordPayload struct {
	Username  string   `json:"username"`
	AvatarUrl string   `json:"avatarUrl"`
	Content   string   `json:"content"`
	Embeds    []Embeds `json:"embeds"`
}

// Discord embed color constants.
//
// These constants represent hexadecimal color codes converted to decimal integers
// for use in Discord embed color fields. Discord embeds use a colored bar on the
// left side of the embed, and these colors are commonly used to indicate status
// or category of the notification.
const (
	// ORANGE represents an orange color (#E67E22).
	// Commonly used for warnings or interrupted/stopped operations.
	ORANGE int32 = 15105570 // #E67E22
	
	// RED represents a red color (#ED4245).
	// Commonly used for errors, failures, or rejected states.
	RED int32 = 15548997 // #ED4245
	
	// GRAY represents a gray color (#99AAB5).
	// Commonly used for neutral or informational messages.
	GRAY int32 = 10070709 // #99AAB5
	
	// YELLOW represents a yellow color (#FEE75C).
	// Commonly used for pending states or new/created items.
	YELLOW int32 = 16705372 // #FEE75C
	
	// BLURPLE represents Discord's brand color (#5865F2).
	// Commonly used for completed operations or Discord-related notifications.
	BLURPLE int32 = 5793266 // #5865F2
	
	// WHITE represents a white color (#FFFFFF).
	// Commonly used for unmapped or unknown statuses.
	WHITE int32 = 16777215 // #FFFFFF
	
	// GREEN represents a green color (#57F287).
	// Commonly used for success, approved, or completed states.
	GREEN int32 = 5763713 // #57F287
	
	// BLACK represents a dark gray/black color (#23272A).
	// Commonly used for dark-themed embeds or neutral states.
	BLACK int32 = 2303786 // #23272A
)
