package azuredevops

import (
	"testing"

	"azuredevops-notify/internal/notify"

	"github.com/stretchr/testify/assert"
)

func TestReleaseHandlerProcessReleaseStatus(t *testing.T) {
	tests := []struct {
		name             string
		deploymentStatus string
		wantLevel        notify.Level
		wantTitle        string
	}{
		{name: "succeeded", deploymentStatus: "succeeded", wantLevel: notify.LevelSuccess, wantTitle: "Concluída"},
		{name: "failed", deploymentStatus: "failed", wantLevel: notify.LevelFailure, wantTitle: "Falhada"},
		{name: "stopped", deploymentStatus: "stopped", wantLevel: notify.LevelWarning, wantTitle: "Interrompida"},
		{name: "unmapped", deploymentStatus: "any", wantLevel: notify.LevelUnmapped, wantTitle: "[Status não mapeado: any]"},
	}

	h := ReleaseHandler{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AzureRequest{Resource: Resource{Deployment: Deployment{DeploymentStatus: tt.deploymentStatus}}}
			level, title := h.processReleaseStatus(&req)

			assert.Equal(t, tt.wantLevel, level)
			assert.Equal(t, tt.wantTitle, title)
		})
	}
}
