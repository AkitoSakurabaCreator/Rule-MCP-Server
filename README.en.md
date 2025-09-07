# Rule MCP Server

[![npm version](https://badge.fury.io/js/rule-mcp-server.svg)](https://badge.fury.io/js/rule-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Model Context Protocol (MCP) server that enables AI agents (Cursor, Claude Code, Cline) to retrieve and apply common coding rules across projects.

## Features

- **Project-specific rule management**
- **Language-specific global rule management**
- **Configurable global rule application per project**
- **Code rule violation validation**
- **Rule distribution via MCP protocol**
- **RESTful API endpoints**
- **React-based Web UI** for rule management
- **Multi-language support (i18n)**: English, Japanese, Chinese, Hindi, Spanish, Arabic
- **Dark theme support**: Light/Dark mode toggle
- **Grafana monitoring dashboard**: System metrics, MCP request monitoring, alerting
- **Clean architecture** for high maintainability

## 🚀 Quick Start

### Option 1: Using Existing Server (Recommended)

If a Rule MCP Server is already running, simply configure your AI agent:

#### 1. Install MCP Server

```bash
# Via pnpm dlx (recommended - no installation required)
pnpm dlx rule-mcp-server

# Or global installation
pnpm add -g rule-mcp-server
```

#### 2. Configure AI Agent

##### Cursor
```bash
# Copy configuration template
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json
```

##### Claude Code
```bash
# Add MCP server to Claude Code (stdio)
claude mcp add rule-mcp-server --env RULE_SERVER_URL=http://localhost:18080 -- pnpm dlx rule-mcp-server

# With API key
claude mcp add rule-mcp-server \
  --env RULE_SERVER_URL=http://localhost:18080 \
  --env MCP_API_KEY=your_api_key \
  -- pnpm dlx rule-mcp-server

# Reference: Anthropic official documentation
# https://docs.anthropic.com/en/docs/claude-code/mcp
```

#### 3. Start Using!

Restart your AI agent (Cursor/Claude Code) and you can now automatically retrieve and apply coding rules.

**📦 npm package**: [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

---

### Option 2: Self-Hosted Server

To run your own Rule MCP Server for production use:

#### 1. Environment Setup

```bash
# Create production environment file
cp env.production.example .env.production

# Edit .env.production file as needed
nano .env.production
```

**Required production settings:**
- `JWT_SECRET`: Strong random string (generation method below)
- `ALLOWED_ORIGINS`: Comma-separated allowed origins
- `ENV=production`: Set to production mode

**Secret key generation:**
```bash
# Using OpenSSL
openssl rand -hex 32

# Using Python
python3 -c "import secrets; print(secrets.token_hex(32))"

# Using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
```

#### 2. Start Server

##### Production Environment (Recommended)
```bash
# Start with production Docker Compose
docker compose -f docker-compose.prod.yml up -d

# Verify operation
curl http://localhost:18080/api/v1/health
# -> {"status":"ok"} if running
```

**Production access:**
- Web UI: http://localhost:13000
- API: http://localhost:18080/api/v1
- Database: localhost:15432

##### Development Environment (For Developers)
```bash
# Create development environment file
cp env.development.example .env.development

# Start with development Docker Compose
docker compose -f docker-compose.dev.yml up -d

# Verify operation
curl http://localhost:18080/api/v1/health
# -> {"status":"ok"} if running
```

**Development access:**
- Web UI: http://localhost:13000
- API: http://localhost:18080/api/v1
- Database: localhost:15432

##### Team Usage
```bash
# Start in production mode (for team sharing)
docker compose -f docker-compose.prod.yml up -d

# Verify operation
curl http://localhost:18080/api/v1/health
# -> {"status":"ok"} if running
```

**Team access:**
- Web UI: http://[server-IP]:13000
- API: http://[server-IP]:18080/api/v1
- Database: [server-IP]:15432

#### LAN Sharing Example (Team Usage)
- Start server on LAN host and set client environment variables to LAN IP:
```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "env": {
        "RULE_SERVER_URL": "http://192.168.1.20:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      }
    }
  }
}
```

Reference: Use Makefile with `make docker-up` / `make docker-down`

## Technology Stack

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

- Go 1.21 or higher installed
- Node.js 18 or higher installed
- Docker and Docker Compose installed
- PostgreSQL 15 or higher (automatically installed with Docker)

### Initial Admin Account

The following initial admin account is automatically created on first system startup:

- **Username**: `admin`
- **Password**: `admin123`
- **Email**: `admin@rulemcp.com`
- **Role**: Administrator (admin)

**Important**: Change the password after first login.

### Developer Setup

```bash
git clone https://github.com/AkitoSakurabaCreator/Rule-MCP-Server.git
cd Rule-MCP-Server

# Copy environment file
cp env.template .env.production

# Backend dependencies
go mod tidy

# Frontend dependencies
cd frontend
npm install
cd ..
```

### Startup Methods

#### Production Environment (Recommended)

```bash
# Start with production Docker Compose
docker compose -f docker-compose.prod.yml up -d

# Stop
docker compose -f docker-compose.prod.yml down
```

#### Development Environment

```bash
# Start with development Docker Compose
docker compose -f docker-compose.dev.yml up -d

# Stop
docker compose -f docker-compose.dev.yml down
```

#### Local Development (Without Docker)

```bash
# Backend (port 18080)
go run ./cmd/server

# Frontend (separate terminal)
cd frontend
npm start
```

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

This server operates as an MCP server, allowing direct rule retrieval from Cursor and Cline.

### MCP Endpoints

```
POST /mcp/request    # HTTP MCP requests
GET  /mcp/ws         # WebSocket MCP connection
```

### MCP Methods

- **`getRules`**: Retrieve project rules
- **`validateCode`**: Validate code for rule violations
- **`getProjectInfo`**: Get project information

### Standard MCP Server Configuration (Recommended)

This project provides a **standard MCP (Model Context Protocol) server**.

#### **Features**
- ✅ **Standard MCP SDK**: Full compliance using `@modelcontextprotocol/sdk`
- ✅ **StdioServerTransport**: Compatible with standard MCP clients
- ✅ **pnpm package support**: Easy installation with `pnpm dlx`
- ✅ **Docker support**: Stable operation in production
- ✅ **Template configuration**: Easy setup

#### **1. Install MCP Server**

##### **Via pnpm (Recommended)**
```bash
# Global installation
pnpm add -g rule-mcp-server

# Or via pnpm dlx (no installation required)
pnpm dlx rule-mcp-server
```

**📦 npm package**: Published as [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

##### **Development Build**
```bash
# Install dependencies
make install-mcp

# Build MCP server
make build-mcp
```

#### **2. Environment-Specific Configuration**

##### **Using pnpm Package (Recommended)**
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
- **Claude Code**: Add via CLI (`claude mcp add ...`)

#### **4. Available Tools**
The standard MCP server provides the following tools:

| Tool Name            | Description                    | Required Parameters     |
|---------------------|--------------------------------|------------------------|
| `getRules`          | Get project rules              | `project_id`           |
| `validateCode`      | Validate code                  | `project_id`, `code`   |
| `getProjectInfo`    | Get project information        | `project_id`           |
| `autoDetectProject` | Auto-detect project            | `path`                 |
| `scanLocalProjects` | Scan local projects            | `base_path` (optional) |
| `getGlobalRules`    | Get global rules               | `language`             |

#### **5. Available Resources**
The standard MCP server provides the following resources:

| Resource URI                      | Description                   |
|----------------------------------|-------------------------------|
| `rule://projects/list`           | Project list                  |
| `rule://{project_id}/rules`      | Project-specific rules        |
| `rule://{project_id}/info`       | Project information           |
| `rule://global-rules/{language}` | Language-specific global rules|

### Legacy HTTP Configuration (Compatibility)

For using the legacy HTTP API:

```json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "-H", "Content-Type: application/json",
        "-d", "{\"id\":\"${requestId}\",\"method\":\"${method}\",\"params\":${params}}",
        "http://localhost:18081/mcp/request"
      ],
      "env": {
        "MCP_SERVER_URL": "http://localhost:18081"
      },
      "description": "Rule MCP Server for AI agents to retrieve and apply common rules"
    }
  }
}
```

## MCP (Model Context Protocol) Features

### Basic Features

- **`getRules`**: Retrieve rules by project ID and optional language
- **`validateCode`**: Validate code for rule violations
- **`getProjectInfo`**: Get project information

### Project Auto-Detection Feature 🆕

Advanced feature that allows AI agents to **automatically recognize projects** and retrieve appropriate rules.

#### **Auto-Detection Priority**

1. **Directory name-based detection** (95% confidence)
   - Search using directory name as project ID
   - Excluded directories: `node_modules`, `vendor`, `dist`, `build`, `target`, `.git`, `.vscode`

2. **Git repository name detection** (90% confidence)
   - Parse origin URL from `.git/config`
   - SSH format: `git@github.com:username/repo-name.git`
   - HTTPS format: `https://github.com/username/repo-name.git`

3. **Language-specific file detection** (85% confidence)
   - `go.mod` → Go project
   - `package.json` → Node.js project
   - `requirements.txt` → Python project
   - `pom.xml` → Java project
   - `Cargo.toml` → Rust project
   - `composer.json` → PHP project
   - `Gemfile` → Ruby project

4. **Default project** (70% confidence)
   - Fallback when detection fails

#### **New MCP Methods**

##### **`autoDetectProject`**
Automatically detects project from specified path.

```json
{
  "id": "auto-detect",
  "method": "autoDetectProject",
  "params": {
    "path": "/path/to/your/project"
  }
}
```

**Response example:**
```json
{
  "id": "auto-detect",
  "result": {
    "project": {
      "project_id": "web-app",
      "name": "Web Application",
      "language": "javascript"
    },
    "rules": [...],
    "detection_method": "directory_name",
    "confidence": 0.95,
    "message": "Project detected from directory name 'web-app'"
  }
}
```

##### **`scanLocalProjects`**
Recursively scans local directory to detect projects.

```json
{
  "id": "scan-local",
  "method": "scanLocalProjects",
  "params": {
    "base_path": "/home/user/projects"
  }
}
```

**Response example:**
```json
{
  "id": "scan-local",
  "result": {
    "projects": [
      {
        "project": {...},
        "rules": [...],
        "detection_method": "git_repository",
        "confidence": 0.90,
        "message": "Project detected from Git repository name"
      }
    ],
    "count": 3
  }
}
```

#### **Usage Examples**

##### **Cursor/Cline Configuration**

###### **1. cursor-mcp-config.json (Recommended)**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "localhost",
        "DB_PORT": "15432",
        "DB_USER": "rule_mcp_user",
        "DB_PASSWORD": "rule_mcp_password",
        "DB_NAME": "rule_mcp_db"
      }
    }
  }
}
```

**File placement:**
- **Cursor**: `~/.cursor/mcp-servers/rule-mcp.json`
- **Cline**: `~/.cline/mcp-servers/rule-mcp.json`

###### **2. mcp-client-config.json (Complete version)**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "localhost",
        "DB_PORT": "15432",
        "DB_USER": "rule_mcp_user",
        "DB_PASSWORD": "rule_mcp_password",
        "DB_NAME": "rule_mcp_db"
      },
      "cwd": "/path/to/your/RuleMCPServer"
    }
  }
}
```

###### **3. No environment variables (JSON file mode)**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081"
      }
    }
  }
}
```

##### **Configuration Explanation**

| Parameter | Description             | Required | Default                    |
|-----------|-------------------------|----------|----------------------------|
| `command` | Command to execute      | ✅        | `go`                       |
| `args`    | Command arguments       | ✅        | `["run", "./cmd/server"]`  |
| `env`     | Environment variables   | ❌        | None                       |
| `cwd`     | Working directory       | ❌        | Current directory          |

##### **Environment-Specific Configuration Examples**

###### **Development Environment (JSON file mode)**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081"
      }
    }
  }
}
```

###### **Production Environment (PostgreSQL connection)**
```json
{
  "mcpServers": {
    "rule-mcp": {
      "command": "go",
      "args": ["run", "./cmd/server"],
      "env": {
        "PORT": "18081",
        "DB_HOST": "your-db-host",
        "DB_PORT": "5432",
        "DB_USER": "your-db-user",
        "DB_PASSWORD": "your-db-password",
        "DB_NAME": "your-db-name"
      }
    }
  }
}
```

##### **Configuration File Priority**

1. **Project-specific**: `./cursor-mcp-config.json`
2. **User settings**: `~/.cursor/mcp-servers/rule-mcp.json`
3. **Global settings**: `~/.cursor/mcp-servers/rule-mcp.json`

##### **Troubleshooting**

###### **Common Issues and Solutions**

| Issue                      | Cause                | Solution                    |
|----------------------------|----------------------|-----------------------------|
| `Method not found`         | Old MCP handler      | Restart server              |
| `Connection refused`       | Port in use          | Check with `lsof -i :18081` |
| `Database connection failed` | DB connection misconfig | Check environment variables |
| `Permission denied`        | File permissions     | Grant execute permission with `chmod +x` |

##### **AI Agent Auto-Detection**
AI agents can automatically recognize projects as follows:

1. **Automatic working directory detection**
2. **Appropriate rule set retrieval**
3. **Language-specific rule application**
4. **Project-specific rule application**

This allows AI agents to **always apply appropriate rules** for code review and suggestions **without understanding context**.

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

- **Light theme**: Bright and readable
- **Dark theme**: Eye-friendly night mode
- **Automatic setting save**: Saved to browser localStorage

### Web UI Features

- **Project management**: Create, edit, delete
- **Rule management**: Project-specific rules, global rules
- **Code validation**: Real-time rule violation checking
- **Responsive design**: Mobile and tablet support

## 📊 Monitoring & Metrics

### Grafana Monitoring Dashboard

Rule MCP Server includes a comprehensive monitoring system.

#### Monitored Metrics

1. **Nginx Request Rate**
   - HTTP request processing speed
   - Request count by status code

2. **Nginx Response Time**
   - 95th and 50th percentile response times

3. **System Statistics**
   - Total users
   - Total projects
   - Total rules

4. **MCP Request Rate**
   - Request count by MCP method
   - Success/error statistics

5. **MCP Response Time**
   - MCP request processing time

6. **Active Sessions & API Keys**
   - Currently active session count
   - Active API key count

7. **System Load**
   - CPU usage

8. **MCP Success & Error Rate**
   - Request success and error rates

#### Access Methods

##### Development Environment
- **Grafana**: http://localhost:14000
  - Username: `admin`
  - Password: `admin123`
- **Prometheus**: http://localhost:19090

##### Production Environment
- **Grafana**: http://localhost:14000
  - Username: Environment variable `GRAFANA_ADMIN_USER` (default: `admin`)
  - Password: Environment variable `GRAFANA_ADMIN_PASSWORD` (default: `admin123`)
- **Prometheus**: http://localhost:19090

#### Alert Configuration

The following alerts are configured:

- **High Error Rate**: When MCP request error rate exceeds 10%
- **High Response Time**: When 95th percentile response time exceeds 1 second
- **High System Load**: When system load exceeds 80%
- **Database Connection Failure**: When PostgreSQL connection is lost
- **Nginx High Error Rate**: When 5xx error rate exceeds 5%
- **Disk Space Low**: When disk usage exceeds 90%

#### Customization

##### Dashboard Editing
1. Login to Grafana
2. Open the dashboard
3. Click the "Settings" button in the top right
4. Select "JSON Model" to edit the JSON

##### Adding New Metrics
1. Add metrics to `internal/interface/handler/metrics_handler.go`
2. Add panels to the dashboard
3. Configure Prometheus queries

##### Adding Alert Rules
1. Add rules to `prometheus/rules/rule-mcp-alerts.yml`
2. Restart Prometheus

For details, see [Grafana README](grafana/README.md).

## Rule Definition

### Rule Creation Items

1. **Rule ID**: Unique rule identifier (e.g., `no-console-log`)
2. **Name**: Rule display name (e.g., `No Console Log`)
3. **Description**: Rule content, purpose, reason (e.g., `Console.log statements should not be in production code`)
4. **Type**: Rule category (naming, formatting, security, performance, etc.)
5. **Severity**: Importance level (error, warning, info)
6. **Pattern**: Regular expression pattern to detect (e.g., `console\.log`)
7. **Message**: Fix instruction for violations (e.g., `Console.log detected. Use proper logging framework in production.`)
8. **Active**: Rule enable/disable

### Sample Rule

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

### Startup Examples

#### Local Development

```bash
# Backend (development environment)
PORT=18081 go run ./cmd/server

# Frontend
cd frontend && npm start

# Using Makefile
make run           # Backend (port 18081)
make run-frontend  # Frontend
make run-prod      # Production environment (port 18080)
```

#### Using Docker

```bash
# Build Docker image
make docker-build

# Start services
make docker-up

# Stop services
make docker-down

# Check logs
make docker-logs

# Clean up resources
make docker-clean
```

## Development

### Testing

```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend && npm test
```

### Building

```bash
# Backend
go build -o rule-mcp-server ./cmd/server

# Frontend
cd frontend && npm run build
```

### Code Quality

```bash
# Go language quality check
make glb

# Full repository quality check
make glb_repo
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Error

```bash
# Restart database container
docker-compose restart postgres

# Complete database recreation
docker-compose down -v
docker-compose up -d postgres
```

#### 2. Port Conflict

```bash
# Check ports in use
lsof -i :18081
lsof -i :13000

# Start on different port
PORT=18082 go run ./cmd/server
```

#### 3. MCP Server Not Responding

```bash
# Check server startup
curl http://localhost:18081/api/v1/health

# Test MCP endpoint
curl -X POST http://localhost:18081/mcp/request \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"getRules","params":{"project_id":"web-app"}}'
```

#### 4. Frontend Build Error

```bash
# Reinstall dependencies
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### 5. Authentication Error

```bash
# Check API key
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/projects

# Check permission level
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/auth/me
```

### Log Checking

```bash
# Backend logs
docker logs rule-mcp-server

# Database logs
docker logs rule-mcp-postgres

# Frontend logs
cd frontend && npm start
```

## Permission Management System

### Access Levels

Rule MCP Server provides three access levels:

#### **Public (Open)**
- **Authentication**: Not required
- **Permissions**: View public rules/projects, code validation
- **Use case**: Personal use, open source projects
- **MCP access**: No restrictions (with rate limiting)

#### **User**
- **Authentication**: API key required
- **Permissions**: Create, edit, delete personal rules/projects
- **Use case**: Individual developers, small teams
- **MCP access**: No restrictions

#### **Admin**
- **Authentication**: Admin API key required
- **Permissions**: Full permissions (user management, global rule management)
- **Use case**: Team leaders, system administrators
- **MCP access**: No restrictions

### Authentication Methods

#### **API Key Authentication**
```bash
# Header authentication
curl -H "X-API-Key: your_api_key" http://localhost:18081/api/v1/projects

# MCP request authentication
curl -X POST http://localhost:18081/mcp/request \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"createRule","params":{...}}'
```

#### **Session Authentication**
```bash
# Login
curl -X POST http://localhost:18081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password"}'

# Session authentication
curl -H "Cookie: session=session_token" http://localhost:18081/api/v1/projects
```

### Team Collaboration Features

#### **Project Member Management**
```bash
# Add team member
curl -X POST http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: admin_key" \
  -H "Content-Type: application/json" \
  -d '{"username":"developer1","role":"member"}'

# Check team member permissions
curl http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: user_key"
```

#### **Project Visibility**
- **Public**: Viewable by all users
- **Team**: Viewable only by team members
- **Private**: Viewable only by project owner

#### **Permission Granularity**
```json
{
  "permissions": {
    "read": true,    // View rules/projects
    "write": true,   // Create/edit rules/projects
    "delete": false, // Delete rules/projects
    "admin": false   // Member management
  }
}
```

## Configuration Files

### **Authentication Configuration** (`config/auth.yaml`)
Defines permission levels, rate limiting, security settings

### **Environment Variables** (`config/environment.md`)
List of environment variables for server, database, authentication, MCP configuration

### **MCP Client Configuration**
- **`config/simple-mcp-config.json`**: Simple MCP configuration (no authentication, beginner-friendly)
- **`config/mcp-client-config.json`**: Complete MCP configuration (authentication and team features)

### **Database Schema** (`init.sql`)
Complete schema including permission management tables, user management, team collaboration features

## Security Features

### **Rate Limiting**
- **Public**: 50 req/min (limited)
- **User**: 100 req/min
- **Admin**: 200 req/min

### **Audit Logging**
- Authentication attempt records
- Permission change records
- Rule change records
- Retention period: 365 days

### **Password Policy**
- Minimum 8 characters
- Uppercase, lowercase, numbers, special characters required
- Session timeout: 8 hours

### **API Key Security**
- bcrypt hashing
- Expiration settings
- HTTPS required (production)
- Usage logging

## Usage Examples

### **Personal Use (Public)**
```bash
# Get rules without authentication
curl http://localhost:18081/api/v1/rules?project_id=web-app

# Code validation via MCP
curl -X POST http://localhost:18081/mcp/request \
  -H "Content-Type: application/json" \
  -d '{"id":"test","method":"validateCode","params":{"project_id":"web-app","code":"console.log(\"test\")"}}'
```

### **Team Use (User/Admin)**
```bash
# Create rule with API key
curl -X POST http://localhost:18081/api/v1/rules \
  -H "X-API-Key: user_key" \
  -H "Content-Type: application/json" \
  -d '{"project_id":"team-project","rule_id":"no-todo","name":"No TODO","type":"quality","severity":"warning","pattern":"TODO:","message":"TODO comment detected"}'

# Team member management
curl -X POST http://localhost:18081/api/v1/projects/team-project/members \
  -H "X-API-Key: admin_key" \
  -H "Content-Type: application/json" \
  -d '{"username":"new_dev","role":"member"}'
```

### **MCP Client Configuration**

#### Cursor
```json
// ~/.cursor/mcp.json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": [
        "dlx",
        "rule-mcp-server"
      ],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      },
      "description": "Standard MCP Server for Rule Management - provides coding rules and validation tools for AI agents",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

#### Claude Desktop
```json
// ~/Library/Application Support/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "rule-mcp-server": {
      "command": "pnpm",
      "args": [
        "dlx",
        "rule-mcp-server"
      ],
      "env": {
        "RULE_SERVER_URL": "http://localhost:18080",
        "MCP_API_KEY": "${MCP_API_KEY:-}"
      },
      "description": "Standard MCP Server for Rule Management - provides coding rules and validation tools for AI agents",
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

Note: `MCP_API_KEY` works without setting (Public access). Set only when using team operations or management APIs.

## 🚀 Quick Start

### Using pnpm Package (Recommended)

```bash
# 1. Install MCP server (or auto-install via pnpm dlx)
pnpm add -g rule-mcp-server

# 2. Create configuration file
cp config/pnpm-mcp-config.template.json ~/.cursor/mcp.json

# 3. Start using with AI agents (Cursor/Claude Desktop/Cline)!
```

**📦 npm package**: Published as [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

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

# 2. Start server
make run        # Backend (port 18081)
make run-frontend  # Frontend (port 3000)

# 3. Build MCP server
make install-mcp && make build-mcp

# 4. Create configuration file
cp config/standard-mcp-config.template.json ~/.cursor/mcp_settings.json
```

## 🌟 Key Features

### ✅ Standard MCP Compliance

- **Full MCP compatibility**: Using `@modelcontextprotocol/sdk`
- **Standard tools and resources**: Perfect integration with AI agents
- **StdioServerTransport**: Standard communication protocol

### ✅ Rich Features

- **6 MCP tools**: Rule retrieval, code validation, project auto-detection, etc.
- **5 MCP resources**: Project list, rule information, global rules
- **Project auto-detection**: AI automatically applies appropriate rules

### ✅ Production Ready

- **Docker environment**: Stable production operation
- **PostgreSQL**: High-performance database
- **Authentication & Authorization**: API keys, session authentication
- **Multi-language support**: 6 language support (Japanese, English, Chinese, etc.)

### ✅ Developer Friendly

- **Clean architecture**: High maintainability design
- **Web UI**: React-based management interface
- **Template configuration**: Easy setup
- **Rich documentation**: Detailed usage instructions

## 🎯 Use Cases

### Individual Developers
```bash
# Easy use without authentication
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

# 4. Start development server
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

### Community

- **💬 Discord**: [Rule MCP Server Community](https://discord.gg/dCAUC8m6dw)
- **🐦 X (formerly Twitter)**: [@_sakuraba_akito](https://x.com/_sakuraba_akito)

## 📋 Contribution Guidelines

### Contribution Flow

1. **Create Issue**: Bug reports or feature proposals
2. **Fork**: Fork the repository
3. **Create Branch**: `feature/your-feature` or `fix/your-fix`
4. **Development**: Code changes and test additions
5. **Pull Request**: Submit with detailed description

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📜 License

MIT License - See [LICENSE](LICENSE) file for details.

## 🔒 Security

If you discover security issues, please report them following the instructions in [SECURITY.md](SECURITY.md).

## 📝 Code of Conduct

This project follows [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## 📈 Changelog

See [CHANGELOG.md](CHANGELOG.md) for the latest changelog.

---

**⭐ If this project is helpful, please give it a star on GitHub!**

**🚀 Let's achieve better code quality through AI agent integration!**

## GitHub Repository

This project is hosted on GitHub: [https://github.com/AkitoSakurabaCreator/Rule-MCP-Server](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server)

### Repository Information

- **Repository**: [AkitoSakurabaCreator/Rule-MCP-Server](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server)
- **License**: MIT
- **Language**: TypeScript (46.3%), Go (40.6%), HTML (4.5%), JavaScript (4.3%), MDX (2.7%), CSS (0.8%), Other (0.8%)
- **Stars**: 0
- **Forks**: 0
- **Issues**: 0
- **Pull Requests**: 0

### Getting Started with GitHub

1. **Star the repository** if you find it useful
2. **Fork the repository** to contribute
3. **Create issues** for bug reports or feature requests
4. **Submit pull requests** for contributions
5. **Join discussions** for questions and ideas

### Development

- **Clone**: `git clone https://github.com/AkitoSakurabaCreator/Rule-MCP-Server.git`
- **Contribute**: Follow the contribution guidelines
- **Report Issues**: Use GitHub Issues
- **Discuss**: Use GitHub Discussions

### Community

- **Discord**: [Rule MCP Server Community](https://discord.gg/dCAUC8m6dw)
- **Twitter**: [@_sakuraba_akito](https://x.com/_sakuraba_akito)
- **GitHub**: [@AkitoSakurabaCreator](https://github.com/AkitoSakurabaCreator)