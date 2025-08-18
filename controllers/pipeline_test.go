package controllers

import (
	models "discord-azure-integration/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestProcessStatusExpected struct {
	color int32
	title string
}

type TestProcessStatusMock struct {
	AzureResponse models.AzureRequest
	Expected      TestProcessStatusExpected
}

func TestProcessStatus(t *testing.T) {
	var mockAzureResponse = []TestProcessStatusMock{
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Result: "succeeded",
				},
			},
			Expected: TestProcessStatusExpected{
				color: models.GREEN,
				title: "Concluída",
			},
		},
		{
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Result: "failed",
				},
			},
			Expected: TestProcessStatusExpected{
				color: models.RED,
				title: "Falhada",
			},
		}, {
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Result: "stopped",
				},
			},
			Expected: TestProcessStatusExpected{
				color: models.ORANGE,
				title: "Interrompida",
			},
		}, {
			AzureResponse: models.AzureRequest{
				Resource: models.Resource{
					Result: "any",
				},
			},
			Expected: TestProcessStatusExpected{
				color: models.WHITE,
				title: "[Status não mapeado: any]",
			},
		},
	}

	for _, a := range mockAzureResponse {
		var pr = PipelineController{
			Response: &a.AzureResponse,
		}

		color, title := pr.processStatus()

		assert.Equal(t, a.Expected.color, color)
		assert.Equal(t, a.Expected.title, title)
	}

}
