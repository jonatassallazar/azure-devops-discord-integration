package config

import (
	"errors"
	"log"
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
	c.loadDotEnv()

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

	return c.validate()
}

// loadDotEnv layers whatever dotenv files exist onto the process
// environment. Dotenv is a local-development convenience: in a container
// (Kubernetes ConfigMap/Secret, `docker run --env-file`, ...) the variables
// are already in the environment and no file exists, which is not an error.
// godotenv.Load is first-wins and never overrides an already-set variable,
// so real environment variables always take priority over file contents.
func (c *Config) loadDotEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	files := []string{".env." + env + ".local"}
	if env != "test" {
		files = append(files, ".env.local")
	}
	files = append(files, ".env."+env, ".env")

	var loaded []string
	for _, file := range files {
		if err := godotenv.Load(file); err == nil {
			loaded = append(loaded, file)
		}
	}

	if len(loaded) == 0 {
		log.Println("config: no env file found, reading configuration from environment variables only")
		return
	}

	log.Println("config: loaded env files", loaded)
}

// validate rejects a configuration the service cannot do anything useful
// with. Every webhook URL is individually optional — an unset one just
// disables that destination for that event category — but with none set at
// all the process would accept webhooks and deliver nothing, so fail fast
// with a message that names what is missing.
func (c *Config) validate() error {
	urls := []string{
		c.Discord.PRWebhookURL,
		c.Discord.PipelineWebhookURL,
		c.Discord.ReleaseWebhookURL,
		c.GoogleChat.PRWebhookURL,
		c.GoogleChat.PipelineWebhookURL,
		c.GoogleChat.ReleaseWebhookURL,
	}

	for _, url := range urls {
		if url != "" {
			return nil
		}
	}

	return errors.New("no webhook URL configured: set at least one of DISCORD_PR_URL, DISCORD_PIPELINE_URL, DISCORD_RELEASE_URL, GOOGLE_CHAT_PR_URL, GOOGLE_CHAT_PIPELINE_URL, GOOGLE_CHAT_RELEASE_URL")
}
