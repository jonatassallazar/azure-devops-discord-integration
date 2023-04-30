package controllers

import (
	models "discord-azure-integration/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestProcessVotesExpected struct {
	color int32
	title string
}

type TestProcessVotesMock struct {
	AzureResponse models.AzureRequest
	Expected      TestProcessVotesExpected
}

func TestProcessVotes(t *testing.T) {
	var mockAzureResponse = []TestProcessVotesMock{
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: 0},
						{Vote: 0},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: 0,
				title: "",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: -10},
						{Vote: 0},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.RED,
				title: "Reprovado",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: -10},
						{Vote: 10},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.RED,
				title: "Reprovado",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: 10},
						{Vote: 10},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.GREEN,
				title: "Aprovado",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: 0},
						{Vote: 10},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.GREEN,
				title: "Aprovado",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: -5},
						{Vote: 10},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.ORANGE,
				title: "Aguardando Autor",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: -5},
						{Vote: -5},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.ORANGE,
				title: "Aguardando Autor",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Reviewers: []models.Reviewers{
						{Vote: -5},
						{Vote: -10},
					},
				},
			},
			Expected: TestProcessVotesExpected{
				color: models.RED,
				title: "Reprovado",
			},
		},
	}

	for _, a := range mockAzureResponse {
		var pr = PullRequestController{
			Response: &a.AzureResponse,
		}

		color, title := pr.processVotes()

		assert.Equal(t, a.Expected.color, color)
		assert.Equal(t, a.Expected.title, title)
	}

}
