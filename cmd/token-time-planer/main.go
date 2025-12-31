package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"ttp-mcp/internal/handler"
	"ttp-mcp/internal/ttp"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed SPEC.md
var specText string

func main() {
	planPath := flag.String("plan", "", "path to plan.yaml")
	safetyRate := flag.Float64("rate", 1.8, "safety multiplier")
	isMCP := flag.Bool("mcp", false, "run in MCP mode")
	flag.Parse()

	if *isMCP {
		runMCP()
		return
	}

	if *planPath == "" {
		fmt.Println("Usage: token-time-planer -plan <path> [-rate <float>]")
		fmt.Println("   or: token-time-planer -mcp")
		os.Exit(1)
	}

	result, err := ttp.RunPlan(*planPath, *safetyRate)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(result.Report)
}

func runMCP() {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "Token-Time Planner (TTP) Engine",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name: "calculate_plan",
		Description: `Evaluates software engineering progress based on token-time metrics. 
You MUST read the specification from "ttp://spec/plan" before generating the plan.yaml. 
Pass the path of your generated plan.yaml to this tool for precise duration and cost estimation.`,
	}, handler.CalculatePlanHandler)

	s.AddResource(&mcp.Resource{
		URI:      "ttp://spec/plan",
		Name:     "TTP Plan Specification",
		MIMEType: "text/markdown",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "ttp://spec/plan",
					MIMEType: "text/markdown",
					Text:     specText,
				},
			},
		}, nil
	})

	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
