# Rule MCP Server

A MCP (Model Context Protocol) server that allows AI agents (Cursor, Cline) to retrieve and apply common coding rules.

## Features

- Project-specific rule management
- Language-specific global rule management
- Project-specific global rule application settings
- Code rule violation validation
- Rule distribution via MCP
- RESTful API endpoints
- **React-based Web UI** for rule management
- **Multi-language support (i18n)**: English, Japanese, Chinese, Hindi, Spanish, Arabic
- **Dark theme support**: Light/Dark mode toggle
- **Clean Architecture** for high maintainability

## Tech Stack

### Backend

- **Go 1.21+** + **Gin Web Framework**
- **Clean Architecture** (Domain, Usecase, Interface, Infrastructure)
- **PostgreSQL** database
- **MCP (Model Context Protocol)** support

### Frontend

- **React 18** + **TypeScript**
- **Material-UI (MUI)** component library
- **React Router** for routing
- **i18next** for internationalization
- **Axios** for API communication

## Setup

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- Docker and Docker Compose
- PostgreSQL 15 or higher (automatically installed with Docker)

### Initial Admin Account

The following initial admin account is automatically created on first system startup:

- **Username**: `admin`
- **Password**: `admin123`
- **Email**: `admin@rulemcp.com`
- **Role**: Administrator (admin)

**Important**: Please change the password after first login.

### Installation

```bash
git clone <repository-url>
cd RuleMCPServer

# Backend dependencies
go mod tidy

# Frontend dependencies
cd frontend
npm install
cd ..
```

### Running

#### Development Environment

```bash
# Backend (safe port 18081)
PORT=18081 go run ./cmd/server

# Frontend
cd frontend
npm start

# Using Makefile
make run        # Development environment (port 18081)
make run-frontend
```

#### Production Environment (Docker)

```bash
# Create production environment variables file
cp env.prod.example .env.prod
# Edit .env.prod file values for production environment

# Deploy production environment
make -f Makefile.prod deploy

# Check production environment status
make -f Makefile.prod status

# Check production environment logs
make -f Makefile.prod logs

# Stop production environment
make -f Makefile.prod down

# Clean up production environment
make -f Makefile.prod clean
```

#### Production Environment Features

- **Enhanced Security**: Non-root user execution, environment variable configuration
- **Performance Optimization**: Multi-stage build, lightweight Alpine Linux
- **Health Checks**: Service health monitoring
- **Log Management**: Structured log output and rotation
- **Backup**: Automatic database backup functionality
- **Scalability**: Docker Swarm and Kubernetes ready

## Architecture

### Clean Architecture

```
cmd/server/          # Entry point
├── internal/
│   ├── domain/      # Entities, repository interfaces
│   ├── usecase/     # Business logic
│   ├── interface/   # HTTP handlers, MCP handlers
│   └── infrastructure/ # Database implementation
└── frontend/        # React frontend
```

### Layer Structure

- **Domain Layer**: Business entities and rules
- **Usecase Layer**: Application business logic
- **Interface Layer**: HTTP API, MCP protocol
- **Infrastructure Layer**: PostgreSQL, file system

## API Endpoints

### Health Check

```bash
GET /api/v1/health
```

### Rule Retrieval

```bash
GET /api/v1/rules?project={project_id}
```

### Code Validation

```bash
POST /api/v1/rules/validate
Content-Type: application/json

{
  "project_id": "web-app",
  "code": "console.log('test')"
}
```

### Project Management

```bash
# Get project list
GET /api/v1/projects

# Create project
POST /api/v1/projects
Content-Type: application/json

{
  "project_id": "new-project",
  "name": "New Project",
  "description": "A new project description",
  "language": "javascript",
  "apply_global_rules": true
}
```

### Rule Management

```bash
# Create rule
POST /api/v1/rules
Content-Type: application/json

{
  "project_id": "new-project",
  "rule_id": "no-debug-code",
  "name": "No Debug Code",
  "description": "Debug code should not be in production",
  "type": "style",
  "severity": "warning",
  "pattern": "debugger",
  "message": "Debug code detected. Remove before production."
}

# Delete rule
DELETE /api/v1/rules/{project_id}/{rule_id}
```

### Global Rule Management

```bash
# Get language-specific global rules
GET /api/v1/global-rules/{language}

# Create global rule
POST /api/v1/global-rules
Content-Type: application/json

{
  "language": "javascript",
  "rule_id": "no-console-log",
  "name": "No Console Log",
  "description": "Console.log statements should not be in production code",
  "type": "style",
  "severity": "warning",
  "pattern": "console.log",
  "message": "Console.log detected. Use proper logging framework in production."
}

# Delete global rule
DELETE /api/v1/global-rules/{language}/{rule_id}

# Get available languages
GET /api/v1/languages
```

## MCP (Model Context Protocol) Server

This server operates as an MCP server, allowing Cursor and Cline to directly retrieve rules.

### MCP Endpoints

```
POST /mcp/request    # HTTP MCP request
GET  /mcp/ws         # WebSocket MCP connection
```

### MCP Methods

- **`getRules`**: Get project rules
- **`validateCode`**: Validate code rule violations
- **`getProjectInfo`**: Get project information

### Standard MCP Server Configuration (Recommended)

This project provides a **standard MCP (Model Context Protocol) server**.

#### **Features**
- ✅ **Standard MCP SDK**: Full compliance using `@modelcontextprotocol/sdk`
- ✅ **StdioServerTransport**: Compatible with standard MCP clients
- ✅ **Docker Support**: Stable operation in production
- ✅ **Template Configuration**: Easy setup

#### **1. MCP Server Installation**

##### **pnpm (Recommended)**
```bash
# Global installation
pnpm add -g rule-mcp-server

# Or without installation
pnpm dlx rule-mcp-server
```

##### **Development Build**
```bash
# Install dependencies
make install-mcp

# Build MCP server
make build-mcp
```

#### **2. Environment-specific Configuration**

##### **pnpm Package Usage (Recommended)**
```bash
# Use pnpm package configuration template
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json
```

Configuration example (pnpm package):
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": ["dlx", "rule-mcp-server"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

##### **Docker Environment**
```bash
# Start Docker environment
make docker-up

# Use Docker configuration template
cp config/docker-mcp-config.template.json ~/.cursor/mcp_settings.json
```

Configuration example (Docker environment):
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/your/RuleMCPServer/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management (Docker)",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

##### **Development Environment**
```bash
# Start development server
make run

# Use standard configuration template
cp config/standard-mcp-config.template.json ~/.cursor/mcp_settings.json
```

Configuration example (Development environment):
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/your/RuleMCPServer/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18081",
        "MCP_API_KEY": ""
      },
      "description": "Standard MCP Server for Rule Management (Development)",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

#### **3. Configuration File Placement**
- **Cursor**: `~/.cursor/mcp_settings.json`
- **Cline**: `~/.cline/mcp_settings.json`
- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json`

#### **4. Available Tools**
The standard MCP server provides the following tools:

| Tool Name            | Description                         | Required Parameters         |
|---------------------|-------------------------------------|----------------------------|
| `getRules`          | Get project rules                   | `project_id`               |
| `validateCode`      | Code validation                     | `project_id`, `code`       |
| `getProjectInfo`    | Get project information             | `project_id`               |
| `autoDetectProject` | Auto-detect project                 | `path`                     |
| `scanLocalProjects` | Scan local projects                 | `base_path` (optional)     |
| `getGlobalRules`    | Get global rules                    | `language`                 |

#### **5. Available Resources**
The standard MCP server provides the following resources:

| Resource URI                      | Description                   |
|-----------------------------------|-------------------------------|
| `rule://projects/list`           | Project list                  |
| `rule://{project_id}/rules`      | Project-specific rules        |
| `rule://{project_id}/info`       | Project information           |
| `rule://global-rules/{language}` | Language-specific global rules|

## Frontend Features

### Multi-language Support (i18n)

Supports the following languages:

- **English (en)**: Default language
- **Japanese (ja)**: Full support
- **Chinese (zh-CN)**: Full support
- **Hindi (hi)**: Full support
- **Spanish (es)**: Full support
- **Arabic (ar)**: Full support (RTL support)

### Theme Switching

- **Light Theme**: Bright and readable
- **Dark Theme**: Eye-friendly night mode
- **Automatic Settings Save**: Saved to browser localStorage

### Web UI Features

- **Project Management**: Create, edit, delete
- **Rule Management**: Project-specific rules, global rules
- **Code Validation**: Real-time rule violation checking
- **Responsive Design**: Mobile and tablet support

## Rule Definition

### Rule Creation Items

1. **Rule ID**: Unique identifier for the rule (e.g., `no-console-log`)
2. **Name**: Display name for the rule (e.g., `No Console Log`)
3. **Description**: Rule content, purpose, and reason (e.g., `Console.log statements should not be in production code`)
4. **Type**: Rule category (naming, formatting, security, performance, etc.)
5. **Severity**: Importance level (error, warning, info)
6. **Pattern**: Regular expression pattern to detect (e.g., `console\.log`)
7. **Message**: Correction instructions for violations (e.g., `Console.log detected. Use proper logging framework in production.`)
8. **Active**: Rule enabled/disabled

### Sample Rules

```json
{
  "web-app": {
    "project_id": "web-app",
    "rules": [
      {
        "id": "no-console-log",
        "name": "No Console Log",
        "description": "Console.log statements should not be in production code",
        "type": "style",
        "severity": "warning",
        "pattern": "console.log",
        "message": "Console.log detected. Use proper logging framework in production."
      }
    ]
  }
}
```

## Configuration

### Environment Variables

- `PORT`: Server port number (default: 8080)
- `HOST`: Server host address (default: 0.0.0.0)
- `ENVIRONMENT`: Execution environment (development/production, default: development)
- `LOG_LEVEL`: Log level (default: info)

### Database Configuration

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_NAME`: Database name (default: rule_mcp_db)
- `DB_USER`: Database user (default: rule_mcp_user)
- `DB_PASSWORD`: Database password

### Port Configuration

To avoid port conflicts for developers, the following ports are used:

- **Rule MCP Server**: 18081 (development) / 18080 (production)
- **PostgreSQL**: 15432 (host) → 5432 (container)
- **Web UI**: 13000 (host) → 80 (container)

This avoids conflicts with common ports (8080, 5432, 3000).

## 🚀 Quick Start

### pnpm Package Usage (Recommended)

```bash
# 1. Install MCP server
pnpm add -g rule-mcp-server

# 2. Create configuration file
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json

# 3. Start using with AI agents (Cursor/Cline)!
```

### Docker Environment

```bash
# 1. Clone repository
git clone https://github.com/AkitoSakurabaCreator/Rule-MCP-Server.git
cd Rule-MCP-Server

# 2. Start Docker environment
make docker-up

# 3. Build MCP server
make install-mcp && make build-mcp

# 4. Create configuration file
cp config/docker-mcp-config.template.json ~/.cursor/mcp_settings.json
# Edit path: ${PROJECT_PATH} → /path/to/your/Rule-MCP-Server

# 5. Start using with AI agents (Cursor/Cline)!
```

### Development Environment

```bash
# 1. Install dependencies
go mod tidy
cd frontend && npm install && cd ..

# 2. Start servers
make run        # Backend (port 18081)
make run-frontend  # Frontend (port 3000)

# 3. Build MCP server
make install-mcp && make build-mcp

# 4. Create configuration file
cp config/standard-mcp-config.template.json ~/.cursor/mcp_settings.json
```

## 🌟 Key Features

### ✅ Standard MCP Compliance
- **Full MCP Compatibility**: Using `@modelcontextprotocol/sdk`
- **Standard Tools & Resources**: Perfect integration with AI agents
- **StdioServerTransport**: Standard communication protocol

### ✅ Rich Features
- **6 MCP Tools**: Rule retrieval, code validation, project auto-detection, etc.
- **5 MCP Resources**: Project list, rule information, global rules
- **Project Auto-detection**: AI automatically applies appropriate rules

### ✅ Production Ready
- **Docker Environment**: Stable production operation
- **PostgreSQL**: High-performance database
- **Authentication & Authorization**: API key, session authentication
- **Multi-language Support**: 6 language support (Japanese, English, Chinese, etc.)

### ✅ Developer Friendly
- **Clean Architecture**: Highly maintainable design
- **Web UI**: React-based management interface
- **Template Configuration**: Easy setup
- **Comprehensive Documentation**: Detailed usage instructions

## 🎯 Use Cases

### Individual Developers
```bash
# Simple usage without authentication
curl http://localhost:18080/api/v1/rules?project_id=my-project
```

### Team Development
```bash
# Team management with API key
curl -H "X-API-Key: team_key" http://localhost:18080/api/v1/projects
```

### AI Agent Integration
```json
// Automatic rule application in Cursor/Cline
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "node",
      "args": ["/path/to/Rule-MCP-Server/cmd/mcp-server/build/index.js"],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080"
      }
    }
  }
}
```

## 📊 System Requirements

### Minimum Requirements
- **OS**: Linux, macOS, Windows
- **Go**: 1.21+
- **Node.js**: 18+
- **Docker**: 20.10+ (recommended)

### Recommended Requirements
- **Memory**: 2GB or more
- **Storage**: 1GB or more
- **CPU**: 2 cores or more

## 🔧 Development & Contribution

### Development Environment Setup

```bash
# 1. Fork & Clone
git clone https://github.com/your-username/Rule-MCP-Server.git
cd Rule-MCP-Server

# 2. Create development branch
git checkout -b feature/your-feature

# 3. Install dependencies
make deps
make install-mcp

# 4. Start development servers
make run
make run-frontend
```

### Running Tests

```bash
# Backend tests
make test

# Frontend tests
cd frontend && npm test

# MCP server tests
cd cmd/mcp-server && npm test
```

### Code Quality Check

```bash
# Go language
make fmt
make lint

# TypeScript
cd cmd/mcp-server && npm run lint
```

## 🤝 Contribution Guidelines

### Contribution Flow

1. **Create Issue**: Bug reports or feature proposals
2. **Fork**: Fork the repository
3. **Create Branch**: `feature/your-feature` or `fix/your-fix`
4. **Development**: Code changes and test additions
5. **Pull Request**: Submit with detailed description

### Commit Message Convention

```bash
# Feature addition
feat: add project auto-detection feature

# Bug fix
fix: resolve MCP server connection issue

# Documentation update
docs: update README with Docker setup

# Refactoring
refactor: improve error handling in MCP handlers
```

### Development Best Practices

- **Testing**: Always add tests for new features
- **Documentation**: Always update documentation for API changes
- **Type Safety**: Use TypeScript type definitions appropriately
- **Error Handling**: Appropriate error messages and log output

## 📈 Roadmap

### v1.1.0 (Planned)
- [ ] **Kubernetes Support**: Helm Chart provision
- [ ] **Metrics**: Prometheus/Grafana integration
- [ ] **Plugin System**: Custom rule extensions

### v1.2.0 (Planned)
- [ ] **AI Integration**: GPT-4 automatic rule generation
- [ ] **IDE Extension**: VS Code Extension
- [ ] **Cloud Support**: AWS/GCP/Azure support

### v2.0.0 (Planned)
- [ ] **Distributed Architecture**: Microservices
- [ ] **Real-time Collaboration**: WebSocket utilization
- [ ] **Machine Learning**: Rule recommendation system

## 🏆 Community

### Contributors

This project is supported by the following amazing contributors:

- [@AkitoSakurabaCreator](https://github.com/AkitoSakurabaCreator) - Project founder & maintainer

### Acknowledgments

- **Model Context Protocol**: Amazing protocol that enables standard AI agent integration
- **Go Community**: High-performance backend development environment
- **React Community**: Excellent frontend development experience
- **Docker**: Consistent development and production environment

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

## 🆘 Support & Community

### Issue Reports & Questions

- **🐛 Bug Reports**: [GitHub Issues](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/issues)
- **💡 Feature Proposals**: [GitHub Issues](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/issues)
- **❓ Questions & Discussions**: [GitHub Discussions](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/discussions)

### Documentation

- **📚 Detailed Documentation**: [GitHub Wiki](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server/wiki)
- **🎥 Tutorials**: [YouTube Playlist](https://youtube.com/playlist?list=...)
- **📖 API Reference**: [API Docs](https://api-docs.rulemcp.com)

### Community

- **💬 Discord**: [Rule MCP Server Community](https://discord.gg/...)
- **🐦 Twitter**: [@RuleMCPServer](https://twitter.com/RuleMCPServer)
- **📧 Mailing List**: [Google Groups](https://groups.google.com/g/rule-mcp-server)

---

**⭐ If this project is helpful, please give it a star on GitHub!**

**🚀 Let's achieve better code quality through AI agent integration!**
