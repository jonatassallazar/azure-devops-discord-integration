package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// envKeys is every variable LoadEnvironment reads, cleared before each case
// so a value inherited from the developer's shell can't decide the result.
var envKeys = []string{
	"APP_ENV",
	"GIN_MODE",
	"PUBLIC_BASE_URL",
	"AZURE_ORGANIZATION",
	"AZURE_PROJECT",
	"AZURE_PAT_TOKEN",
	"DISCORD_PR_URL",
	"DISCORD_PIPELINE_URL",
	"DISCORD_RELEASE_URL",
	"GOOGLE_CHAT_PR_URL",
	"GOOGLE_CHAT_PIPELINE_URL",
	"GOOGLE_CHAT_RELEASE_URL",
}

// isolate runs the test in an empty temp directory (so no dotenv file is
// found) with every configuration variable unset, restoring both afterwards.
func isolate(t *testing.T) {
	t.Helper()

	for _, key := range envKeys {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}

	wd, err := os.Getwd()
	assert.NoError(t, err)

	assert.NoError(t, os.Chdir(t.TempDir()))
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

// A container deployment (Kubernetes ConfigMap/Secret, docker --env-file)
// has no dotenv file on disk: environment variables alone must be enough.
func TestLoadEnvironmentWithoutEnvFile(t *testing.T) {
	isolate(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("DISCORD_PR_URL", "https://discord.example/pr")
	t.Setenv("AZURE_PROJECT", "sample-project")
	t.Setenv("PUBLIC_BASE_URL", "https://notify.example.com")

	var cfg Config
	assert.NoError(t, cfg.LoadEnvironment())

	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "https://discord.example/pr", cfg.Discord.PRWebhookURL)
	assert.Equal(t, "sample-project", cfg.Azure.Project)
	assert.Equal(t, "https://notify.example.com", cfg.PublicBaseURL)
}

// With no dotenv file and no webhook URL there is nothing to deliver to, so
// startup should fail with a message naming the variables to set.
func TestLoadEnvironmentWithoutAnyWebhookURL(t *testing.T) {
	isolate(t)

	var cfg Config
	err := cfg.LoadEnvironment()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DISCORD_PR_URL")
}

func TestLoadEnvironmentFromDotEnvFile(t *testing.T) {
	isolate(t)

	dir, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(filepath.Join(dir, ".env"),
		[]byte("GOOGLE_CHAT_RELEASE_URL=https://chat.example/release\nAZURE_PAT_TOKEN=fake-token\n"), 0o600))

	var cfg Config
	assert.NoError(t, cfg.LoadEnvironment())

	assert.Equal(t, "https://chat.example/release", cfg.GoogleChat.ReleaseWebhookURL)
	assert.Equal(t, "fake-token", cfg.Azure.PATToken)
}

// godotenv never overrides an already-set variable, so a real environment
// variable wins over the dotenv file — the behaviour a container relies on.
func TestEnvironmentVariableWinsOverDotEnvFile(t *testing.T) {
	isolate(t)
	t.Setenv("DISCORD_PR_URL", "https://discord.example/from-environment")

	dir, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(filepath.Join(dir, ".env"),
		[]byte("DISCORD_PR_URL=https://discord.example/from-file\n"), 0o600))

	var cfg Config
	assert.NoError(t, cfg.LoadEnvironment())

	assert.Equal(t, "https://discord.example/from-environment", cfg.Discord.PRWebhookURL)
}
