package main

import (
	"bytes"
	models "discord-azure-integration/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestCreateRequestRoute(t *testing.T) {
	router := setupRouter()

	fakeTime := time.Now()

	azureContent := models.AzureRequest{
		SubscriptionId: "2f5389a6-c75a-4048-8eba-33c0767d9ad6",
		NotificationId: 1094,
		ID:             "1c75aaf4-a8de-4748-835b-8bf7964ac0fa",
		EventType:      "git.pullrequest.created",
		PublisherId:    "tfs",
		Message: models.Text{
			Text: "Iassam da Silva de Souza created pull request 181109 (feature/cetei-6703/shared-playbooks-and-activities-on-create > release/v1.95) in empodera_app\r\nhttps://totvstfs.visualstudio.com/Empodera/_git/empodera_app/",
		},
		DetailedMessage: models.Text{
			Text: "Iassam da Silva de Souza created pull request 181109 (feature/cetei-6703/shared-playbooks-and-activities-on-create > release/v1.95) in empodera_app\r\nhttps://totvstfs.visualstudio.com/Empodera/_git/empodera_app/\r\n**Observação: A branch será ajustada após o deploy**\r\n",
		},
		Resource: models.Resource{
			Repository: models.Repository{
				ID:   "0f8a3f57-0bed-4451-86df-09d919e0d169",
				Name: "empodera_app",
				Url:  "https://totvstfs.visualstudio.com/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9/_apis/git/repositories/0f8a3f57-0bed-4451-86df-09d919e0d169",
				Project: models.Project{
					ID:             "a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9",
					Name:           "Empodera",
					Description:    "Repositorio do time do Empodera.",
					Url:            "https://totvstfs.visualstudio.com/_apis/projects/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9",
					State:          "wellFormed",
					Revision:       3762,
					Visibility:     "private",
					LastUpdateTime: fakeTime,
				},
				Size:            6476068,
				RemoteUrl:       "https://totvstfs.visualstudio.com/Empodera/_git/empodera_app",
				SshUrl:          "totvstfs@vs-ssh.visualstudio.com:v3/totvstfs/Empodera/empodera_app",
				WebUrl:          "https://totvstfs.visualstudio.com/Empodera/_git/empodera_app",
				IsDisabled:      false,
				IsInMaintenance: false,
			},
			PullRequestId: 181109,
			CodeReviewId:  181108,
			Status:        "active",
			CreatedBy: models.CreatedBy{
				DisplayName: "Iassam da Silva de Souza",
				Url:         "https://spsprodcus1.vssps.visualstudio.com/Aecd9a5a6-03de-436d-acc4-a73f4c80ca8f/_apis/Identities/a4dd11b6-b7bc-64e7-8a8e-1617179bef68",
				ID:          "a4dd11b6-b7bc-64e7-8a8e-1617179bef68",
				UniqueName:  "BRSSI0002@totvspartners.com.br",
				ImageUrl:    "https://totvstfs.visualstudio.com/_api/_common/identityImage?id=a4dd11b6-b7bc-64e7-8a8e-1617179bef68",
			},
			CreationDate:  fakeTime,
			Title:         "feature/cetei-6703/shared-playbooks-and-activities-on-create > release/v1.95",
			Description:   "**Observação: A branch será ajustada após o deploy**",
			SourceRefName: "refs/heads/feature/cetei-6703/shared-playbooks-and-activities-on-create",
			TargetRefName: "refs/heads/release/v1.95",
			MergeStatus:   "succeeded",
			IsDraft:       false,
			MergeID:       "9634d63f-a0ab-49de-8767-52e80ac18240",
			LastMergeSourceCommit: models.MergeCommit{
				CommitID: "cc759c13386c9a6e0aac1a70bb52a22c66961987",
				Url:      "https://totvstfs.visualstudio.com/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9/_apis/git/repositories/0f8a3f57-0bed-4451-86df-09d919e0d169/commits/cc759c13386c9a6e0aac1a70bb52a22c66961987",
			},
			LastMergeTargetCommit: models.MergeCommit{
				CommitID: "b893fe4dcf62be158a14f4b9eeae841bb8ad5321",
				Url:      "https://totvstfs.visualstudio.com/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9/_apis/git/repositories/0f8a3f57-0bed-4451-86df-09d919e0d169/commits/b893fe4dcf62be158a14f4b9eeae841bb8ad5321",
			},
			LastMergeCommit: models.MergeCommit{
				CommitID: "2fe080e092155e0c60d38460f8e4df9c962b889b",
				Author: models.UserAuthor{
					Name:  "Iassam da Silva de Souza",
					Email: "BRSSI0002@totvspartners.com.br",
					Date:  fakeTime,
				},
				Committer: models.UserAuthor{
					Name:  "Iassam da Silva de Souza",
					Email: "BRSSI0002@totvspartners.com.br",
					Date:  fakeTime,
				},
				Comment: "Merge pull request 181109 from feature/cetei-6703/shared-playbooks-and-activities-on-create into release/v1.95",
				Url:     "https://totvstfs.visualstudio.com/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9/_apis/git/repositories/0f8a3f57-0bed-4451-86df-09d919e0d169/commits/2fe080e092155e0c60d38460f8e4df9c962b889b",
			},
			Reviewers:          []models.Reviewers{},
			Url:                "https://totvstfs.visualstudio.com/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9/_apis/git/repositories/0f8a3f57-0bed-4451-86df-09d919e0d169/pullRequests/181109",
			SupportsIterations: true,
			ArtifactId:         "vstfs:///Git/PullRequestId/a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9%2f0f8a3f57-0bed-4451-86df-09d919e0d169%2f181109",
		},
		ResourceVersion: "1.0",
		CreatedDate:     fakeTime,
	}

	json_data, err := json.Marshal(azureContent)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/pull-request/created", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// assert.Equal(t, json_data, w.Body)
}
