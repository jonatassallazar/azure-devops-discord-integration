package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AzureConfig struct {
	Organization string
	Project      string
	PATToken     string
}

type DiscordConfig struct {
	PRWebhookURL       string
	PipelineWebhookURL string
	ReleaseWebhookURL  string
}

type GoogleChatConfig struct {
	PRWebhookURL       string
	PipelineWebhookURL string
	ReleaseWebhookURL  string
}

type Config struct {
	AppEnv     string
	GinMode    string
	Azure      AzureConfig
	Discord    DiscordConfig
	GoogleChat GoogleChatConfig
}

func (c *Config) LoadEnvironment() error {
	if err := c.loadDotEnv(); err != nil {
		return err
	}

	c.AppEnv = os.Getenv("APP_ENV")
	c.GinMode = os.Getenv("GIN_MODE")

	c.Azure = AzureConfig{
		Organization: os.Getenv("AZURE_ORGANIZATION"),
		Project:      os.Getenv("AZURE_PROJECT"),
		PATToken:     os.Getenv("AZURE_PAT_TOKEN"),
	}

	c.Discord = DiscordConfig{
		PRWebhookURL:       os.Getenv("DISCORD_PR_URL"),
		PipelineWebhookURL: os.Getenv("DISCORD_PIPELINE_URL"),
		ReleaseWebhookURL:  os.Getenv("DISCORD_RELEASE_URL"),
	}

	c.GoogleChat = GoogleChatConfig{
		PRWebhookURL:       os.Getenv("GOOGLE_CHAT_PR_URL"),
		PipelineWebhookURL: os.Getenv("GOOGLE_CHAT_PIPELINE_URL"),
		ReleaseWebhookURL:  os.Getenv("GOOGLE_CHAT_RELEASE_URL"),
	}

	return nil
}

func (c *Config) loadDotEnv() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	errLocalEnv := godotenv.Load(".env." + env + ".local")

	var errLocal error
	if env != "test" {
		errLocal = godotenv.Load(".env.local")
	}

	errEnv := godotenv.Load(".env." + env)
	errDefault := godotenv.Load()

	if errLocalEnv != nil && errLocal != nil && errEnv != nil && errDefault != nil {
		return errors.New("no env file was loaded")
	}

	return nil
}
