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
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatarUrl"`
	Content   string  `json:"content"`
	Embeds    []embed `json:"embeds"`
}

type embed struct {
	Author      author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int32   `json:"color"`
	Fields      []field `json:"fields"`
	Thumbnail   linkURL `json:"thumbnail"`
	Image       linkURL `json:"image"`
	Footer      footer  `json:"footer"`
}

type author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
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
	IconURL string `json:"icon_url"`
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

	body := payload{
		Username: msg.Source,
		Embeds: []embed{
			{
				Author:      author{Name: msg.Author.Name, URL: msg.Author.URL, IconURL: msg.Author.IconURL},
				Title:       msg.Title,
				URL:         msg.URL,
				Description: msg.Description,
				Color:       levelColors[msg.Level],
				Fields:      fields,
				Thumbnail:   linkURL{URL: msg.ThumbnailURL},
				Footer:      footer{Text: msg.FooterText},
			},
		},
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
