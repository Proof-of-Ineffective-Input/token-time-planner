# Token-Time Planner (TTP) 🚀

[English](README.md) | [中文](README_ZH.md)

**别再靠感觉估时了。** TTP 是一个专为 AI 驱动开发设计的 MCP 服务端。它将工程严谨性引入 Agent 工作流，用 **diff tokens** 和 **重生成倍率** 取代模糊的“大概几天”。

## 💡 为什么需要 TTP？

在 AI Agent 时代，开发的瓶颈不再是打字速度，而是模型的吞吐量 (TPS) 和上下文窗口管理。TTP 将开发过程建模为一系列 Token 生成事件，让你清晰掌握：

- ⏳ **实际耗时**：基于模型真实速度的预测。
- 💰 **API 成本**：精确到 Token 的费用预估。
- 🧠 **上下文压力**：跨文件操作时的 Token 累积情况。

## ✨ 核心特性

- **上下文感知**：自动计算跨文件操作时的输入 Token 累积。
- **实时指标**：从 OpenRouter 抓取最新的 TPS 和价格数据。
- **子任务隔离**：支持复杂重构中的上下文重置与时间惩罚建模。
- **安全缓冲**：内置安全系数，对冲 AI 幻觉和反复调试的时间成本。
- **原生集成**：通过 MCP 协议完美适配 Cursor, Windsurf 和 Roo Code。
- **内置规范**：通过 `ttp://spec.md` 资源直接访问完整的技术规范文档。

## 🛠️ `plan.yaml` 规范

TTP 驱动核心是一个声明式的 [`plan.yaml`](merge-plan.yaml:1)。

```yaml
plan:
  task_summary: "重构鉴权逻辑"
  model_id: "google/gemini-3-flash-preview"
  safety_rate: 1.5
  files:
    - path: "internal/auth/service.go"
      action: "modify"
      predicted_diff_tokens: 2500
      regen_rate: 2 # 你预计会反复生成几次？
      description: "更新 JWT 校验逻辑"
```

## 📐 计算逻辑（极简版）

我们基于以下公式推导：

1. **输入**：`当前上下文 + 预计 Diff`
2. **输出**：`预计 Diff × 重生成倍率`
3. **时长**：`输出 / TPS × 安全系数`

## 🚀 快速开始

### 1. 编译

```bash
go build -ldflags="-s -w" -o token-time-planer.exe cmd/token-time-planer/main.go
```

*注意：在 Windows 上，如果正在运行 MCP，编译前请清空 [`.roo/mcp.json`](.roo/mcp.json:1) 以解除文件锁定。*

### 2. 配置 MCP

在你的 MCP 设置中添加：

```json
{
  "mcpServers": {
    "token-time-planer": {
      "command": "~/ttp-mcp/token-time-planer.exe",
      "args": ["-mcp"]
    }
  }
}
```

## 📄 开源协议

采用 [GPLv3](LICENSE:1) 协议。
