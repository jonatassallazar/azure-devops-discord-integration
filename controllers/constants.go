package controllers

// HEADER_APP_JSON is the Content-Type header value for JSON content.
// It is used when making HTTP requests that send JSON payloads.
const HEADER_APP_JSON = "application/json"

// HTTP route constants for the Azure DevOps webhook endpoints.
//
// These routes define the paths where the server listens for webhook events
// from Azure DevOps for different types of events (pull requests, pipelines, releases).

// Pull Request Routes

// CREATED_ROUTE is the HTTP route endpoint for pull request created events.
// Azure DevOps sends webhook notifications to this endpoint when a new pull request is created.
const CREATED_ROUTE = "/pull-request/created"

// REVIEW_ROUTE is the HTTP route endpoint for pull request review events.
// Azure DevOps sends webhook notifications to this endpoint when a pull request review is submitted.
const REVIEW_ROUTE = "/pull-request/review"

// STATUS_ROUTE is the HTTP route endpoint for pull request status change events.
// Azure DevOps sends webhook notifications to this endpoint when a pull request status changes
// (e.g., completed, abandoned, etc.).
const STATUS_ROUTE = "/pull-request/status"

// Pipeline Routes

// PIPELINE_ROUTE is the HTTP route endpoint for pipeline events.
// Azure DevOps sends webhook notifications to this endpoint when pipeline status changes occur
// (e.g., succeeded, failed, stopped).
const PIPELINE_ROUTE = "/pipeline/"

// Release Routes

// RELEASE_ROUTE is the HTTP route endpoint for release events.
// Azure DevOps sends webhook notifications to this endpoint when release status changes occur.
const RELEASE_ROUTE = "/release/"
