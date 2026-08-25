// Package googlechat implements notify.Sink for Google Chat incoming
// webhooks using the Cards v2 message format.
package googlechat

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

// Google Chat cards have no per-card background color, unlike Discord
// embeds, so the level is rendered as a colored-circle indicator in the
// header title instead.
var levelIndicator = map[notify.Level]string{
	notify.LevelPending:   "🟡",
	notify.LevelSuccess:   "🟢",
	notify.LevelFailure:   "🔴",
	notify.LevelWarning:   "🟠",
	notify.LevelCompleted: "🔵",
	notify.LevelUnmapped:  "⚪",
}

type payload struct {
	CardsV2 []cardWrapper `json:"cardsV2"`
}

type cardWrapper struct {
	CardID string `json:"cardId"`
	Card   card   `json:"card"`
}

type card struct {
	Header   cardHeader `json:"header"`
	Sections []section  `json:"sections"`
}

type cardHeader struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type section struct {
	Widgets []widget `json:"widgets"`
}

type widget struct {
	DecoratedText *decoratedText `json:"decoratedText,omitempty"`
	TextParagraph *textParagraph `json:"textParagraph,omitempty"`
	ButtonList    *buttonList    `json:"buttonList,omitempty"`
}

type decoratedText struct {
	TopLabel string `json:"topLabel,omitempty"`
	Text     string `json:"text"`
}

type textParagraph struct {
	Text string `json:"text"`
}

type buttonList struct {
	Buttons []button `json:"buttons"`
}

type button struct {
	Text    string  `json:"text"`
	OnClick onClick `json:"onClick"`
}

type onClick struct {
	OpenLink openLink `json:"openLink"`
}

type openLink struct {
	URL string `json:"url"`
}

// Sink delivers notify.Messages to a Google Chat incoming webhook.
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
	var widgets []widget

	if msg.Description != "" {
		widgets = append(widgets, widget{TextParagraph: &textParagraph{Text: msg.Description}})
	}

	for _, f := range msg.Fields {
		widgets = append(widgets, widget{DecoratedText: &decoratedText{TopLabel: f.Name, Text: f.Value}})
	}

	if msg.FooterText != "" {
		widgets = append(widgets, widget{TextParagraph: &textParagraph{Text: msg.FooterText}})
	}

	if msg.URL != "" {
		widgets = append(widgets, widget{
			ButtonList: &buttonList{
				Buttons: []button{
					{Text: "Abrir no Azure DevOps", OnClick: onClick{OpenLink: openLink{URL: msg.URL}}},
				},
			},
		})
	}

	title := msg.Title
	if indicator, ok := levelIndicator[msg.Level]; ok {
		title = fmt.Sprintf("%s %s", indicator, title)
	}

	body := payload{
		CardsV2: []cardWrapper{
			{
				CardID: "notification",
				Card: card{
					Header: cardHeader{
						Title:    title,
						Subtitle: msg.Author.Name,
						ImageURL: msg.Author.IconURL,
					},
					Sections: []section{{Widgets: widgets}},
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("googlechat: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.WebhookURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("googlechat: build request: %w", err)
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("googlechat: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("googlechat: webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
