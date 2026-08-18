package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/KillerBeast69/Mnemocast/ingest"
	"github.com/KillerBeast69/Mnemocast/store"
	"github.com/jackc/pgx/v5/pgxpool"
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
	// ---------------------------------------------------------
	// 1. CONNECT TO POSTGRES
	// ---------------------------------------------------------
	ctx := context.Background()
	connStr := "postgres://mnemocast_admin:mnemocast_password@localhost:5433/mnemocast_db?sslmode=disable"

	fmt.Println("Connecting to the database...")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize our generated sqlc queries struct using the connection pool
	queries := store.New(pool)
	fmt.Println("Connected successfully!")

	// ---------------------------------------------------------
	// 2. FETCH RSS FEED
	// ---------------------------------------------------------
	channelID := "UCsBjURrPoezykLs9EqgamOA" // Fireship
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID

	fmt.Printf("\nFetching the latest video for channel %s...\n", channelID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
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

	// ---------------------------------------------------------
	// 3. PROCESS THE FEED & SAVE TO DB
	// ---------------------------------------------------------
	if len(feed.Entries) > 0 {
		latestVideo := feed.Entries[0]

		// A. Save the Channel first (Foreign Key requirement)
		err = queries.CreateChannel(ctx, store.CreateChannelParams{
			ChannelID: channelID,
			Title:     feed.Title,
		})
		if err != nil {
			log.Printf("Failed to save channel to DB: %v", err)
		} else {
			fmt.Printf("✓ Channel '%s' saved/verified in DB.\n", feed.Title)
		}

		// B. Save the Video
		err = queries.CreateVideo(ctx, store.CreateVideoParams{
			VideoID:   latestVideo.VideoID,
			ChannelID: channelID,
			Title:     latestVideo.Title,
			Url:       latestVideo.Link.Href,
		})
		if err != nil {
			log.Printf("Failed to save video to DB: %v", err)
		} else {
			fmt.Printf("✓ Video '%s' saved/verified in DB.\n", latestVideo.Title)
		}

		// ---------------------------------------------------------
		// 4. FETCH THE TRANSCRIPT
		// ---------------------------------------------------------
		fmt.Println("\nAttempting to fetch transcript...")
		transcript, err := ingest.FetchTranscript(latestVideo.VideoID)
		if err != nil {
			log.Printf("Error fetching transcript: %v\n", err)
		} else {
			runes := []rune(transcript)
			previewLimit := 200
			if len(runes) < 200 {
				previewLimit = len(runes)
			}
			fmt.Printf("\n--- Transcript Preview ---\n%s...\n", string(runes[:previewLimit]))
			fmt.Printf("\nTotal transcript length: %d characters\n", len(transcript))
		}
	} else {
		fmt.Println("No videos found.")
	}
}
