package azuredevops

const (
	RouteCreatedPR       = "/pull-request/created"
	RouteReviewedPR      = "/pull-request/review"
	RouteStatusUpdatedPR = "/pull-request/status"
	RoutePipeline        = "/pipeline/"
	RouteRelease         = "/release/"

	// RouteAvatar serves identity avatars fetched from Azure DevOps with
	// the configured PAT; see AvatarProxy for why the raw Azure URL cannot
	// be handed to the chat platforms directly.
	RouteAvatar = "/avatar/:ref"
)
