package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigUrls struct {
	DiscordEnvPRUrl    string
	DiscordEnvBuildUrl string
}

func LoadEnvironment(C *ConfigUrls) {
	loadDotEnv()

	C.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	C.DiscordEnvBuildUrl = os.Getenv("DISCORD_BUILD_URL")
}

func loadDotEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
