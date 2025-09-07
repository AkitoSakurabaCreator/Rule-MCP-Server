# Rule MCP Server

[![npm version](https://badge.fury.io/js/rule-mcp-server.svg)](https://badge.fury.io/js/rule-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A MCP (Model Context Protocol) server that allows AI agents (Cursor, Claude Code, Cline) to retrieve and apply common coding rules.

**📦 npm package**: Published as [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

## 🚀 Quick Start

### Pattern 1: Server Already Running

If the Rule MCP Server is already running, you only need to configure your AI agent.

#### 1. Installation

```bash
# Via pnpm dlx (recommended, no installation required)
pnpm dlx rule-mcp-server

# Or global installation
pnpm add -g rule-mcp-server
```

#### 2. AI Agent Configuration

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
```

#### 3. Start Using!

Restart your AI agent (Cursor/Claude Code) and it will automatically retrieve and apply coding rules.

---

### Pattern 2: Set Up Your Own Server

If you want to set up and run your own Rule MCP Server, please refer to the [main repository](https://github.com/AkitoSakurabaCreator/Rule-MCP-Server).

## Configuration Examples

### Cursor Configuration

Add the following to `~/.cursor/mcp.json`:

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

### Claude Code Configuration

```bash
# Add MCP server to Claude Code (stdio)
claude mcp add rule-mcp-server --env RULE_SERVER_URL=http://localhost:18080 -- pnpm dlx rule-mcp-server

# With API key
claude mcp add rule-mcp-server \
  --env RULE_SERVER_URL=http://localhost:18080 \
  --env MCP_API_KEY=your_api_key \
  -- pnpm dlx rule-mcp-server

# Reference: Anthropic official docs
# https://docs.anthropic.com/ja/docs/claude-code/mcp
```

### Environment Variables

- `RULE_SERVER_URL`: Rule MCP Server URL (default: http://localhost:18080)
- `MCP_API_KEY`: API key (optional, required for authentication)

Note: `MCP_API_KEY` is optional (Public access works without it). Set it only for team operations or when using management APIs.

## Prerequisite

The MCP client configuration assumes the Rule MCP Server is running.

```bash
curl http://localhost:18080/api/v1/health
```

If the server is not running, start it with Docker:

```bash
docker compose up -d
```

For LAN operation, change `RULE_SERVER_URL` to the LAN host IP.

## Available Tools

- `getRules`: Get project rules
- `validateCode`: Code validation
- `getProjectInfo`: Get project information
- `autoDetectProject`: Auto-detect project
- `scanLocalProjects`: Scan local projects
- `getGlobalRules`: Get global rules

## License

MIT License
