package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KillerBeast69/Mnemocast/ingest"
)

// defining structs to tell Go exactly how to map the XML data into variables

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title   string `xml:"title"`
	VideoID string `xml:"videoId"` // youtube's custom tag for video ID
	Link    Link   `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"` // the "attr" option tells Go to grab the attribute inside the tag, not the text between it
}

func main() {
	// Gordan Ramsay's YouTube channel ID
	channelID := "UCIEv3lZ_tNXHzL3ox-_uUGQ"
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelID

	fmt.Printf("Fetching the latest video for channel %s...\n", channelID)

	// make the HTTP GET request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err) //  if there is no internat or the URL is invalid, crash and print the error
	}

	defer resp.Body.Close() // Ensure we close the connection when the function exits

	// read the raw XML data from the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err) // if we can't read the response body, crash and print the error
	}

	// unmarshal(parse) the raw XML data into our Go structs
	var feed Feed
	if err := xml.Unmarshal(body, &feed); err != nil {
		panic(err) // if we can't parse the XML, crash and print the error
	}

	// print the results
	fmt.Printf("Feed Title: %s\n", feed.Title)
	if len(feed.Entries) > 0 {
		latestVideo := feed.Entries[0]
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Channel: %s\n", feed.Title)
		fmt.Printf("Latest Video Title: %s\n", latestVideo.Title)
		fmt.Printf("Video ID: %s\n", latestVideo.VideoID)

		// fetch the transcript for the latest video
		fmt.Println("\nattempting to fetch transcript...")
		transcript, err := ingest.FetchTranscript(latestVideo.VideoID)
		if err != nil {
			fmt.Printf("Error fetching transcript: %v\n", err)
		} else {
			previewLimit := 200
			if len(transcript) < 200 {
				previewLimit = len(transcript)
			}
			fmt.Printf("Transcript Preview: %s\n", transcript[:previewLimit])
			fmt.Printf("total transcript length: %d characters\n", len(transcript))
		}
		fmt.Printf("Link: %s\n", latestVideo.Link.Href)
	} else {
		fmt.Println("No videos found :/.")
	}
}
