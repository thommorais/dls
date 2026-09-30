package main

import (
	"context"
	"fmt"
	"log"

	youtubedata "google.golang.org/api/youtube/v3"
)

// channelVideoIDs returns every video in every public playlist of a channel,
// each once, in the order first seen. A video that sits in several playlists
// would otherwise be sent to the model several times.
func channelVideoIDs(ctx context.Context, svc *youtubedata.Service, handle string) ([]string, error) {
	ch, err := svc.Channels.List([]string{"id"}).ForHandle(handle).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("channel %s: %w", handle, explainYouTubeAPIError(err))
	}
	if len(ch.Items) == 0 {
		return nil, fmt.Errorf("no channel found for handle %q", handle)
	}
	channelID := ch.Items[0].Id

	type playlist struct{ id, title string }
	var playlists []playlist
	pageToken := ""
	for {
		resp, err := svc.Playlists.List([]string{"snippet"}).ChannelId(channelID).MaxResults(50).PageToken(pageToken).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("listing playlists: %w", explainYouTubeAPIError(err))
		}
		for _, p := range resp.Items {
			playlists = append(playlists, playlist{p.Id, p.Snippet.Title})
		}
		if pageToken = resp.NextPageToken; pageToken == "" {
			break
		}
	}
	if len(playlists) == 0 {
		return nil, fmt.Errorf("channel %s has no public playlists", handle)
	}

	seen := map[string]bool{}
	var ids []string
	for _, p := range playlists {
		added := 0
		pageToken = ""
		for {
			resp, err := svc.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(p.id).MaxResults(50).PageToken(pageToken).Context(ctx).Do()
			if err != nil {
				return nil, fmt.Errorf("playlist %q: %w", p.title, explainYouTubeAPIError(err))
			}
			for _, item := range resp.Items {
				id := item.ContentDetails.VideoId
				if !seen[id] {
					seen[id] = true
					ids = append(ids, id)
					added++
				}
			}
			if pageToken = resp.NextPageToken; pageToken == "" {
				break
			}
		}
		log.Printf("playlist %q: %d new videos", p.title, added)
	}
	return ids, nil
}
