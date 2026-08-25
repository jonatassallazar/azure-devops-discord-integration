package azuredevops

import (
	"fmt"
	"time"

	"azuredevops-notify/internal/notify"
)

func (a *AzureRequest) toPRMessage(title string, level notify.Level) notify.Message {
	return notify.Message{
		Source: "Azure Pull Request",
		Author: notify.Author{
			Name:    a.Resource.CreatedBy.DisplayName,
			URL:     a.Resource.CreatedBy.Url,
			IconURL: a.Resource.CreatedBy.ImageUrl,
		},
		Title:       title,
		URL:         fmt.Sprintf("%s/pullrequest/%d", a.Resource.Repository.RemoteUrl, a.Resource.PullRequestId),
		Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
		Level:       level,
		Fields:      a.getFields(),
		FooterText:  fmt.Sprintf("Status do Commit: %s", a.getMergeStatusText()),
	}
}

func (a *AzureRequest) toPipelineMessage(title string, level notify.Level, repository AzureRepository) notify.Message {
	return notify.Message{
		Source: "Azure Pipeline Build",
		Author: notify.Author{
			Name:    a.Resource.RequestedFor.DisplayName,
			URL:     a.Resource.RequestedFor.Url,
			IconURL: a.Resource.RequestedFor.ImageUrl,
		},
		Title:       title,
		URL:         a.Resource.Links.Web.Href,
		Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
		Level:       level,
		Fields: []notify.Field{
			{Name: "Pipeline", Value: a.Resource.Definition.Name},
			{Name: "TriggeredBy", Value: repository.Name},
			{Name: "Branch", Value: a.Resource.TriggerInfo.SourceBranch},
			{Name: "Merge", Value: a.Resource.TriggerInfo.Message},
		},
	}
}

func (a *AzureRequest) toReleaseMessage(title string, level notify.Level) notify.Message {
	return notify.Message{
		Source: "Azure Release Build",
		Author: notify.Author{
			Name:    a.Resource.Deployment.RequestedFor.DisplayName,
			URL:     a.Resource.Deployment.RequestedFor.Url,
			IconURL: a.Resource.Deployment.RequestedFor.ImageUrl,
		},
		Title:       title,
		URL:         a.Resource.Environment.ReleaseDefinition.Links.Web.Href,
		Description: fmt.Sprintf("Release: %s", a.Resource.Environment.ReleaseDefinition.Name),
		Level:       level,
		Fields: []notify.Field{
			{Name: "Release Nº", Value: a.Resource.Environment.Release.Name},
			{Name: "Description", Value: a.Resource.Environment.Name},
			{Name: "Iniciado em", Value: a.Resource.Deployment.StartedOn.UTC().Format(time.RFC1123)},
			{Name: "Finalizado em", Value: a.Resource.Deployment.CompletedOn.UTC().Format(time.RFC1123)},
		},
	}
}

func (a *AzureRequest) getFields() []notify.Field {
	fields := []notify.Field{{Name: "Título", Value: a.Resource.Title}}

	for _, r := range a.Resource.Reviewers {
		fields = append(fields, notify.Field{Name: r.DisplayName, Value: r.getVoteText(), Inline: true})
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
