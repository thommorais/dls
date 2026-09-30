package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const openAIChatURL = "https://api.openai.com/v1/chat/completions"

// openAIClient is the one call the ingest needs: force a single function
// call and hand back its JSON arguments.
type openAIClient struct {
	apiKey string
	model  string
	http   *http.Client
}

// The model name is not hardcoded: it is read from OPENAI_MODEL so switching
// models never needs a code change.
func newOpenAIClientFromEnv() (openAIClient, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return openAIClient{}, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		return openAIClient{}, fmt.Errorf("OPENAI_MODEL is not set")
	}
	// A full transcript of a live show is long, so the response can take minutes.
	return openAIClient{apiKey: key, model: model, http: &http.Client{Timeout: 10 * time.Minute}}, nil
}

func (c openAIClient) callFunction(ctx context.Context, name, description string, schema map[string]any, prompt string) (json.RawMessage, error) {
	body, err := json.Marshal(map[string]any{
		"model":    c.model,
		"messages": []map[string]any{{"role": "user", "content": prompt}},
		"tools": []map[string]any{{
			"type":     "function",
			"function": map[string]any{"name": name, "description": description, "parameters": schema},
		}},
		"tool_choice": map[string]any{"type": "function", "function": map[string]any{"name": name}},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIChatURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai: %s: %s", resp.Status, bytes.TrimSpace(payload))
	}

	var parsed struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				ToolCalls []struct {
					Function struct {
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return nil, fmt.Errorf("openai: decoding response: %w", err)
	}
	if len(parsed.Choices) == 0 || len(parsed.Choices[0].Message.ToolCalls) == 0 {
		return nil, fmt.Errorf("openai: model did not return a function call")
	}
	if parsed.Choices[0].FinishReason == "length" {
		return nil, fmt.Errorf("openai: output was cut off before the function call finished")
	}
	return json.RawMessage(parsed.Choices[0].Message.ToolCalls[0].Function.Arguments), nil
}
