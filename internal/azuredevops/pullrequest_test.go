package azuredevops

import (
	"testing"

	"azuredevops-notify/internal/notify"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestHandlerProcessVotes(t *testing.T) {
	tests := []struct {
		name      string
		reviewers []Reviewers
		wantLevel notify.Level
		wantTitle string
	}{
		{
			name:      "no votes",
			reviewers: []Reviewers{{Vote: 0}, {Vote: 0}},
			wantLevel: 0,
			wantTitle: "",
		},
		{
			name:      "rejected alone",
			reviewers: []Reviewers{{Vote: -10}, {Vote: 0}},
			wantLevel: notify.LevelFailure,
			wantTitle: "Reprovado",
		},
		{
			name:      "rejected wins over approved",
			reviewers: []Reviewers{{Vote: -10}, {Vote: 10}},
			wantLevel: notify.LevelFailure,
			wantTitle: "Reprovado",
		},
		{
			name:      "approved",
			reviewers: []Reviewers{{Vote: 10}, {Vote: 10}},
			wantLevel: notify.LevelSuccess,
			wantTitle: "Aprovado",
		},
		{
			name:      "approved with a no-vote",
			reviewers: []Reviewers{{Vote: 0}, {Vote: 10}},
			wantLevel: notify.LevelSuccess,
			wantTitle: "Aprovado",
		},
		{
			name:      "waiting wins over approved",
			reviewers: []Reviewers{{Vote: -5}, {Vote: 10}},
			wantLevel: notify.LevelWarning,
			wantTitle: "Aguardando Autor",
		},
		{
			name:      "waiting alone",
			reviewers: []Reviewers{{Vote: -5}, {Vote: -5}},
			wantLevel: notify.LevelWarning,
			wantTitle: "Aguardando Autor",
		},
		{
			name:      "rejected wins over waiting",
			reviewers: []Reviewers{{Vote: -5}, {Vote: -10}},
			wantLevel: notify.LevelFailure,
			wantTitle: "Reprovado",
		},
	}

	h := PullRequestHandler{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AzureRequest{Resource: Resource{Reviewers: tt.reviewers}}
			level, title := h.processVotes(&req)

			assert.Equal(t, tt.wantLevel, level)
			assert.Equal(t, tt.wantTitle, title)
		})
	}
}
