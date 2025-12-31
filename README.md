# Token-Time Planner (TTP) 🚀

[English](README.md) | [中文](README_ZH.md)

**Stop guessing, start calculating.** TTP is an MCP server that brings engineering rigor to AI-driven development. It replaces vague "gut feelings" with precise metrics based on **diff tokens** and **regeneration rates**.

## 💡 Why TTP?

In the age of AI Agents, the bottleneck isn't how fast you can type—it's how fast the model can think and how much context it needs to carry. TTP models your development workflow as a sequence of token-generation events, giving you a realistic view of:

- ⏳ **How long** a task will actually take.
- 💰 **How much** it will cost in API credits.
- 🧠 **Context pressure** across multiple files.

## ✨ Key Features

- **Context-Aware**: Automatically tracks how tokens accumulate as you move through a project.
- **Live Metrics**: Pulls real-time TPS (Tokens Per Second) and pricing from OpenRouter.
- **Subtask Logic**: Model complex refactors with context resets and time penalties.
- **Safety Buffers**: Built-in multipliers to account for the "AI hallucination tax."
- **IDE Native**: Works seamlessly with Cursor, Windsurf, and Roo Code via MCP.

## 🛠️ The `plan.yaml`

TTP runs on a simple, declarative [`plan.yaml`](merge-plan.yaml:1). Define your steps, and let the engine do the math.

```yaml
plan:
  task_summary: "Refactor auth logic"
  model_id: "google/gemini-3-flash-preview"
  safety_rate: 1.5
  files:
    - path: "internal/auth/service.go"
      action: "modify"
      predicted_diff_tokens: 2500
      regen_rate: 2 # How many times will you iterate?
      description: "Update JWT validation"
```

## 📐 The Math (Simplified)

We use a few core principles to calculate your plan:

1. **Input**: `Context + New Diff`
2. **Output**: `Diff × Regeneration Rate`
3. **Time**: `Output / TPS × Safety Multipliers`

## 🚀 Getting Started

### 1. Build

```bash
go build -ldflags="-s -w" -o token-time-planer.exe cmd/token-time-planer/main.go
```

*Note: If using Roo Code, clear [`.roo/mcp.json`](.roo/mcp.json:1) before building to release the file lock.*

### 2. Configure MCP

Add this to your settings:

```json
{
  "mcpServers": {
    "token-time-planer": {
      "command": "E:/Dev/ttp-mcp/token-time-planer.exe",
      "args": ["-mcp"]
    }
  }
}
```

## 📄 License

Licensed under [GPLv3](LICENSE:1).
