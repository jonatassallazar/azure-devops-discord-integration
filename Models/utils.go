package models

import "fmt"

func (a *AzureRequest) ConvertToDiscordPayload(title string, approved int8, reproved bool) DiscordPayload {
	body := DiscordPayload{
		Username:  "Azure Pull Request",
		AvatarUrl: "https://pbs.twimg.com/profile_images/1145617831905681408/XNKktHjN_400x400.png",
		Content:   a.Resource.Title,
		Embeds: []Embeds{
			{
				Author: Author{
					Name:    a.Resource.CreatedBy.DisplayName,
					Url:     a.Resource.CreatedBy.Url,
					IconUrl: a.Resource.CreatedBy.ImageUrl,
				},
				Title:       title,
				Url:         a.Resource.Url,
				Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
				Color:       16705372,
				Fields:      a.getFields(),
			},
		},
	}

	return body
}

func (a *AzureRequest) getFields() []Field {
	fields := []Field{{Name: a.Resource.Status, Value: a.DetailedMessage.Text}}

	for _, i := range a.Resource.Reviewers {
		fields = append(fields, Field{Name: i.DisplayName, Value: i.getVoteText()})
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
