package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type MCPServer struct {
	rules map[string]ProjectRules
}

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetRulesParams struct {
	ProjectID string `json:"project_id"`
}

type GetRulesResult struct {
	Rules []Rule `json:"rules"`
}

func NewMCPServer() *MCPServer {
	return &MCPServer{
		rules: make(map[string]ProjectRules),
	}
}

func (s *MCPServer) HandleRequest(ctx context.Context, requestData []byte) ([]byte, error) {
	var request MCPRequest
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	var response MCPResponse
	response.JSONRPC = "2.0"
	response.ID = request.ID

	switch request.Method {
	case "getRules":
		result, err := s.handleGetRules(request.Params)
		if err != nil {
			response.Error = &MCPError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			response.Result = result
		}
	case "validateCode":
		result, err := s.handleValidateCode(request.Params)
		if err != nil {
			response.Error = &MCPError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			response.Result = result
		}
	default:
		response.Error = &MCPError{
			Code:    -32601,
			Message: "Method not found",
		}
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseData, nil
}

func (s *MCPServer) handleGetRules(params interface{}) (*GetRulesResult, error) {
	paramsData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	var getRulesParams GetRulesParams
	if err := json.Unmarshal(paramsData, &getRulesParams); err != nil {
		return nil, fmt.Errorf("failed to parse params: %w", err)
	}

	rules, err := loadRules(getRulesParams.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}

	return &GetRulesResult{
		Rules: rules.Rules,
	}, nil
}

func (s *MCPServer) handleValidateCode(params interface{}) (*ValidationResult, error) {
	paramsData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	var validateParams struct {
		ProjectID string `json:"project_id"`
		Code      string `json:"code"`
	}

	if err := json.Unmarshal(paramsData, &validateParams); err != nil {
		return nil, fmt.Errorf("failed to parse params: %w", err)
	}

	result := validateCode(validateParams.ProjectID, validateParams.Code)
	return &result, nil
}

func (s *MCPServer) Start(ctx context.Context) error {
	log.Println("MCP Server started")

	select {
	case <-ctx.Done():
		log.Println("MCP Server stopped")
		return nil
	}
}
