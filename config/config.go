package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type ConfigServer struct {
	AppEnv                string
	GinMode               string
	DiscordEnvPRUrl       string
	DiscordEnvPipelineUrl string
	DiscordEnvReleaseUrl  string
	AzureOrganization     string
	AzureProject          string
	AzurePAT              string
}

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
