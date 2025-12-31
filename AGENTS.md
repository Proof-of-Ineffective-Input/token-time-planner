# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Project Commands

- **Build**: `go build -ldflags="-s -w" -o token-time-planer.exe cmd/token-time-planer/main.go`
- **Run CLI**: `./token-time-planer -plan <path> [-rate <float>]`
- **Run MCP**: `./token-time-planer -mcp`
- **Release**: Push a tag matching `v*` to trigger GitHub Actions (builds for Windows, Linux, Darwin).

## Non-Obvious Patterns

- **Windows Build Lock**: You MUST empty [`.roo/mcp.json`](.roo/mcp.json:1) before building to release the file lock; restore it after compilation to remount the MCP server.
- **Context Accumulation**: Input tokens for each file include the `diff + output` of all preceding files in the plan.
- **Subtask Isolation**: `subtask: true` resets the context to 0 and applies a 1.2x time multiplier.
- **Buffer Priority**: `safety_rate` in `plan.yaml` overrides the `-rate` CLI flag.
- **TPS Scraper**: Metrics are fetched from `openrouter.ai/api/frontend/stats/throughput-comparison`. Defaults to 50 TPS on failure.
- **Pricing**: Fetched from `openrouter.ai/api/v1/models`.
- **Embedded Spec**: The specification is embedded from [`cmd/token-time-planer/SPEC.md`](cmd/token-time-planer/SPEC.md:1) and served via MCP resource `ttp://spec.md`.
