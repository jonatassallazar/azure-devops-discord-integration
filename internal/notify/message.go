// Package notify defines a vendor-neutral notification model and the Sink
// interface that outbound integrations (Discord, Google Chat, ...) implement.
package notify

// Level is a vendor-neutral notification status. The zero value pairs with
// an empty title to mean "nothing to notify" - callers should not send a
// Message with the zero Level.
type Level int

const (
	_ Level = iota
	LevelPending
	LevelSuccess
	LevelFailure
	LevelWarning
	LevelCompleted
	LevelUnmapped
)

type Author struct {
	Name    string
	URL     string
	IconURL string
}

type Field struct {
	Name   string
	Value  string
	Inline bool
}

// Message is a vendor-neutral notification built by a source integration and
// delivered by one or more Sinks.
type Message struct {
	Source       string
	Author       Author
	Title        string
	URL          string
	Description  string
	Level        Level
	Fields       []Field
	ThumbnailURL string
	FooterText   string
}
