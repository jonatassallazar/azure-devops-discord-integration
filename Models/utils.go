package models

import "fmt"

func (a *AzureRequest) ConvertToDiscordPayload(approved int8, reproved bool) DiscordPayload {
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
				Title:       "Pull Request Criado",
				Url:         a.Resource.Url,
				Description: fmt.Sprintf("Projeto %s", a.Resource.Repository.Name),
				Color:       16705372,
				Fields: []Field{
					{Name: "PR Aberto", Value: a.DetailedMessage.Text},
				},
			},
		},
	}

	return body
}
