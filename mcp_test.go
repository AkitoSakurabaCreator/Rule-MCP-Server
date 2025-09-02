package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestMCPServerHandleRequest(t *testing.T) {
	// テスト用のルールファイルを作成
	testRules := `{
		"mcp-test": {
			"project_id": "mcp-test",
			"rules": [
				{
					"id": "test-rule",
					"name": "Test Rule",
					"description": "A test rule",
					"type": "test",
					"severity": "warning",
					"pattern": "test_pattern",
					"message": "Test pattern detected"
				}
			]
		}
	}`

	err := os.WriteFile("rules.json", []byte(testRules), 0644)
	if err != nil {
		t.Fatalf("Failed to create test rules file: %v", err)
	}
	defer os.Remove("rules.json")

	server := NewMCPServer()
	ctx := context.Background()

	// getRules メソッドのテスト
	getRulesRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "getRules",
		Params: GetRulesParams{
			ProjectID: "mcp-test",
		},
	}

	requestData, err := json.Marshal(getRulesRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	responseData, err := server.HandleRequest(ctx, requestData)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	var response MCPResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got: %v", response.Error)
	}

	if response.ID != "1" {
		t.Errorf("Expected response ID '1', got '%s'", response.ID)
	}

	// validateCode メソッドのテスト
	validateRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      "2",
		Method:  "validateCode",
		Params: map[string]interface{}{
			"project_id": "mcp-test",
			"code":       "test_pattern",
		},
	}

	requestData, err = json.Marshal(validateRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	responseData, err = server.HandleRequest(ctx, requestData)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	err = json.Unmarshal(responseData, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got: %v", response.Error)
	}

	if response.ID != "2" {
		t.Errorf("Expected response ID '2', got '%s'", response.ID)
	}
}

func TestMCPServerInvalidMethod(t *testing.T) {
	server := NewMCPServer()
	ctx := context.Background()

	// 存在しないメソッドのテスト
	invalidRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      "3",
		Method:  "invalidMethod",
		Params:  nil,
	}

	requestData, err := json.Marshal(invalidRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	responseData, err := server.HandleRequest(ctx, requestData)
	if err != nil {
		t.Fatalf("Failed to handle request: %v", err)
	}

	var response MCPResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error for invalid method, got nil")
	}

	if response.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got %d", response.Error.Code)
	}
}

func TestMCPServerInvalidJSON(t *testing.T) {
	server := NewMCPServer()
	ctx := context.Background()

	// 無効なJSONのテスト
	invalidJSON := []byte(`{"invalid": json}`)

	_, err := server.HandleRequest(ctx, invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}
