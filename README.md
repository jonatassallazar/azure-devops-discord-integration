# Azure DevOps Discord Integration

Integration between Azure DevOps and chat apps that notifies about important events on the
Azure DevOps platform — Pull Requests, Pipelines, and Releases — via Discord and/or Google
Chat incoming webhooks.

## 📋 Features

- **Pull Requests**: Notifications about creation, reviews, and status changes
- **Pipelines**: Notifications about pipeline execution status
- **Releases**: Notifications about release status
- **Multiple outbound destinations**: Discord and Google Chat can both be configured at once;
  each event category (PR / pipeline / release) is delivered to whichever are set

## 🚀 Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose (optional, for container execution)
- Azure DevOps account with appropriate permissions
- A Discord webhook and/or a Google Chat webhook configured

## 📦 Installation

### Local Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/azure-devops-discord-integration.git
cd azure-devops-discord-integration
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables (see [Configuration](#-configuration) section)

4. Run the server:
```bash
go run ./cmd/server
```

### Docker Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/azure-devops-discord-integration.git
cd azure-devops-discord-integration
```

2. Configure environment variables in `.env` or `.env.local` file

3. Run with Docker Compose:
```bash
docker-compose up -d
```

## ⚙️ Configuration

Copy [`.env.example`](.env.example) to `.env` in the project root and fill in the values you
need — it documents every supported variable with inline comments, so it's the fastest way to
get onboarded:

```bash
cp .env.example .env
```

At least one webhook URL (Discord or Google Chat) should be set per event category you want
notifications for; both may be set to deliver to both destinations at once.

### Environment File Loading Order

The system loads environment files in the following order. Note that the **first file found
wins** — dotenv loading does not override a variable that's already set, so if a variable
appears in more than one file, the earliest file in this list takes priority:

1. `.env.{APP_ENV}.local`
2. `.env.local` (skipped entirely when `APP_ENV=test`)
3. `.env.{APP_ENV}`
4. `.env`

### How to Get Azure DevOps Personal Access Token (PAT)

1. Access Azure DevOps
2. Go to **User Settings** > **Personal Access Tokens**
3. Click on **New Token**
4. Configure the necessary permissions (Code: Read, Build: Read, Release: Read)
5. Copy the generated token

### How to Create Discord Webhooks

1. Access your Discord server settings
2. Go to **Integrations** > **Webhooks**
3. Click on **New Webhook**
4. Configure the name and channel
5. Copy the webhook URL

### How to Create Google Chat Webhooks

1. Open the Google Chat space you want to notify
2. Go to the space name menu > **Apps & integrations** > **Manage webhooks**
3. Click **Add a webhook**, give it a name, and save
4. Copy the generated webhook URL

## 🔌 Endpoints

The server exposes the following endpoints that should be configured as Service Hooks in Azure DevOps:

### Pull Requests

- `POST /pull-request/created` - Notifies when a Pull Request is created
- `POST /pull-request/review` - Notifies when a Pull Request is reviewed
- `POST /pull-request/status` - Notifies when a Pull Request status changes

### Pipelines

- `POST /pipeline/` - Notifies about pipeline status changes

### Releases

- `POST /release/` - Notifies about release status changes

## 🔧 Azure DevOps Configuration

To configure Service Hooks in Azure DevOps:

1. Access your project in Azure DevOps
2. Go to **Project Settings** > **Service hooks**
3. Click on **Create subscription**
4. Select the service (Pull Request, Build, Release, etc.)
5. Configure the desired events
6. In the action, select **Web Hooks**
7. Enter your server URL: `http://your-server:8080/corresponding-endpoint`
8. Test the connection and save

## 🏗️ Project Structure

```
azure-devops-discord-integration/
├── cmd/server/                   # Application entry point (main.go)
├── internal/
│   ├── config/                   # Environment variable loading
│   ├── azuredevops/              # Azure DevOps webhook handlers, models, event mapping
│   ├── notify/                   # Vendor-neutral notification model + Sink interface + fan-out dispatcher
│   │   ├── discord/              # Discord Sink implementation
│   │   └── googlechat/           # Google Chat Sink implementation
│   └── httpserver/               # gin server/router wiring
├── Dockerfile                    # Docker build configuration
└── docker-compose.yml            # Docker Compose configuration
```

Inbound events come from Azure DevOps today (`internal/azuredevops`); outbound delivery goes
through the `notify.Sink` interface, currently implemented by Discord and Google Chat
(`internal/notify/discord`, `internal/notify/googlechat`). Adding another outbound
destination means implementing `notify.Sink` in a new package and wiring it into
`internal/httpserver/router.go` — no changes needed elsewhere.

## 🧪 Testing

Run the full suite with:

```bash
go test ./...
# or: make test
```

To run tests with coverage:

```bash
go test -cover ./...
# or: make test-cover
```

`cmd/server`'s tests are end-to-end: they spin up a local stub server and, by default, point
both the Discord and Google Chat sinks (plus the Azure REST call) at it, so the suite runs
fully offline with no manual env setup. Use `TEST_SINKS` to exercise just one destination:

```bash
TEST_SINKS=discord go test ./cmd/server/...    # or: make test-discord
TEST_SINKS=googlechat go test ./cmd/server/... # or: make test-googlechat
TEST_SINKS=both go test ./cmd/server/...       # or: make test-e2e (also the default)
```

`internal/azuredevops`'s tests are pure unit tests with no env or network dependency:

```bash
go test ./internal/azuredevops/...             # or: make test-unit
```

## 🐳 Docker

### Build Image

```bash
docker build -t azure-devops-discord-integration .
```

### Run Container

```bash
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  --name azure-devops-discord \
  azure-devops-discord-integration
```

## 📝 Development

### Development Requirements

- Go 1.22+
- Testing tools (already included in dependencies)

### Run in Development Mode

```bash
APP_ENV=development GIN_MODE=debug go run ./cmd/server
```

## 🔒 Security

- **Never** commit `.env` or `.env.*` files to the repository
- Keep your tokens and webhooks secure
- Use environment variables or secrets in production environments
- Consider using HTTPS in production

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or pull requests.

## 📧 Support

For questions and support, open an issue in the repository.
