package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigUrls struct {
	DiscordEnvPRUrl    string
	DiscordEnvBuildUrl string
}

func LoadEnvironment(C *ConfigUrls) error {
	err := loadDotEnv()
	if err != nil {
		return err
	}

	C.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	C.DiscordEnvBuildUrl = os.Getenv("DISCORD_BUILD_URL")

	return nil
}

func loadDotEnv() error {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

	return nil
}
