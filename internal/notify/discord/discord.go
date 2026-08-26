// Package discord implements notify.Sink for Discord incoming webhooks.
package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"azuredevops-notify/internal/notify"
)

const contentTypeJSON = "application/json"

var levelColors = map[notify.Level]int32{
	notify.LevelPending:   16705372, // #FEE75C
	notify.LevelSuccess:   5763713,  // #57F287
	notify.LevelFailure:   15548997, // #ED4245
	notify.LevelWarning:   15105570, // #E67E22
	notify.LevelCompleted: 5793266,  // #5865F2
	notify.LevelUnmapped:  16777215, // #FFFFFF
}

type payload struct {
	Username string `json:"username,omitempty"`
	// Discord's webhook field is avatar_url, not avatarUrl - spelled the
	// old way it was silently ignored as an unknown field.
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []embed `json:"embeds"`
}

// The optional objects are pointers so an unset one is omitted entirely:
// an "image"/"thumbnail" carrying an empty url is a broken image for
// Discord to render, not an absent one.
type embed struct {
	Author      *author  `json:"author,omitempty"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Color       int32    `json:"color"`
	Fields      []field  `json:"fields"`
	Thumbnail   *linkURL `json:"thumbnail,omitempty"`
	Image       *linkURL `json:"image,omitempty"`
	Footer      *footer  `json:"footer,omitempty"`
}

type author struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type linkURL struct {
	URL string `json:"url"`
}

type footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// Sink delivers notify.Messages to a Discord incoming webhook.
type Sink struct {
	WebhookURL string
	HTTPClient *http.Client
}

// New builds a Sink with a sane default HTTP client timeout.
func New(webhookURL string) *Sink {
	return &Sink{
		WebhookURL: webhookURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Sink) Send(ctx context.Context, msg notify.Message) error {
	fields := make([]field, len(msg.Fields))
	for i, f := range msg.Fields {
		fields[i] = field{Name: f.Name, Value: f.Value, Inline: f.Inline}
	}

	item := embed{
		Title:       msg.Title,
		URL:         msg.URL,
		Description: msg.Description,
		Color:       levelColors[msg.Level],
		Fields:      fields,
	}

	if msg.Author.Name != "" {
		item.Author = &author{Name: msg.Author.Name, URL: msg.Author.URL, IconURL: msg.Author.IconURL}
	}
	if msg.ThumbnailURL != "" {
		item.Thumbnail = &linkURL{URL: msg.ThumbnailURL}
	}
	if msg.FooterText != "" {
		item.Footer = &footer{Text: msg.FooterText}
	}

	body := payload{
		Username: msg.Source,
		Embeds:   []embed{item},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.WebhookURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("discord: build request: %w", err)
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("discord: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord: webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
