# Azure DevOps Discord Integration

Integration between Azure DevOps and Discord that allows notifying a Discord server about important events that occur on the Azure DevOps platform, such as Pull Requests, Pipelines, and Releases.

## 📋 Features

- **Pull Requests**: Notifications about creation, reviews, and status changes
- **Pipelines**: Notifications about pipeline execution status
- **Releases**: Notifications about release status

## 🚀 Prerequisites

- Go 1.17 or higher
- Docker and Docker Compose (optional, for container execution)
- Azure DevOps account with appropriate permissions
- Discord webhook configured

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
go run main.go
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

Create a `.env` file in the project root with the following variables:

```env
# Application environment (development, production, test)
APP_ENV=development

# Gin mode (debug, release, test)
GIN_MODE=debug

# Discord Webhook URLs
DISCORD_PR_URL=https://discord.com/api/webhooks/your-pr-webhook
DISCORD_PIPELINE_URL=https://discord.com/api/webhooks/your-pipeline-webhook
DISCORD_RELEASE_URL=https://discord.com/api/webhooks/your-release-webhook

# Azure DevOps Configuration
AZURE_ORGANIZATION=your-organization
AZURE_PROJECT=your-project
AZURE_PAT_TOKEN=your-personal-access-token
```

### Environment File Loading Order

The system loads environment files in the following order (the last one has priority):

1. `.env`
2. `.env.{APP_ENV}`
3. `.env.local` (ignored if `APP_ENV=test`)
4. `.env.{APP_ENV}.local`

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
├── config/           # Configuration and environment variable loading
├── controllers/      # Controllers for each event type
├── models/           # Data models (Azure DevOps and Discord)
├── server/           # Server configuration and routes
├── main.go          # Application entry point
├── Dockerfile       # Docker build configuration
└── docker-compose.yml # Docker Compose configuration
```

## 🧪 Testing

Run tests with:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test -cover ./...
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

- Go 1.17+
- Testing tools (already included in dependencies)

### Run in Development Mode

```bash
APP_ENV=development GIN_MODE=debug go run main.go
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
