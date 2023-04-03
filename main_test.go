package main

import (
	"bytes"
	config "discord-azure-integration/Config"
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

const REVIEW_ROUTE = "/pull-request/review"

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
	req, err := http.NewRequest(http.MethodPost, "/pull-request/created", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
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
	req, err := http.NewRequest(http.MethodPost, REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
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
	req, err := http.NewRequest(http.MethodPost, REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// Should not send message to Discord Webhook
func TestReviewNeutralRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadNeutralPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// Should not send message to Discord Webhook
func TestReviewWaitingForAuthorRequestRoute(t *testing.T) {
	r, _ := prepareRouter()

	json_data, err := json.Marshal(fakePayloadWaitingPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, REVIEW_ROUTE, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
		return
	}
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
