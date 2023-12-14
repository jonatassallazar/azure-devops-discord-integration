package main

import (
	"bytes"
	config "discord-azure-integration/config"
	"discord-azure-integration/controllers"
	"discord-azure-integration/server"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareRouter() (*gin.Engine, error) {
	var c config.ConfigServer

	err := c.LoadEnvironment()
	if err != nil {
		return nil, err
	}

	var s = server.Server{
		ConfigServer: &c,
	}

	r := s.SetupRouter()

	return r, nil
}

func TestConfigE2ETesting(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	_, err := prepareRouter()
	assert.Nil(t, err)
}

// Should send message to Discord Webhook with created flag
func TestCreateRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadCreatePR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.CREATED_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with approved flag
func TestReviewApprovedRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadApprovedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with rejected flag
func TestReviewRejectedRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadRejectedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should NOT send message to Discord Webhook
func TestReviewNeutralRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadNeutralPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Should send message to Discord Webhook with waiting flag
func TestReviewWaitingForAuthorRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadWaitingPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with complete flag
func TestCompletedPullRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadCompletedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.STATUS_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with conflict flag
func TestConflictPullRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadConflictPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.STATUS_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should NOT send message to Discord Webhook
func TestOrdinaryUpdatePullRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadOrdinaryUpdatePR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.STATUS_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Should send message to Discord Webhook with Pipeline Details
func TestPipelineRouteWithSuccessResultSucceeded(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadPipelineUpdateSucceeded)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.PIPELINE_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with Pipeline Details
func TestPipelineRouteWithSuccessResultFailed(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadPipelineUpdateFailed)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.PIPELINE_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with Pipeline Details
func TestPipelineRouteWithSuccessResultStopped(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadPipelineUpdateStopped)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.PIPELINE_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with Pipeline Details
func TestPipelineRouteWithSuccessResultDefault(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadPipelineUpdateDefault)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.PIPELINE_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Should send message to Discord Webhook with Pipeline Details
func TestReleaseRouteWithSuccess(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadReleaseSuccess)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, controllers.RELEASE_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
