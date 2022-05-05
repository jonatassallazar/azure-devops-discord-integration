package config

import (
	"os"
)

type ConfigUrls struct {
	DiscordEnvPRUrl    string
	DiscordEnvBuildUrl string
}

func LoadEnvironment(C *ConfigUrls) {
	C.DiscordEnvPRUrl = os.Getenv("DISCORD_PR_URL")
	C.DiscordEnvBuildUrl = os.Getenv("DISCORD_BUILD_URL")

}
