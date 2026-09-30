package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	kkyoutube "github.com/kkdai/youtube/v2"
)

type transcriptSegment struct {
	StartMs int
	Text    string
}

// fetchTranscript reads the video's own caption track. The library's
// GetTranscript calls YouTube's get_transcript endpoint, which answers 400 for
// every video, so the caption track url from the video info is used instead.
func fetchTranscript(ctx context.Context, video *kkyoutube.Video, lang string) ([]transcriptSegment, error) {
	track, err := pickCaptionTrack(video.CaptionTracks, lang)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, track.BaseURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching captions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching captions: unexpected status %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading captions: %w", err)
	}
	return parseTimedText(body)
}

func pickCaptionTrack(tracks []kkyoutube.CaptionTrack, lang string) (kkyoutube.CaptionTrack, error) {
	available := make([]string, 0, len(tracks))
	for _, t := range tracks {
		if t.LanguageCode == lang || strings.HasPrefix(t.LanguageCode, lang+"-") {
			return t, nil
		}
		available = append(available, t.LanguageCode)
	}
	if len(tracks) == 0 {
		return kkyoutube.CaptionTrack{}, fmt.Errorf("the video has no captions")
	}
	return kkyoutube.CaptionTrack{}, fmt.Errorf("no captions in %q, available: %s (use --lang)", lang, strings.Join(available, ", "))
}

// parseTimedText reads YouTube's timedtext format 3: <p t="ms"> paragraphs
// whose words sit in <s> elements. Paragraphs with no words are layout only.
func parseTimedText(body []byte) ([]transcriptSegment, error) {
	var doc struct {
		Paragraphs []struct {
			StartMs int      `xml:"t,attr"`
			Text    string   `xml:",chardata"`
			Words   []string `xml:"s"`
		} `xml:"body>p"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("parsing captions: %w", err)
	}

	segments := make([]transcriptSegment, 0, len(doc.Paragraphs))
	for _, p := range doc.Paragraphs {
		text := strings.Join(strings.Fields(p.Text+" "+strings.Join(p.Words, " ")), " ")
		if text == "" {
			continue
		}
		segments = append(segments, transcriptSegment{StartMs: p.StartMs, Text: text})
	}
	return segments, nil
}
