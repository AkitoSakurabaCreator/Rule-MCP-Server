#!/usr/bin/env node
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ErrorCode,
  ListResourcesRequestSchema,
  ListResourceTemplatesRequestSchema,
  ListToolsRequestSchema,
  McpError,
  ReadResourceRequestSchema,
} from '@modelcontextprotocol/sdk/types.js';
import axios from 'axios';

// 環境変数から設定を取得
const RULE_SERVER_URL = process.env.RULE_SERVER_URL || 'http://localhost:18080';
const MCP_API_KEY = process.env.MCP_API_KEY || '';

interface RuleServerResponse {
  id: string;
  result?: any;
  error?: {
    code: number;
    message: string;
  };
}

interface ValidationArgs {
  project_id: string;
  code: string;
  language?: string;
}

interface RuleArgs {
  project_id: string;
  language?: string;
}

interface ProjectInfoArgs {
  project_id: string;
}

interface AutoDetectArgs {
  path: string;
}

interface ScanProjectsArgs {
  base_path?: string;
}

const isValidRuleArgs = (args: any): args is RuleArgs =>
  typeof args === 'object' &&
  args !== null &&
  typeof args.project_id === 'string' &&
  (args.language === undefined || typeof args.language === 'string');

const isValidValidationArgs = (args: any): args is ValidationArgs =>
  typeof args === 'object' &&
  args !== null &&
  typeof args.project_id === 'string' &&
  typeof args.code === 'string' &&
  (args.language === undefined || typeof args.language === 'string');

const isValidProjectInfoArgs = (args: any): args is ProjectInfoArgs =>
  typeof args === 'object' &&
  args !== null &&
  typeof args.project_id === 'string';

const isValidAutoDetectArgs = (args: any): args is AutoDetectArgs =>
  typeof args === 'object' &&
  args !== null &&
  typeof args.path === 'string';

const isValidScanProjectsArgs = (args: any): args is ScanProjectsArgs =>
  typeof args === 'object' &&
  args !== null &&
  (args.base_path === undefined || typeof args.base_path === 'string');

class RuleMCPServer {
  private server: Server;
  private axiosInstance;

  constructor() {
    this.server = new Server(
      {
        name: 'rule-mcp-server',
        version: '1.0.0',
      },
      {
        capabilities: {
          resources: {},
          tools: {},
        },
      }
    );

    this.axiosInstance = axios.create({
      baseURL: RULE_SERVER_URL,
      headers: MCP_API_KEY ? { 'X-API-Key': MCP_API_KEY } : {},
      timeout: 30000,
    });

    this.setupResourceHandlers();
    this.setupToolHandlers();
    
    // エラーハンドリング
    this.server.onerror = (error) => console.error('[MCP Error]', error);
    process.on('SIGINT', async () => {
      await this.server.close();
      process.exit(0);
    });
  }

  private setupResourceHandlers() {
    // 静的リソース: プロジェクト一覧
    this.server.setRequestHandler(ListResourcesRequestSchema, async () => ({
      resources: [
        {
          uri: 'rule://projects/list',
          name: 'Available Projects',
          mimeType: 'application/json',
          description: 'List of all available projects with their rules',
        },
        {
          uri: 'rule://global-rules/javascript',
          name: 'JavaScript Global Rules',
          mimeType: 'application/json',
          description: 'Global coding rules for JavaScript projects',
        },
        {
          uri: 'rule://global-rules/typescript',
          name: 'TypeScript Global Rules',
          mimeType: 'application/json',
          description: 'Global coding rules for TypeScript projects',
        },
        {
          uri: 'rule://global-rules/go',
          name: 'Go Global Rules',
          mimeType: 'application/json',
          description: 'Global coding rules for Go projects',
        },
        {
          uri: 'rule://global-rules/python',
          name: 'Python Global Rules',
          mimeType: 'application/json',
          description: 'Global coding rules for Python projects',
        },
      ],
    }));

    // 動的リソーステンプレート
    this.server.setRequestHandler(
      ListResourceTemplatesRequestSchema,
      async () => ({
        resourceTemplates: [
          {
            uriTemplate: 'rule://{project_id}/rules',
            name: 'Project Rules',
            mimeType: 'application/json',
            description: 'Coding rules for a specific project',
          },
          {
            uriTemplate: 'rule://{project_id}/info',
            name: 'Project Information',
            mimeType: 'application/json',
            description: 'Information about a specific project',
          },
          {
            uriTemplate: 'rule://global-rules/{language}',
            name: 'Language Global Rules',
            mimeType: 'application/json',
            description: 'Global rules for a specific programming language',
          },
        ],
      })
    );

    // リソース読み取り
    this.server.setRequestHandler(
      ReadResourceRequestSchema,
      async (request) => {
        const uri = request.params.uri;

        try {
          if (uri === 'rule://projects/list') {
            // プロジェクト一覧を取得
            const response = await this.axiosInstance.get('/api/v1/projects');
            return {
              contents: [
                {
                  uri: request.params.uri,
                  mimeType: 'application/json',
                  text: JSON.stringify(response.data, null, 2),
                },
              ],
            };
          }

          // プロジェクト固有のルール
          const projectRulesMatch = uri.match(/^rule:\/\/([^/]+)\/rules$/);
          if (projectRulesMatch) {
            const projectId = decodeURIComponent(projectRulesMatch[1]);
            const response = await this.callRuleServer('getRules', {
              project_id: projectId,
            });
            
            return {
              contents: [
                {
                  uri: request.params.uri,
                  mimeType: 'application/json',
                  text: JSON.stringify(response, null, 2),
                },
              ],
            };
          }

          // プロジェクト情報
          const projectInfoMatch = uri.match(/^rule:\/\/([^/]+)\/info$/);
          if (projectInfoMatch) {
            const projectId = decodeURIComponent(projectInfoMatch[1]);
            const response = await this.callRuleServer('getProjectInfo', {
              project_id: projectId,
            });
            
            return {
              contents: [
                {
                  uri: request.params.uri,
                  mimeType: 'application/json',
                  text: JSON.stringify(response, null, 2),
                },
              ],
            };
          }

          // グローバルルール
          const globalRulesMatch = uri.match(/^rule:\/\/global-rules\/([^/]+)$/);
          if (globalRulesMatch) {
            const language = decodeURIComponent(globalRulesMatch[1]);
            const response = await this.axiosInstance.get(`/api/v1/global-rules/${language}`);
            
            return {
              contents: [
                {
                  uri: request.params.uri,
                  mimeType: 'application/json',
                  text: JSON.stringify(response.data, null, 2),
                },
              ],
            };
          }

          throw new McpError(
            ErrorCode.InvalidRequest,
            `Unknown resource URI: ${uri}`
          );
        } catch (error) {
          if (axios.isAxiosError(error)) {
            throw new McpError(
              ErrorCode.InternalError,
              `Rule server error: ${error.response?.data?.message || error.message}`
            );
          }
          throw error;
        }
      }
    );
  }

  private setupToolHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: [
        {
          name: 'getRules',
          description: 'Get coding rules for a specific project',
          inputSchema: {
            type: 'object',
            properties: {
              project_id: {
                type: 'string',
                description: 'The project ID to get rules for',
              },
              language: {
                type: 'string',
                description: 'Programming language (optional)',
              },
            },
            required: ['project_id'],
          },
        },
        {
          name: 'validateCode',
          description: 'Validate code against project rules',
          inputSchema: {
            type: 'object',
            properties: {
              project_id: {
                type: 'string',
                description: 'The project ID to validate against',
              },
              code: {
                type: 'string',
                description: 'The code to validate',
              },
              language: {
                type: 'string',
                description: 'Programming language (optional)',
              },
            },
            required: ['project_id', 'code'],
          },
        },
        {
          name: 'getProjectInfo',
          description: 'Get information about a specific project',
          inputSchema: {
            type: 'object',
            properties: {
              project_id: {
                type: 'string',
                description: 'The project ID to get info for',
              },
            },
            required: ['project_id'],
          },
        },
        {
          name: 'autoDetectProject',
          description: 'Automatically detect project from path and get appropriate rules',
          inputSchema: {
            type: 'object',
            properties: {
              path: {
                type: 'string',
                description: 'The path to detect project from',
              },
            },
            required: ['path'],
          },
        },
        {
          name: 'scanLocalProjects',
          description: 'Scan local directory to detect multiple projects',
          inputSchema: {
            type: 'object',
            properties: {
              base_path: {
                type: 'string',
                description: 'The base path to scan for projects (optional, defaults to current directory)',
              },
            },
          },
        },
        {
          name: 'getGlobalRules',
          description: 'Get global rules for a specific programming language',
          inputSchema: {
            type: 'object',
            properties: {
              language: {
                type: 'string',
                description: 'Programming language (javascript, typescript, go, python, etc.)',
              },
            },
            required: ['language'],
          },
        },
      ],
    }));

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      try {
        switch (request.params.name) {
          case 'getRules':
            if (!isValidRuleArgs(request.params.arguments)) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Invalid getRules arguments'
              );
            }
            return await this.handleGetRules(request.params.arguments);

          case 'validateCode':
            if (!isValidValidationArgs(request.params.arguments)) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Invalid validateCode arguments'
              );
            }
            return await this.handleValidateCode(request.params.arguments);

          case 'getProjectInfo':
            if (!isValidProjectInfoArgs(request.params.arguments)) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Invalid getProjectInfo arguments'
              );
            }
            return await this.handleGetProjectInfo(request.params.arguments);

          case 'autoDetectProject':
            if (!isValidAutoDetectArgs(request.params.arguments)) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Invalid autoDetectProject arguments'
              );
            }
            return await this.handleAutoDetectProject(request.params.arguments);

          case 'scanLocalProjects':
            if (!isValidScanProjectsArgs(request.params.arguments)) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Invalid scanLocalProjects arguments'
              );
            }
            return await this.handleScanLocalProjects(request.params.arguments);

          case 'getGlobalRules':
            if (!request.params.arguments?.language) {
              throw new McpError(
                ErrorCode.InvalidParams,
                'Language parameter is required'
              );
            }
            return await this.handleGetGlobalRules(request.params.arguments as { language: string });

          default:
            throw new McpError(
              ErrorCode.MethodNotFound,
              `Unknown tool: ${request.params.name}`
            );
        }
      } catch (error) {
        if (error instanceof McpError) {
          throw error;
        }
        if (axios.isAxiosError(error)) {
          return {
            content: [
              {
                type: 'text',
                text: `Rule server error: ${error.response?.data?.message || error.message}`,
              },
            ],
            isError: true,
          };
        }
        throw new McpError(
          ErrorCode.InternalError,
          `Unexpected error: ${error instanceof Error ? error.message : String(error)}`
        );
      }
    });
  }

  private async callRuleServer(method: string, params: any): Promise<any> {
    const response = await this.axiosInstance.post<RuleServerResponse>('/mcp/request', {
      id: `mcp-${Date.now()}`,
      method,
      params,
    });

    if (response.data.error) {
      throw new Error(response.data.error.message);
    }

    return response.data.result;
  }

  private async handleGetRules(args: RuleArgs) {
    try {
      const result = await this.callRuleServer('getRules', args);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to get rules: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async handleValidateCode(args: ValidationArgs) {
    try {
      const result = await this.callRuleServer('validateCode', args);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to validate code: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async handleGetProjectInfo(args: ProjectInfoArgs) {
    try {
      const result = await this.callRuleServer('getProjectInfo', args);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to get project info: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async handleAutoDetectProject(args: AutoDetectArgs) {
    try {
      const result = await this.callRuleServer('autoDetectProject', args);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to auto-detect project: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async handleScanLocalProjects(args: ScanProjectsArgs) {
    try {
      const result = await this.callRuleServer('scanLocalProjects', args);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(result, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to scan local projects: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async handleGetGlobalRules(args: { language: string }) {
    try {
      const response = await this.axiosInstance.get(`/api/v1/global-rules/${args.language}`);
      return {
        content: [
          {
            type: 'text',
            text: JSON.stringify(response.data, null, 2),
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to get global rules: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('Rule MCP Server running on stdio');
  }
}

const server = new RuleMCPServer();
server.run().catch(console.error);
