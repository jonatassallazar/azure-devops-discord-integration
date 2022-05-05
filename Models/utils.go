package models

import "fmt"

func (a *AzureRequest) ConvertToDiscordPayload(title string, color int32) DiscordPayload {
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
				Url:         fmt.Sprintf("%s/pullrequest/%d", a.Resource.Repository.RemoteUrl, a.CodeReviewId),
				Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
				Color:       color,
				Fields:      a.getFields(),
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
