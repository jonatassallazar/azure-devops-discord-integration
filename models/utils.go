package models

import (
	"fmt"
	"time"
)

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

func (a *AzureRequest) ConvertToDiscordPayloadPipeline(title string, color int32) DiscordPayload {
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
					{Name: "Branch", Value: a.Resource.TriggerInfo.SourceBranch},
					{Name: "Merge", Value: a.Resource.TriggerInfo.Message},
				},
			},
		},
	}

	return body
}

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

func (a *AzureRequest) getFields() []Field {
	fields := []Field{{Name: "Título", Value: a.Resource.Title}}

	for _, i := range a.Resource.Reviewers {
		fields = append(fields, Field{Name: i.DisplayName, Value: i.getVoteText(), Inline: true})
	}

	return fields
}

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
