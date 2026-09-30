package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const openAIChatURL = "https://api.openai.com/v1/chat/completions"

// openAIClient is the one call the ingest needs: force a single function
// call and hand back its JSON arguments.
type openAIClient struct {
	apiKey string
	model  string
	url    string
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
	return openAIClient{apiKey: key, model: model, url: openAIChatURL, http: &http.Client{Timeout: 10 * time.Minute}}, nil
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

	payload, err := c.post(ctx, body)
	if err != nil {
		return nil, err
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

const maxAttempts = 4

// post retries what a retry can fix (rate limits, 5xx, dropped connections)
// and fails at once on what it cannot (bad key, bad model, oversized input).
func (c openAIClient) post(ctx context.Context, body []byte) ([]byte, error) {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		payload, retryAfter, err := c.postOnce(ctx, body)
		if err == nil {
			return payload, nil
		}
		lastErr = err

		var apiErr *openAIError
		retryable := !errors.As(err, &apiErr) || apiErr.retryable()
		if !retryable || attempt == maxAttempts || ctx.Err() != nil {
			break
		}

		wait := retryAfter
		if wait == 0 {
			wait = time.Duration(1<<attempt) * time.Second
		}
		log.Printf("openai: attempt %d of %d failed (%v), retrying in %s", attempt, maxAttempts, err, wait)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}
	return nil, lastErr
}

func (c openAIClient) postOnce(ctx context.Context, body []byte) ([]byte, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("openai: reading response: %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		return payload, 0, nil
	}

	seconds, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
	return nil, time.Duration(seconds) * time.Second, newOpenAIError(resp.StatusCode, payload)
}

type openAIError struct {
	status  int
	code    string
	message string
}

func newOpenAIError(status int, payload []byte) *openAIError {
	e := &openAIError{status: status, message: strings.TrimSpace(string(payload))}
	var body struct {
		Error struct {
			Message string `json:"message"`
			Code    string `json:"code"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	if json.Unmarshal(payload, &body) == nil && body.Error.Message != "" {
		e.message = body.Error.Message
		e.code = body.Error.Code
		if e.code == "" {
			e.code = body.Error.Type
		}
	}
	return e
}

// A 429 that says the quota is spent will not clear in a few seconds, so it
// is not retried.
func (e *openAIError) retryable() bool {
	if e.status == http.StatusTooManyRequests {
		return e.code != "insufficient_quota"
	}
	return e.status >= 500
}

func (e *openAIError) Error() string {
	msg := fmt.Sprintf("openai: %d %s", e.status, e.message)
	switch {
	case e.status == http.StatusUnauthorized:
		msg += " (check OPENAI_API_KEY)"
	case e.status == http.StatusNotFound || e.code == "model_not_found":
		msg += " (check OPENAI_MODEL: the name may be wrong, or the account has no access to it)"
	case e.code == "insufficient_quota":
		msg += " (the account is out of credit or over its spend limit)"
	case e.code == "context_length_exceeded":
		msg += " (the transcript does not fit this model's context window, use a model with a larger one)"
	}
	return msg
}
