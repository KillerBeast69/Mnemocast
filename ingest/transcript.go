package ingest

import (
	"fmt"
	"strings"

	"github.com/kkdai/youtube/v2"
)

// FetchTranscript takes youtube video ID and returns the full english transcript as a single string
func FetchTranscript(videoID string) (string, error) {
	client := youtube.Client{}

	// fetch the video metadata
	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch video metadata: %v", err)
	}

	if len(video.CaptionTracks) == 0 {
		return "", fmt.Errorf("no caption tracks available for video ID: %s", videoID)
	}

	var targetTrack *youtube.CaptionTrack
	fmt.Println(" [Debug] available caption tracks:")
	for i, track := range video.CaptionTracks {
		fmt.Printf(" -> %s, (Code: %s)\n", track.Name, track.LanguageCode)

		// grab the first track that starts with "en" (english)
		if targetTrack == nil && strings.HasPrefix(track.LanguageCode, "en") {
			targetTrack = &video.CaptionTracks[i]
		}
	}

	if targetTrack == nil {
		return "", fmt.Errorf("no english caption track found for video ID: %s", videoID)
	}

	fmt.Printf(" [Debug] selected track: %s\n", targetTrack.LanguageCode)

	// try to fetch the specific track we found
	transcript, err := client.GetTranscript(video, targetTrack.LanguageCode)
	if err != nil {
		return "", fmt.Errorf("failed to fetch english transcript: %v", err)
	}

	// combine the text chunks into one large strings for the LLM
	var fullText strings.Builder
	for _, chunk := range transcript {
		fullText.WriteString(chunk.Text + " ")
	}

	return strings.TrimSpace(fullText.String()), nil
}
