package impl

import (
	"encoding/json"
	"fmt"

	mcpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/mcp"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPImpl struct {
	mcpcel.MCP
}

func (m *MCPImpl) Parse(content []byte) (*mcpcel.MCPRequest, error) {
	var jsonRPCReq mcp.JSONRPCRequest
	if err := json.Unmarshal(content, &jsonRPCReq); err != nil {
		return nil, err
	}

	mcpReq := &mcpcel.MCPRequest{
		Method: jsonRPCReq.Method,
		ID:     jsonRPCReq.ID,
	}

	switch mcpReq.Method {
	case string(mcp.MethodInitialize):
		params, err := unmarshalParams[mcp.InitializeParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.Initialize = params

	case string(mcp.MethodPing):
		break // Do nothing, no params

	// All *List are using Paginated
	case string(mcp.MethodResourcesList),
		string(mcp.MethodResourcesTemplatesList),
		string(mcp.MethodToolsList),
		string(mcp.MethodPromptsList):
		params, err := unmarshalParams[mcp.PaginatedParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.Paginated = params

	case string(mcp.MethodResourcesRead):
		params, err := unmarshalParams[mcp.ReadResourceParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.ResourceRead = params

	case string(mcp.MethodToolsCall):
		params, err := unmarshalParams[mcp.CallToolParams](content)
		if err != nil {
			return nil, err
		}
		fmt.Println("mcp.MethodToolsCall params", params)
		mcpReq.ToolCall = params

	case string(mcp.MethodSetLogLevel):
		params, err := unmarshalParams[mcp.SetLevelParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.SetLogLevel = params

	case string(mcp.MethodPromptsGet):
		params, err := unmarshalParams[mcp.GetPromptParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.PromptGet = params

	case string(mcp.MethodElicitationCreate):
		params, err := unmarshalParams[mcp.ElicitationParams](content)
		if err != nil {
			return nil, err
		}
		mcpReq.Elicitation = params
	}

	return mcpReq, nil
}

func unmarshalParams[T any](content []byte) (*T, error) {
	var enveloppe struct {
		Params T `json:"params"`
	}
	if err := json.Unmarshal(content, &enveloppe); err != nil {
		return nil, err
	}
	return &enveloppe.Params, nil
}
