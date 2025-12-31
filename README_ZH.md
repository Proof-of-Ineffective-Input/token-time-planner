# Token-Time Planner (TTP)

> 让 AI 彻底摆脱"人类时间感"，用 token 量精准估算编码任务的时长与成本。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://go.dev/)
[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io)

---

## 📖 背景

在 **Token 时代**，实际编码效率由输出速度（TPS）和上下文增长决定。TTP 通过结构化规划，让 AI 工作流回归理性预期。

---

## ✨ 核心特性

### 📊 结构化规划（plan.yaml）
- **拓扑排序**：任务按 `files` 列表顺序执行，默认保留上下文。
- **上下文感知**：区分普通任务与隔离子任务（Subtask）。
- **模型绑定**：每个任务明确分配 `model_id`，用于精确 TPS 匹配。

### 🧮 外部计算器工具（calculate_plan）
- **动态 TPS 抓取**：自动从 OpenRouter 抓取实时 TPS，若抓取失败则回退至默认 **50 tok/s**。
- **累加成本计算**：考虑上下文随执行顺序增长带来的输入 Token 成本。

---

## 📋 plan.yaml 格式详解

### 完整示例

```yaml
plan:
  task_summary: 实现用户头像上传功能
  confidence: high
  model_id: anthropic/claude-3.5-sonnet  # 针对此任务分配的主模型
  safety_buffer: 1.5                     # 全局安全倍率，覆盖 CLI 默认值
  total_files: 3
  estimated_total_diff_tokens: 25000
  estimated_total_regen_rounds: 2
  files:
    - path: backend/models/user.go
      action: modify
      subtask: false                   # 默认 false，继承前序上下文
      predicted_diff_tokens: 3000
      predicted_regen_times: 1
      description: 添加 avatar_url 字段
      
    - path: backend/api/upload.go
      action: create
      subtask: true                    # 标记为子任务：隔离上下文，1.2x 启动倍率
      predicted_diff_tokens: 12000
      predicted_regen_times: 2
      description: 实现 S3 上传逻辑（逻辑较独立，建议开启新上下文）
      
    - path: frontend/components/Upload.tsx
      action: create
      subtask: false                   # 继承 upload.go 的上下文进行联调
      predicted_diff_tokens: 10000
      predicted_regen_times: 3
      description: 实现上传 UI 组件
```

### 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `model_id` | string | ✅ | 针对此任务分配的模型 ID（用于抓取 TPS） |
| `safety_buffer` | float | ❌ | 全局安全倍率（默认 1.8）。YAML 中的定义优先级高于 CLI 参数 |
| `files[].subtask` | bool | ✅ | 是否为隔离子任务。`true` 则不继承上下文，且计算时增加 1.2x 倍率 |
| `files[].path` | string | ✅ | 文件相对路径。**必须按执行依赖顺序排列** |

---

## 🧮 calculate_plan 计算逻辑

### 1. TPS 获取优先级
1.  实时抓取 OpenRouter 对应 `model_id` 的 `Avg Throughput`。
2.  若抓取失败，使用默认值 **50 tok/s**。

### 2. 时长预估公式
- **普通任务**：`Time = (Diff * Regen) / TPS`
- **子任务**：`Time = (Diff * Regen * 1.2) / TPS`（补偿重新读取上下文的开销）

### 3. 成本预估（上下文累加）
- 除非 `subtask: true`，否则每个任务的 $Input\_Tokens$ 会包含之前所有任务的 $Diff + Output$。

---

## 🎯 系统提示词 (System Prompt)

```markdown
你是一个 Token-era 软件工程师。

核心规则：
1. **严禁**提及"小时/天/周"。所有估算基于 **diff token** 和 **再生轮次**。
2. 规划 `plan.yaml` 时，必须考虑代码依赖关系，按**执行顺序**排列 `files`。
3. 对于逻辑独立、容易导致上下文爆炸的模块，应设置 `subtask: true` 以隔离风险。
```
