package ttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Plan struct {
	Plan struct {
		TaskSummary              string      `yaml:"task_summary"`
		ModelID                  string      `yaml:"model_id"`
		EstimatedTotalDiffTokens int         `yaml:"estimated_total_diff_tokens"`
		SafetyRate               float64     `yaml:"safety_rate"`
		Files                    []FileEntry `yaml:"files"`
	} `yaml:"plan"`
}

type FileEntry struct {
	Path                string  `yaml:"path"`
	Action              string  `yaml:"action"`
	Subtask             bool    `yaml:"subtask"`
	PredictedDiffTokens int     `yaml:"predicted_diff_tokens"`
	RegenRate           int     `yaml:"regen_rate"`
	SafetyRate          float64 `yaml:"safety_rate"`
	Description         string  `yaml:"description"`
}

type Result struct {
	Report string
}

type ModelMetrics struct {
	TPS           float64
	InputPrice    float64
	OutputPrice   float64
	CanonicalSlug string
}

func RunPlan(planPath string, safetyRate float64) (Result, error) {
	data, err := os.ReadFile(planPath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to read plan: %w", err)
	}

	var p Plan
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Result{}, fmt.Errorf("failed to parse yaml: %w", err)
	}

	if p.Plan.SafetyRate > 0 {
		safetyRate = p.Plan.SafetyRate
	}

	metrics, err := GetModelMetrics(p.Plan.ModelID)
	if err != nil {
		metrics = ModelMetrics{TPS: 50.0}
	}

	return Calculate(&p, metrics, safetyRate), nil
}

func Calculate(p *Plan, metrics ModelMetrics, globalSafetyRate float64) Result {
	var totalInputTokens, totalOutputTokens, totalTimeSeconds, currentContext, totalWeightedRate float64

	for _, f := range p.Plan.Files {
		rate := float64(f.RegenRate)
		if rate < 1.0 {
			rate = 1.0
		}

		fileRate := f.SafetyRate
		if fileRate <= 0 {
			fileRate = globalSafetyRate
		}

		diff := float64(f.PredictedDiffTokens)
		multiplier := 1.0
		if f.Subtask {
			currentContext = 0
			multiplier = 1.2
		}

		taskInput := currentContext + diff
		taskOutput := diff * rate
		taskTime := (taskOutput / metrics.TPS) * multiplier * fileRate

		totalInputTokens += taskInput
		totalOutputTokens += taskOutput
		totalTimeSeconds += taskTime
		totalWeightedRate += fileRate
		currentContext += diff + taskOutput
	}

	totalTimeMinutes := totalTimeSeconds / 60.0
	totalCost := (totalInputTokens * metrics.InputPrice) + (totalOutputTokens * metrics.OutputPrice)
	avgWeightedRate := 0.0
	if len(p.Plan.Files) > 0 {
		avgWeightedRate = totalWeightedRate / float64(len(p.Plan.Files))
	}

	report := fmt.Sprintf(`===== Token-Time Planner - Result =====
Summary: %s
Model: %s (TPS: %.1f)
Stats: %d files

Estimated Tokens:
  - Input:  %.0f
  - Output: %.0f
  - Total:  %.0f

Estimated Cost: $%.4f
Estimated Time: %.1f minutes
(Weighted Safety Buffer: %.2fx)
==========================================`,
		p.Plan.TaskSummary, p.Plan.ModelID, metrics.TPS,
		len(p.Plan.Files),
		totalInputTokens, totalOutputTokens, totalInputTokens+totalOutputTokens,
		totalCost,
		totalTimeMinutes,
		avgWeightedRate,
	)

	return Result{Report: report}
}

func GetModelMetrics(modelID string) (ModelMetrics, error) {
	metrics := ModelMetrics{TPS: 50.0}
	
	input, output, slug, err := fetchPricing(modelID)
	if err == nil {
		metrics.InputPrice = input
		metrics.OutputPrice = output
		metrics.CanonicalSlug = slug
	}

	targetSlug := modelID
	if metrics.CanonicalSlug != "" {
		targetSlug = metrics.CanonicalSlug
	}

	if tps, err := fetchTPS(targetSlug); err == nil && tps > 0 {
		metrics.TPS = tps
	}

	return metrics, nil
}

func fetchTPS(modelID string) (float64, error) {
	apiURL := fmt.Sprintf("https://openrouter.ai/api/frontend/stats/throughput-comparison?permaslug=%s", strings.ReplaceAll(modelID, "/", "%%2F"))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", "https://openrouter.ai/")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var tr struct {
		Data []struct {
			Y map[string]float64 `json:"y"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return 0, err
	}
	if len(tr.Data) == 0 {
		return 0, fmt.Errorf("empty data")
	}

	latest := tr.Data[len(tr.Data)-1]
	var total float64
	var count int
	for _, val := range latest.Y {
		total += val
		count++
	}
	if count == 0 {
		return 0, fmt.Errorf("no providers")
	}
	return total / float64(count), nil
}

func fetchPricing(modelID string) (float64, float64, string, error) {
	resp, err := http.Get("https://openrouter.ai/api/v1/models")
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var mr struct {
		Data []struct {
			ID            string `json:"id"`
			CanonicalSlug string `json:"canonical_slug"`
			Pricing       struct {
				Prompt     string `json:"prompt"`
				Completion string `json:"completion"`
			} `json:"pricing"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return 0, 0, "", err
	}

	for _, m := range mr.Data {
		if m.ID == modelID || (len(modelID) > len(m.ID) && strings.HasPrefix(modelID, m.ID)) {
			var input, output float64
			fmt.Sscanf(m.Pricing.Prompt, "%f", &input)
			fmt.Sscanf(m.Pricing.Completion, "%f", &output)
			return input, output, m.CanonicalSlug, nil
		}
	}
	return 0, 0, "", fmt.Errorf("not found")
}
