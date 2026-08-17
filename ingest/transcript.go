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
		return "", fmt.Errorf("failed to fetch video metadata: %w", err)
	}

	// try to fetch the englist transcript
	transcript, err := client.GetTranscript(video, "en")
	if err != nil {
		return "", fmt.Errorf("failed to fetch english transcript: %w")
	}

	// combine the text chunks into one large strings for the LLM
	var fullText strings.Builder
	for _, chunk := range transcript {
		fullText.WriteString(chunk.Text + " ")
	}

	return strings.TrimSpace(fullText.String()), nil
}
