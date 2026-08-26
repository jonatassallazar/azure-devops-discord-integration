package httpserver

import (
	"azuredevops-notify/internal/azuredevops"
	"azuredevops-notify/internal/notify"
	"azuredevops-notify/internal/notify/discord"
	"azuredevops-notify/internal/notify/googlechat"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupRouter() *gin.Engine {
	if s.Config.GinMode != "" {
		gin.SetMode(s.Config.GinMode)
	}

	r := gin.Default()

	prHandler := azuredevops.PullRequestHandler{
		Dispatcher: buildDispatcher(s.Config.Discord.PRWebhookURL, s.Config.GoogleChat.PRWebhookURL),
	}
	pipelineHandler := azuredevops.PipelineHandler{
		Dispatcher: buildDispatcher(s.Config.Discord.PipelineWebhookURL, s.Config.GoogleChat.PipelineWebhookURL),
		Azure:      s.Config.Azure,
	}
	releaseHandler := azuredevops.ReleaseHandler{
		Dispatcher: buildDispatcher(s.Config.Discord.ReleaseWebhookURL, s.Config.GoogleChat.ReleaseWebhookURL),
	}

	r.GET(RouteHealth, healthCheck)

	r.POST(azuredevops.RouteCreatedPR, prHandler.CreatedPR)
	r.POST(azuredevops.RouteReviewedPR, prHandler.ReviewedPR)
	r.POST(azuredevops.RouteStatusUpdatedPR, prHandler.StatusUpdatedPR)
	r.POST(azuredevops.RoutePipeline, pipelineHandler.PipelineStatusReport)
	r.POST(azuredevops.RouteRelease, releaseHandler.ReleaseStatusReport)

	return r
}

// buildDispatcher wires up a sink for every configured webhook URL. An
// empty URL disables that sink for the event category without erroring,
// so a Discord-only deployment keeps working with no config changes.
func buildDispatcher(discordURL, googleChatURL string) *notify.Dispatcher {
	var sinks []notify.Sink

	if discordURL != "" {
		sinks = append(sinks, discord.New(discordURL))
	}
	if googleChatURL != "" {
		sinks = append(sinks, googlechat.New(googleChatURL))
	}

	return &notify.Dispatcher{Sinks: sinks}
}
