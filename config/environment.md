# Environment Variables Configuration

## Server Configuration

```bash
# Server Configuration
PORT=18081
HOST=0.0.0.0
ENVIRONMENT=development
LOG_LEVEL=info
```

## Database Configuration

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=15432
DB_NAME=rule_mcp_db
DB_USER=rule_mcp_user
DB_PASSWORD=your_secure_password
DB_SSL_MODE=disable
```

## Authentication Configuration

```bash
# Authentication Configuration
AUTH_ENABLED=true
AUTH_DEFAULT_ACCESS_LEVEL=public
AUTH_REQUIRE_HTTPS=false  # Set to true in production
AUTH_SESSION_TIMEOUT_MINUTES=480
AUTH_MAX_CONCURRENT_SESSIONS=5
```

## API Key Configuration

```bash
# API Key Configuration
API_KEY_ENABLED=true
API_KEY_EXPIRATION_DAYS=365
API_KEY_LENGTH=32
API_KEY_HASH_ALGORITHM=bcrypt
```

## Rate Limiting

```bash
# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_BURST_LIMIT=20
RATE_LIMIT_MCP_REQUESTS_PER_MINUTE=200
RATE_LIMIT_MCP_BURST_LIMIT=50
```

## Security Configuration

```bash
# Security Configuration
SECURITY_PASSWORD_MIN_LENGTH=8
SECURITY_PASSWORD_REQUIRE_UPPERCASE=true
SECURITY_PASSWORD_REQUIRE_LOWERCASE=true
SECURITY_PASSWORD_REQUIRE_NUMBERS=true
SECURITY_PASSWORD_REQUIRE_SPECIAL_CHARS=true
```

## Team Collaboration

```bash
# Team Collaboration
TEAM_MAX_MEMBERS_PER_PROJECT=50
TEAM_DEFAULT_VISIBILITY=public
```

## Audit Logging

```bash
# Audit Logging
AUDIT_LOG_ENABLED=true
AUDIT_LOG_RETENTION_DAYS=365
AUDIT_LOG_AUTH_ATTEMPTS=true
AUDIT_LOG_PERMISSION_CHANGES=true
AUDIT_LOG_RULE_MODIFICATIONS=true
```

## MCP Protocol Configuration

```bash
# MCP Protocol Configuration
MCP_ENABLED=true
MCP_DEFAULT_ACCESS_LEVEL=public
MCP_PUBLIC_METHODS=getRules,getProjectInfo,validateCode,getGlobalRules
MCP_PROTECTED_METHODS=createRule,updateRule,deleteRule,createProject,updateProject,deleteProject
```

## CORS Configuration

```bash
# CORS Configuration
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:13000,http://localhost:18000
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-API-Key
```

## Development Settings

```bash
# Development Settings
DEV_MODE=true
DEV_SKIP_AUTH=false
DEV_MOCK_DATA=false
DEV_DEBUG_MCP=false
```

## Production Configuration Example

```bash
# Production Configuration
ENVIRONMENT=production
PORT=18080
AUTH_REQUIRE_HTTPS=true
LOG_LEVEL=warn
DEV_MODE=false
DEV_SKIP_AUTH=false
```

## Docker Environment Variables

```bash
# Docker Compose Environment Variables
POSTGRES_PASSWORD=your_secure_password
POSTGRES_USER=rule_mcp_user
POSTGRES_DB=rule_mcp_db
POSTGRES_PORT=15432
```

## MCP Client Configuration

```bash
# MCP Client Environment Variables
MCP_SERVER_URL=http://localhost:18081
MCP_API_KEY=your_api_key_here
MCP_ACCESS_LEVEL=user
AUTO_INJECT=true
```
