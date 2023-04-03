package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type ConfigServer struct {
	AppEnv             string
	GinMode            string
	DiscordEnvPRUrl    string
	DiscordEnvBuildUrl string
}

func (c *ConfigServer) LoadEnvironment() error {
	err := c.loadDotEnv()
	if err != nil {
		return err
	}

	c.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	c.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	c.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	c.DiscordEnvBuildUrl = os.Getenv("DISCORD_BUILD_URL")

	return nil
}

func (c *ConfigServer) loadDotEnv() error {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	var err1 error
	var err2 error
	var err3 error
	var err4 error

	err1 = godotenv.Load(".env." + env + ".local")
	if env != "test" {
		err2 = godotenv.Load(".env.local")
	}
	err3 = godotenv.Load(".env." + env)
	err4 = godotenv.Load()

	if err1 != nil && err2 != nil && err3 != nil && err4 != nil {
		return errors.New("no env file was loaded")
	}

	return nil
}
