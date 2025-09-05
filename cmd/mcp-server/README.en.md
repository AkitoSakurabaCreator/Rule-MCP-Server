# Rule MCP Server

[![npm version](https://badge.fury.io/js/rule-mcp-server.svg)](https://badge.fury.io/js/rule-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A MCP (Model Context Protocol) server that allows AI agents (Cursor, Cline) to retrieve and apply common coding rules.

**📦 npm package**: Published as [rule-mcp-server](https://www.npmjs.com/package/rule-mcp-server)

## Installation

### pnpm (Recommended)

```bash
pnpm add -g rule-mcp-server
```

### pnpm dlx (No Installation Required)

```bash
pnpm dlx rule-mcp-server
```

## Usage

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

### Environment Variables

- `RULE_SERVER_URL`: Rule MCP Server URL (default: http://localhost:18080)
- `MCP_API_KEY`: API key (optional, required for authentication)

## Available Tools

- `getRules`: Get project rules
- `validateCode`: Code validation
- `getProjectInfo`: Get project information
- `autoDetectProject`: Auto-detect project
- `scanLocalProjects`: Scan local projects
- `getGlobalRules`: Get global rules

## License

MIT License
