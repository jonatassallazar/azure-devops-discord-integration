package models

import (
	"fmt"
	"time"
)

// ConvertToDiscordPayloadPR converts an Azure DevOps pull request webhook payload
// into a Discord webhook payload format.
//
// The method creates a Discord embed with:
//   - Author information from the pull request creator
//   - Title and URL linking to the pull request
//   - Description showing the repository name
//   - Color based on the provided color parameter
//   - Fields containing pull request title and reviewer information
//   - Footer showing the merge status
//
// Parameters:
//   - title: The title to display in the Discord embed
//   - color: The color code (int32) for the embed's left border
//
// Returns:
//   - DiscordPayload: A formatted Discord webhook payload ready to be sent
func (a *AzureRequest) ConvertToDiscordPayloadPR(title string, color int32) DiscordPayload {
	body := DiscordPayload{
		Username:  "Azure Pull Request",
		AvatarUrl: "",
		Content:   "",
		Embeds: []Embeds{
			{
				Author: Author{
					Name:    a.Resource.CreatedBy.DisplayName,
					Url:     a.Resource.CreatedBy.Url,
					IconUrl: a.Resource.CreatedBy.ImageUrl,
				},
				Title:       title,
				Url:         fmt.Sprintf("%s/pullrequest/%d", a.Resource.Repository.RemoteUrl, a.Resource.PullRequestId),
				Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
				Color:       color,
				Fields:      a.getFields(),
				Footer: Footer{
					Text: fmt.Sprintf("Status do Commit: %s", a.getMergeStatusText()),
				},
			},
		},
	}

	return body
}

// ConvertToDiscordPayloadPipeline converts an Azure DevOps pipeline webhook payload
// into a Discord webhook payload format.
//
// The method creates a Discord embed with:
//   - Author information from the person who requested the pipeline
//   - Title and URL linking to the pipeline build
//   - Description showing the repository name
//   - Color based on the provided color parameter
//   - Fields containing pipeline name, repository name, branch, and commit message
//
// Parameters:
//   - title: The title to display in the Discord embed
//   - color: The color code (int32) for the embed's left border
//   - repository: The Azure repository information fetched from the API
//
// Returns:
//   - DiscordPayload: A formatted Discord webhook payload ready to be sent
func (a *AzureRequest) ConvertToDiscordPayloadPipeline(title string, color int32, repository AzureRepository) DiscordPayload {
	body := DiscordPayload{
		Username:  "Azure Pipeline Build",
		AvatarUrl: "",
		Content:   "",
		Embeds: []Embeds{
			{
				Author: Author{
					Name:    a.Resource.RequestedFor.DisplayName,
					Url:     a.Resource.RequestedFor.Url,
					IconUrl: a.Resource.RequestedFor.ImageUrl,
				},
				Title:       title,
				Url:         a.Resource.Links.Web.Href,
				Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
				Color:       color,
				Fields: []Field{
					{Name: "Pipeline", Value: a.Resource.Definition.Name},
					{Name: "TriggeredBy", Value: repository.Name},
					{Name: "Branch", Value: a.Resource.TriggerInfo.SourceBranch},
					{Name: "Merge", Value: a.Resource.TriggerInfo.Message},
				},
			},
		},
	}

	return body
}

// ConvertToDiscordPayloadRelease converts an Azure DevOps release webhook payload
// into a Discord webhook payload format.
//
// The method creates a Discord embed with:
//   - Author information from the person who requested the deployment
//   - Title and URL linking to the release definition
//   - Description showing the release name
//   - Color based on the provided color parameter
//   - Fields containing release number, environment description, start time, and completion time
//
// Parameters:
//   - title: The title to display in the Discord embed
//   - color: The color code (int32) for the embed's left border
//
// Returns:
//   - DiscordPayload: A formatted Discord webhook payload ready to be sent
func (a *AzureRequest) ConvertToDiscordPayloadRelease(title string, color int32) DiscordPayload {
	body := DiscordPayload{
		Username:  "Azure Release Build",
		AvatarUrl: "",
		Content:   "",
		Embeds: []Embeds{
			{
				Author: Author{
					Name:    a.Resource.Deployment.RequestedFor.DisplayName,
					Url:     a.Resource.Deployment.RequestedFor.Url,
					IconUrl: a.Resource.Deployment.RequestedFor.ImageUrl,
				},
				Title:       title,
				Url:         a.Resource.Environment.ReleaseDefinition.Links.Web.Href,
				Description: fmt.Sprintf("Release: %s", a.Resource.Environment.ReleaseDefinition.Name),
				Color:       color,
				Fields: []Field{
					{Name: "Release Nº", Value: a.Resource.Environment.Release.Name},
					{Name: "Description", Value: a.Resource.Environment.Name},
					{Name: "Iniciado em", Value: a.Resource.Deployment.StartedOn.UTC().Format(time.RFC1123)},
					{Name: "Finalizado em", Value: a.Resource.Deployment.CompletedOn.UTC().Format(time.RFC1123)},
				},
			},
		},
	}

	return body
}

// getFields generates a list of Discord embed fields from the pull request data.
//
// The method creates fields containing:
//   - The pull request title
//   - Each reviewer's name and their vote status (as inline fields)
//
// Returns:
//   - []Field: A slice of Field structures ready to be used in a Discord embed
func (a *AzureRequest) getFields() []Field {
	fields := []Field{{Name: "Título", Value: a.Resource.Title}}

	for _, i := range a.Resource.Reviewers {
		fields = append(fields, Field{Name: i.DisplayName, Value: i.getVoteText(), Inline: true})
	}

	return fields
}

// getVoteText converts a reviewer's vote numeric value into a human-readable text string.
//
// Vote value mappings:
//   - 10: "Aprovado" (Approved)
//   - 5: "Aprovado com Sugestões" (Approved with Suggestions)
//   - 0: "Sem Voto" (No Vote)
//   - -5: "Aguardando o Autor" (Waiting for Author)
//   - -10: "Rejeitado" (Rejected)
//   - Any other value: Empty string
//
// Returns:
//   - string: A text representation of the vote status
func (r *Reviewers) getVoteText() string {
	switch r.Vote {
	case 10:
		return "Aprovado"
	case 5:
		return "Aprovado com Sugestões"
	case 0:
		return "Sem Voto"
	case -5:
		return "Aguardando o Autor"
	case -10:
		return "Rejeitado"
	default:
		return ""
	}
}

// getMergeStatusText converts a pull request merge status into a human-readable text string.
//
// Merge status mappings:
//   - "succeeded": "Sem conflito" (No conflicts)
//   - "conflicts": "Com conflito" (Has conflicts)
//   - "queued": "Aguardando" (Waiting)
//   - "rejectedByPolicy": "Rejeitado pelas regras" (Rejected by policy)
//   - "failure": "Com erros" (Has errors)
//   - Any other status: Empty string
//
// Returns:
//   - string: A text representation of the merge status
func (a *AzureRequest) getMergeStatusText() string {
	switch a.Resource.MergeStatus {
	case "succeeded":
		return "Sem conflito"
	case "conflicts":
		return "Com conflito"
	case "queued":
		return "Aguardando"
	case "rejectedByPolicy":
		return "Rejeitado pelas regras"
	case "failure":
		return "Com erros"
	default:
		return ""
	}
}
