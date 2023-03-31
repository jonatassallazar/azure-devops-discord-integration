package main

import (
	"bytes"
	models "discord-azure-integration/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigE2ETesting(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	_, err := setupRouter()

	assert.Nil(t, err)
}

func TestPingRoute(t *testing.T) {
	router, _ := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestCreateRequestRoute(t *testing.T) {
	router, _ := setupRouter()

	fakeTime := time.Now()

	azureContent := models.AzureRequest{
		SubscriptionId: "2f5389a6-c75a-4048-8eba-33c0767d9ad6",
		NotificationId: 1094,
		ID:             "1c75aaf4-a8de-4748-835b-8bf7964ac0fa",
		EventType:      "git.pullrequest.created",
		PublisherId:    "doom",
		Message: models.Text{
			Text: "John Don created pull request 18111209 (branch > release) in repository",
		},
		DetailedMessage: models.Text{
			Text: "John Don created pull request 18111209 (branch > release) in repository",
		},
		Resource: models.Resource{
			Repository: models.Repository{
				ID:   "0f8a3f57-0bed-4451-86df-09d919e0d169",
				Name: "repository_name",
				Url:  "https://azure.com/",
				Project: models.Project{
					ID:             "a133d6b5-dfbe-49b7-9a8f-d8cb27e3adc9",
					Name:           "Project Doom",
					Description:    "Team Doom repository.",
					Url:            "https://doom.azure.com/",
					State:          "wellFormed",
					Revision:       123,
					Visibility:     "private",
					LastUpdateTime: fakeTime,
				},
				Size:            132,
				RemoteUrl:       "https://doom.azure.com/",
				SshUrl:          "doom@vs-ssh.visualstudio.com:v3/project/repo",
				WebUrl:          "https://doom.azure.com/",
				IsDisabled:      false,
				IsInMaintenance: false,
			},
			PullRequestId: 12333,
			CodeReviewId:  4123123,
			Status:        "active",
			CreatedBy: models.CreatedBy{
				DisplayName: "Jonh Don",
				Url:         "https://johndon.com",
				ID:          "a4dd11b6-b7bc-64e7-8a8e-1617179bef68",
				UniqueName:  "jonh_don@doom.com.br",
				ImageUrl:    "http://johndon.com",
			},
			CreationDate:  fakeTime,
			Title:         "feature > release",
			Description:   "PR Description",
			SourceRefName: "refs/heads/feature",
			TargetRefName: "refs/heads/release",
			MergeStatus:   "succeeded",
			IsDraft:       false,
			MergeID:       "9634d63f-a0ab-49de-8767-52e80ac18240",
			LastMergeSourceCommit: models.MergeCommit{
				CommitID: "dfasfasdfb234khb34kjb14",
				Url:      "https://",
			},
			LastMergeTargetCommit: models.MergeCommit{
				CommitID: "1h3hl41bl4hjb1l2j4hbl12j4hb123",
				Url:      "https://",
			},
			LastMergeCommit: models.MergeCommit{
				CommitID: "fdsufhsdf76g9sdfynnh7bdc7208d3",
				Author: models.UserAuthor{
					Name:  "John Don",
					Email: "john_don@doom.com",
					Date:  fakeTime,
				},
				Committer: models.UserAuthor{
					Name:  "John Don",
					Email: "john_don@doom.com",
					Date:  fakeTime,
				},
				Comment: "Merge pull request 1312312 from feature into release",
				Url:     "https://",
			},
			Reviewers:          []models.Reviewers{},
			Url:                "https://",
			SupportsIterations: true,
			ArtifactId:         "",
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
