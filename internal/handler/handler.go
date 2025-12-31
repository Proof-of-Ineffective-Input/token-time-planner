package handler

import (
	"context"

	"ttp-mcp/internal/ttp"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CalculatePlanInput struct {
	PlanPath string `json:"plan_path" jsonschema:"description:path to plan.yaml"`
}

type CalculatePlanOutput struct {
	Report string `json:"report"`
}

func CalculatePlanHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculatePlanInput) (*mcp.CallToolResult, CalculatePlanOutput, error) {
	result, err := ttp.RunPlan(input.PlanPath, 1.8)
	if err != nil {
		return nil, CalculatePlanOutput{}, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result.Report,
			},
		},
	}, CalculatePlanOutput{Report: result.Report}, nil
}
