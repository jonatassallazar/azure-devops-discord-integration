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

// TEST_SINKS selects which outbound sinks get a webhook URL for this test
// run, so the same route table can be exercised against Discord only,
// Google Chat only, or both:
//
//	go test ./cmd/server/...                     # both (default)
//	TEST_SINKS=discord go test ./cmd/server/...   # Discord only
//	TEST_SINKS=googlechat go test ./cmd/server/... # Google Chat only
//
// See the Makefile's test-discord / test-googlechat / test-e2e targets.
func stubHandler(w http.ResponseWriter, r *http.Request) {
	// The pipeline route resolves its triggering repository with a GET
	// against the Azure DevOps REST API; every other outbound call is a
	// sink POSTing a notification. Both just need a 2xx from the stub.
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(azuredevops.AzureRepository{
			ID:   "7a7c38b3-ac94-4e1e-b47e-d746c762079a",
			Name: "sample-service",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// writeTestEnv points every sink selected by TEST_SINKS, plus the Azure
// REST call, at the local stub server by writing a plain `.env` file (the
// always-attempted fallback per config.Config.LoadEnvironment) into this
// package's directory, so the suite runs fully offline with no manual env
// setup. It returns a restore func that puts back whatever `.env` (if any)
// was already there.
func writeTestEnv(stubURL string) (func(), error) {
	sinks := os.Getenv("TEST_SINKS")
	includeDiscord := sinks != "googlechat"
	includeGoogleChat := sinks != "discord"

	lines := []string{
		"APP_ENV=development",
		"GIN_MODE=test",
		"AZURE_ORGANIZATION=" + stubURL,
		"AZURE_PROJECT=sample-project",
		"AZURE_PAT_TOKEN=fake-token",
	}
	if includeDiscord {
		lines = append(lines,
			"DISCORD_PR_URL="+stubURL,
			"DISCORD_PIPELINE_URL="+stubURL,
			"DISCORD_RELEASE_URL="+stubURL,
		)
	}
	if includeGoogleChat {
		lines = append(lines,
			"GOOGLE_CHAT_PR_URL="+stubURL,
			"GOOGLE_CHAT_PIPELINE_URL="+stubURL,
			"GOOGLE_CHAT_RELEASE_URL="+stubURL,
		)
	}

	const envPath = ".env"

	original, err := os.ReadFile(envPath)
	hadOriginal := err == nil

	content := ""
	for _, line := range lines {
		content += line + "\n"
	}
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		return nil, fmt.Errorf("writing test .env: %w", err)
	}

	return func() {
		if hadOriginal {
			_ = os.WriteFile(envPath, original, 0o600)
		} else {
			_ = os.Remove(envPath)
		}
	}, nil
}

func TestMain(m *testing.M) {
	stub := httptest.NewServer(http.HandlerFunc(stubHandler))

	restore, err := writeTestEnv(stub.URL)
	if err != nil {
		fmt.Println(err)
		stub.Close()
		os.Exit(1)
	}

	code := m.Run()

	restore()
	stub.Close()

	os.Exit(code)
}

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

// Should answer the health probe without touching any sink
func TestHealthCheckRoute(t *testing.T) {
	r, _ := prepareRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, httpserver.RouteHealth, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

// Should send a notification to the configured sink(s) with created flag
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

// Should send a notification to the configured sink(s) with approved flag
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

// Should send a notification to the configured sink(s) with rejected flag
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

// Should NOT send a notification to any sink
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

// Should send a notification to the configured sink(s) with waiting flag
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

// Should send a notification to the configured sink(s) with complete flag
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

// Should send a notification to the configured sink(s) with conflict flag
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

// Should NOT send a notification to any sink
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

// Should send a notification to the configured sink(s) with Pipeline Details
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

// Should send a notification to the configured sink(s) with Pipeline Details
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

// Should send a notification to the configured sink(s) with Pipeline Details
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

// Should send a notification to the configured sink(s) with Pipeline Details
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

// Should send a notification to the configured sink(s) with Pipeline Details
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
