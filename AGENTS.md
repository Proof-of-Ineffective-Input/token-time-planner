# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Project Commands
- Build: `go build -ldflags="-s -w" -o token-time-planer.exe cmd/token-time-planer/main.go`
- Run CLI: `./token-time-planer -plan <path> [-buffer <float>]`
- Run MCP: `./token-time-planer -mcp`

## Non-Obvious Patterns
- **Windows Build Lock**: 编译前必须清空 [`.roo/mcp.json`](.roo/mcp.json:1) 以解除 Windows 文件锁定；编译完成后重新写入配置以挂载 MCP Server。
- **Context Accumulation**: Input tokens for each file include the `diff + output` of all preceding files in the plan.
- **Subtask Isolation**: Setting `subtask: true` resets the context to 0 and applies a 1.2x time multiplier.
- **Buffer Priority**: `safety_buffer` defined in `plan.yaml` silently overrides the `-buffer` CLI flag.
- **TPS Scraper**: Metrics are fetched from `openrouter.ai/api/frontend/stats/throughput-comparison`. If it fails, it defaults to 50 TPS.
- **Pricing**: Pricing is fetched from the standard OpenRouter models API.
