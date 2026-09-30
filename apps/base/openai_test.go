package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const okBody = `{"choices":[{"finish_reason":"stop","message":{"tool_calls":[{"function":{"arguments":"{\"openings\":[],\"facts\":[]}"}}]}}]}`

func clientFor(url string) openAIClient {
	return openAIClient{apiKey: "k", model: "m", url: url, http: &http.Client{Timeout: time.Second}}
}

func TestRetriesRateLimitThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"slow down","code":"rate_limit_exceeded"}}`))
			return
		}
		_, _ = w.Write([]byte(okBody))
	}))
	defer srv.Close()

	// Retry-After 0 falls back to the 2s backoff, so keep the assertion on
	// the outcome, not the timing.
	if _, err := clientFor(srv.URL).callFunction(context.Background(), "f", "d", map[string]any{}, "p"); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Errorf("calls = %d, want 2", calls.Load())
	}
}

func TestDoesNotRetryABadKey(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"Incorrect API key","code":"invalid_api_key"}}`))
	}))
	defer srv.Close()

	_, err := clientFor(srv.URL).callFunction(context.Background(), "f", "d", map[string]any{}, "p")
	if err == nil || !strings.Contains(err.Error(), "check OPENAI_API_KEY") {
		t.Fatalf("err = %v, want a hint about OPENAI_API_KEY", err)
	}
	if calls.Load() != 1 {
		t.Errorf("calls = %d, want 1", calls.Load())
	}
}

func TestDoesNotRetrySpentQuota(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"You exceeded your quota","code":"insufficient_quota"}}`))
	}))
	defer srv.Close()

	_, err := clientFor(srv.URL).callFunction(context.Background(), "f", "d", map[string]any{}, "p")
	if err == nil || !strings.Contains(err.Error(), "out of credit") || calls.Load() != 1 {
		t.Fatalf("err = %v after %d calls, want one call and a credit hint", err, calls.Load())
	}
}

func TestReportsATruncatedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Replace(okBody, `"stop"`, `"length"`, 1)))
	}))
	defer srv.Close()

	_, err := clientFor(srv.URL).callFunction(context.Background(), "f", "d", map[string]any{}, "p")
	if err == nil || !strings.Contains(err.Error(), "cut off") {
		t.Fatalf("err = %v, want a cut-off error", err)
	}
}

func TestVideoIDShape(t *testing.T) {
	for in, want := range map[string]bool{
		extractVideoID("https://www.youtube.com/live/Vme7qk9NECM?si=abc&t=2812"): true,
		extractVideoID("https://youtu.be/Vme7qk9NECM"):                           true,
		extractVideoID("Vme7qk9NECM"):                                            true,
		extractVideoID("not a video"):                                            false,
	} {
		if videoIDShape.MatchString(in) != want {
			t.Errorf("%q: match = %v, want %v", in, !want, want)
		}
	}
}

func TestParseTimedTextSkipsLayoutParagraphs(t *testing.T) {
	body := []byte(`<timedtext format="3"><body>
<p t="2801510" w="1"></p>
<p t="2801520"><s>ligado?</s><s>  Mas  deve</s></p>
<p t="2811319">>> A  música</p>
</body></timedtext>`)
	got, err := parseTimedText(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].StartMs != 2801520 || got[0].Text != "ligado? Mas deve" || got[1].Text != ">> A música" {
		t.Errorf("segments = %+v", got)
	}
}
