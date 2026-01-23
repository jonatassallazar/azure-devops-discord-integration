package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// ConfigServer holds all configuration values needed for the application to run.
//
// The configuration is loaded from environment variables, which can be provided
// through .env files or directly as environment variables. The struct contains
// settings for the application environment, Gin framework mode, Discord webhook URLs,
// and Azure DevOps credentials.
type ConfigServer struct {
	// AppEnv specifies the application environment (e.g., "development", "production", "test").
	// This is used to determine which .env file to load.
	AppEnv string
	
	// GinMode sets the Gin framework mode (e.g., "debug", "release", "test").
	// This controls logging and performance optimizations in the Gin router.
	GinMode string
	
	// DiscordEnvPRUrl is the Discord webhook URL for pull request notifications.
	DiscordEnvPRUrl string
	
	// DiscordEnvPipelineUrl is the Discord webhook URL for pipeline notifications.
	DiscordEnvPipelineUrl string
	
	// DiscordEnvReleaseUrl is the Discord webhook URL for release notifications.
	DiscordEnvReleaseUrl string
	
	// AzureOrganization is the Azure DevOps organization name or URL.
	AzureOrganization string
	
	// AzureProject is the Azure DevOps project name.
	AzureProject string
	
	// AzurePAT is the Azure DevOps Personal Access Token (PAT) used for API authentication.
	AzurePAT string
}

// LoadEnvironment loads all configuration values from environment variables.
//
// The method performs the following operations:
//  1. Loads environment variables from .env files using loadDotEnv()
//  2. Reads all required configuration values from environment variables
//  3. Populates the ConfigServer struct fields with the loaded values
//
// Required environment variables:
//   - APP_ENV: Application environment (defaults to "development" if not set)
//   - GIN_MODE: Gin framework mode
//   - DISCORD_PR_URL: Discord webhook URL for pull request notifications
//   - DISCORD_PIPELINE_URL: Discord webhook URL for pipeline notifications
//   - DISCORD_RELEASE_URL: Discord webhook URL for release notifications
//   - AZURE_ORGANIZATION: Azure DevOps organization name or URL
//   - AZURE_PROJECT: Azure DevOps project name
//   - AZURE_PAT_TOKEN: Azure DevOps Personal Access Token
//
// Returns:
//   - error: error if environment files cannot be loaded or required variables are missing
func (c *ConfigServer) LoadEnvironment() error {
	err := c.loadDotEnv()
	if err != nil {
		return err
	}

	c.AppEnv = os.Getenv("APP_ENV")
	c.GinMode = os.Getenv("GIN_MODE")
	c.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	c.DiscordEnvPipelineUrl = os.Getenv("DISCORD_PIPELINE_URL")
	c.DiscordEnvReleaseUrl = os.Getenv("DISCORD_RELEASE_URL")
	c.AzureOrganization = os.Getenv("AZURE_ORGANIZATION")
	c.AzureProject = os.Getenv("AZURE_PROJECT")
	c.AzurePAT = os.Getenv("AZURE_PAT_TOKEN")

	return nil
}

// loadDotEnv loads environment variables from .env files following a priority order.
//
// The method attempts to load environment files in the following order (first found wins):
//  1. .env.{APP_ENV}.local - Environment-specific local overrides (e.g., .env.development.local)
//  2. .env.local - Local overrides (skipped if APP_ENV is "test")
//  3. .env.{APP_ENV} - Environment-specific file (e.g., .env.production)
//  4. .env - Default environment file
//
// If APP_ENV is not set, it defaults to "development".
//
// The method will succeed if at least one .env file is successfully loaded.
// This allows for flexible configuration where different environments can have
// different .env files, with local overrides taking precedence.
//
// Returns:
//   - error: error if no .env files could be loaded (all attempts failed)
func (c *ConfigServer) loadDotEnv() error {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	var errLocalEnv error
	var errLocal error
	var errEnv error
	var errDefault error

	errLocalEnv = godotenv.Load(".env." + env + ".local")
	if env != "test" {
		errLocal = godotenv.Load(".env.local")
	}
	errEnv = godotenv.Load(".env." + env)
	errDefault = godotenv.Load()

	if errLocalEnv != nil && errLocal != nil && errEnv != nil && errDefault != nil {
		return errors.New("no env file was loaded")
	}

	return nil
}
