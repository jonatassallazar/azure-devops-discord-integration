package azuredevops

import (
	"testing"

	"azuredevops-notify/internal/notify"

	"github.com/stretchr/testify/assert"
)

func TestPipelineHandlerProcessStatus(t *testing.T) {
	tests := []struct {
		name      string
		result    string
		wantLevel notify.Level
		wantTitle string
	}{
		{name: "succeeded", result: "succeeded", wantLevel: notify.LevelSuccess, wantTitle: "Concluída"},
		{name: "failed", result: "failed", wantLevel: notify.LevelFailure, wantTitle: "Falhada"},
		{name: "stopped", result: "stopped", wantLevel: notify.LevelWarning, wantTitle: "Interrompida"},
		{name: "unmapped", result: "any", wantLevel: notify.LevelUnmapped, wantTitle: "[Status não mapeado: any]"},
	}

	h := PipelineHandler{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AzureRequest{Resource: Resource{Result: tt.result}}
			level, title := h.processStatus(&req)

			assert.Equal(t, tt.wantLevel, level)
			assert.Equal(t, tt.wantTitle, title)
		})
	}
}
