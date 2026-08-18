package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/KillerBeast69/Mnemocast/ingest"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title   string `xml:"title"`
	VideoID string `xml:"videoId"`
	Link    Link   `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func main() {
	channelID := "UCsBjURrPoezykLs9EqgamOA" // Fireship
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID

	fmt.Printf("Fetching the latest video for channel %s...\n", channelID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		// Replace panic with log.Fatalf or log.Printf
		log.Fatalf("Failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Title)
	if len(feed.Entries) > 0 {
		latestVideo := feed.Entries[0]
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Channel: %s\n", feed.Title)
		fmt.Printf("Latest Video Title: %s\n", latestVideo.Title)
		fmt.Printf("Video ID: %s\n", latestVideo.VideoID)
		fmt.Printf("Link: %s\n", latestVideo.Link.Href)

		fmt.Println("\nAttempting to fetch transcript...")
		transcript, err := ingest.FetchTranscript(latestVideo.VideoID)
		if err != nil {
			// Log the error instead of crashing, so the poller can continue to the next video
			log.Printf("Error fetching transcript: %v\n", err)
		} else {
			// FIX: Safe substring slicing using runes (Claude's catch!)
			runes := []rune(transcript)
			previewLimit := 200
			if len(runes) < 200 {
				previewLimit = len(runes)
			}
			fmt.Printf("Transcript Preview: %s...\n", string(runes[:previewLimit]))
			fmt.Printf("Total transcript length: %d characters\n", len(transcript))
		}
	} else {
		fmt.Println("No videos found.")
	}
}
