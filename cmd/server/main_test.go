package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"azuredevops-notify/internal/azuredevops"
	"azuredevops-notify/internal/config"
	"azuredevops-notify/internal/httpserver"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func prepareRouter() (*gin.Engine, error) {
	var cfg config.Config

	if err := cfg.LoadEnvironment(); err != nil {
		return nil, err
	}

	s := httpserver.Server{Config: &cfg}

	return s.SetupRouter(), nil
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

	jsonData, err := json.Marshal(fakePayloadCreatePR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteCreatedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadApprovedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteReviewedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadRejectedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteReviewedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadNeutralPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteReviewedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadWaitingPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteReviewedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadCompletedPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteStatusUpdatedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadConflictPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteStatusUpdatedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadOrdinaryUpdatePR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteStatusUpdatedPR, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadPipelineUpdateSucceeded)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RoutePipeline, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadPipelineUpdateFailed)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RoutePipeline, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadPipelineUpdateStopped)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RoutePipeline, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadPipelineUpdateDefault)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RoutePipeline, bytes.NewBuffer(jsonData))
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

	jsonData, err := json.Marshal(fakePayloadReleaseSuccess)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, azuredevops.RouteRelease, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
