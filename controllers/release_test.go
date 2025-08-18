package controllers

import (
	models "discord-azure-integration/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestProcessReleaseStatusExpected struct {
	color int32
	title string
}

type TestProcessReleaseStatusMock struct {
	AzureResponse models.AzureRequest
	Expected      TestProcessReleaseStatusExpected
}

func TestProcessReleaseStatus(t *testing.T) {
	var mockAzureResponse = []TestProcessReleaseStatusMock{
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Deployment: models.Deployment{
						DeploymentStatus: "succeeded",
					},
				},
			},
			Expected: TestProcessReleaseStatusExpected{
				color: models.GREEN,
				title: "Concluída",
			},
		},
		{
			AzureResponse: models.AzureRequest{

				Resource: models.Resource{
					Deployment: models.Deployment{
						DeploymentStatus: "failed",
					},
				},
			},
			Expected: TestProcessReleaseStatusExpected{
				color: models.RED,
				title: "Falhada",
			},
		}, {
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Deployment: models.Deployment{
						DeploymentStatus: "stopped",
					},
				},
			},
			Expected: TestProcessReleaseStatusExpected{
				color: models.ORANGE,
				title: "Interrompida",
			},
		}, {
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Deployment: models.Deployment{
						DeploymentStatus: "any",
					},
				},
			},
			Expected: TestProcessReleaseStatusExpected{
				color: models.WHITE,
				title: "[Status não mapeado: any]",
			},
		},
	}

	for _, a := range mockAzureResponse {
		var pr = ReleaseController{
			Response: &a.AzureResponse,
		}

		color, title := pr.processReleaseStatus()

		assert.Equal(t, a.Expected.color, color)
		assert.Equal(t, a.Expected.title, title)
	}

}
