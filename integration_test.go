package podcasts

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestCompleteIntegration tests the full workflow of creating a podcast feed
func TestCompleteIntegration(t *testing.T) {
	// Create a podcast with multiple episodes
	podcast := &Podcast{
		Title:       "Integration Test Podcast",
		Description: "A test podcast for integration testing",
		Language:    "en-US",
		Link:        "https://example.com/podcast",
		Copyright:   "2024 Test Corporation",
	}

	// Add episodes with various content types
	podcast.AddItem(&Item{
		Title:       "Episode 1: Introduction",
		GUID:        "https://example.com/episode-1",
		PubDate:     NewPubDate(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
		Duration:    NewDuration(time.Minute * 30),
		Description: &CDATAText{Value: "An introduction to our podcast"},
		Enclosure: &Enclosure{
			URL:    "https://example.com/episode-1.mp3",
			Length: "14567890",
			Type:   "audio/mpeg",
		},
	})

	podcast.AddItem(&Item{
		Title:          "Episode 2: Advanced Topics",
		GUID:           "https://example.com/episode-2",
		PubDate:        NewPubDate(time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC)),
		Duration:       NewDuration(time.Hour + time.Minute*15),
		Description:    &CDATAText{Value: "Diving deep into advanced topics"},
		ContentEncoded: &CDATAText{Value: "<p>This episode covers <strong>advanced topics</strong> in detail.</p>"},
		Author:         "John Doe",
		Explicit:       "no",
		Subtitle:       "Advanced content for experienced listeners",
		Summary:        &CDATAText{Value: "A comprehensive look at advanced concepts"},
		Enclosure: &Enclosure{
			URL:    "https://example.com/episode-2.mp3",
			Length: "25678901",
			Type:   "audio/mpeg",
		},
	})

	// Generate feed with all options
	feed, err := podcast.Feed(
		Author("Test Author"),
		Block,
		Explicit,
		Complete,
		NewFeedURL("https://example.com/new-feed"),
		Subtitle("Test Podcast Subtitle"),
		Summary("This is a comprehensive test podcast covering various topics"),
		Owner("Test Owner", "test@example.com"),
		Image("https://example.com/podcast-artwork.jpg"),
	)

	if err != nil {
		t.Fatalf("Failed to create feed: %v", err)
	}

	// Test XML generation
	xmlContent, err := feed.XML()
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Validate XML structure
	if !strings.Contains(xmlContent, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("XML should contain proper XML declaration")
	}

	// Validate RSS structure
	expectedElements := []string{
		"<rss",
		"xmlns:itunes=\"http://www.itunes.com/dtds/podcast-1.0.dtd\"",
		"xmlns:content=\"http://purl.org/rss/1.0/modules/content/\"",
		"version=\"2.0\"",
		"<channel>",
		"<title>Integration Test Podcast</title>",
		"<description>A test podcast for integration testing</description>",
		"<language>en-US</language>",
		"<link>https://example.com/podcast</link>",
		"<copyright>2024 Test Corporation</copyright>",
		"<itunes:author>Test Author</itunes:author>",
		"<itunes:block>yes</itunes:block>",
		"<itunes:explicit>yes</itunes:explicit>",
		"<itunes:complete>yes</itunes:complete>",
		"<itunes:new-feed-url>https://example.com/new-feed</itunes:new-feed-url>",
		"<itunes:subtitle>Test Podcast Subtitle</itunes:subtitle>",
		"<itunes:owner>",
		"<itunes:name>Test Owner</itunes:name>",
		"<itunes:email>test@example.com</itunes:email>",
		"<itunes:image href=\"https://example.com/podcast-artwork.jpg\">",
		"<item>",
		"<title>Episode 1: Introduction</title>",
		"<guid>https://example.com/episode-1</guid>",
		"<pubDate>Mon, 01 Jan 2024 12:00:00 +0000</pubDate>",
		"<itunes:duration>30:00</itunes:duration>",
		"<enclosure url=\"https://example.com/episode-1.mp3\" length=\"14567890\" type=\"audio/mpeg\">",
		"<title>Episode 2: Advanced Topics</title>",
		"<itunes:duration>1:15:00</itunes:duration>",
		"<content:encoded>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(xmlContent, element) {
			t.Errorf("XML should contain element: %s", element)
		}
	}

	// Test that XML is valid by parsing it
	var parsedFeed Feed
	if err := xml.Unmarshal([]byte(xmlContent), &parsedFeed); err != nil {
		t.Fatalf("Generated XML is not valid: %v", err)
	}

	// Verify parsed content matches original
	if parsedFeed.Channel.Title != podcast.Title {
		t.Errorf("Parsed title doesn't match: expected %s, got %s", podcast.Title, parsedFeed.Channel.Title)
	}

	// Note: Items don't have explicit XML tags so they won't be parsed back
	// But we can verify the XML contains the item content
	itemCount := strings.Count(xmlContent, "<item>")
	if itemCount != 2 {
		t.Errorf("Expected 2 items in XML, got %d", itemCount)
	}
}

// TestXMLValidation tests that generated XML is well-formed
func TestXMLValidation(t *testing.T) {
	podcast := &Podcast{
		Title:       "Validation Test",
		Description: "Testing XML validation with special characters: <>&\"'",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	podcast.AddItem(&Item{
		Title:       "Episode with Special Characters: <>&\"'",
		GUID:        "https://example.com/special",
		PubDate:     NewPubDate(time.Now()),
		Description: &CDATAText{Value: "Description with <b>HTML</b> & special chars"},
		Enclosure: &Enclosure{
			URL:    "https://example.com/episode.mp3",
			Length: "12345",
			Type:   "audio/mpeg",
		},
	})

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatalf("Failed to create feed: %v", err)
	}

	xmlContent, err := feed.XML()
	if err != nil {
		t.Fatalf("Failed to generate XML: %v", err)
	}

	// Validate that XML can be parsed
	var parsed Feed
	if err := xml.Unmarshal([]byte(xmlContent), &parsed); err != nil {
		t.Fatalf("Generated XML is not valid: %v\nXML content:\n%s", err, xmlContent)
	}

	// Verify special characters are properly escaped in non-CDATA sections
	if !strings.Contains(xmlContent, "&lt;") || !strings.Contains(xmlContent, "&gt;") {
		t.Error("Special characters should be properly escaped in XML")
	}

	// Verify CDATA sections contain unescaped content
	if !strings.Contains(xmlContent, "<![CDATA[") {
		t.Error("CDATA sections should be present for description content")
	}
}

// TestLargeFeedGeneration tests handling of feeds with many episodes
func TestLargeFeedGeneration(t *testing.T) {
	podcast := &Podcast{
		Title:       "Large Feed Test",
		Description: "Testing with many episodes",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	// Add 100 episodes
	for episodeNum := 1; episodeNum <= 100; episodeNum++ {
		podcast.AddItem(&Item{
			Title:    fmt.Sprintf("Episode %d", episodeNum),
			GUID:     fmt.Sprintf("https://example.com/episode-%d", episodeNum),
			PubDate:  NewPubDate(time.Date(2024, 1, episodeNum%28+1, 12, 0, 0, 0, time.UTC)),
			Duration: NewDuration(time.Minute * time.Duration(20+episodeNum%40)),
			Enclosure: &Enclosure{
				URL:    fmt.Sprintf("https://example.com/episode-%d.mp3", episodeNum),
				Length: fmt.Sprintf("%d", 1000000+episodeNum*10000),
				Type:   "audio/mpeg",
			},
		})
	}

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatalf("Failed to create large feed: %v", err)
	}

	xmlContent, err := feed.XML()
	if err != nil {
		t.Fatalf("Failed to generate XML for large feed: %v", err)
	}

	// Count episodes in generated XML
	episodeCount := strings.Count(xmlContent, "<item>")
	if episodeCount != 100 {
		t.Errorf("Expected 100 episodes in XML, found %d", episodeCount)
	}

	// Ensure XML is still valid
	var parsed Feed
	if err := xml.Unmarshal([]byte(xmlContent), &parsed); err != nil {
		t.Fatalf("Large feed XML is not valid: %v", err)
	}

	// Verify the basic feed structure is intact
	if parsed.Channel.Title != podcast.Title {
		t.Errorf("Large feed title doesn't match")
	}
}

// TestEdgeCases tests various edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyPodcast", func(t *testing.T) {
		podcast := &Podcast{}
		feed, err := podcast.Feed()
		if err != nil {
			t.Errorf("Empty podcast should not cause error: %v", err)
		}

		xmlContent, err := feed.XML()
		if err != nil {
			t.Errorf("Empty podcast XML generation should not fail: %v", err)
		}

		if !strings.Contains(xmlContent, "<channel>") {
			t.Error("Even empty podcast should have channel element")
		}
	})

	t.Run("ZeroDuration", func(t *testing.T) {
		podcast := &Podcast{
			Title:       "Zero Duration Test",
			Description: "Test",
			Link:        "https://example.com",
		}

		podcast.AddItem(&Item{
			Title:    "Zero Duration Episode",
			GUID:     "https://example.com/zero",
			PubDate:  NewPubDate(time.Now()),
			Duration: NewDuration(0),
			Enclosure: &Enclosure{
				URL:  "https://example.com/zero.mp3",
				Type: "audio/mpeg",
			},
		})

		feed, err := podcast.Feed()
		if err != nil {
			t.Errorf("Zero duration should not cause error: %v", err)
		}

		xmlContent, err := feed.XML()
		if err != nil {
			t.Errorf("Zero duration XML generation should not fail: %v", err)
		}

		if !strings.Contains(xmlContent, "<itunes:duration>0:00</itunes:duration>") {
			t.Error("Zero duration should format as 0:00")
		}
	})

	t.Run("VeryLongDuration", func(t *testing.T) {
		podcast := &Podcast{
			Title:       "Long Duration Test",
			Description: "Test",
			Link:        "https://example.com",
		}

		// 25 hours duration
		longDuration := time.Hour*25 + time.Minute*30 + time.Second*45
		podcast.AddItem(&Item{
			Title:    "Very Long Episode",
			GUID:     "https://example.com/long",
			PubDate:  NewPubDate(time.Now()),
			Duration: NewDuration(longDuration),
			Enclosure: &Enclosure{
				URL:  "https://example.com/long.mp3",
				Type: "audio/mpeg",
			},
		})

		feed, err := podcast.Feed()
		if err != nil {
			t.Errorf("Long duration should not cause error: %v", err)
		}

		xmlContent, err := feed.XML()
		if err != nil {
			t.Errorf("Long duration XML generation should not fail: %v", err)
		}

		if !strings.Contains(xmlContent, "<itunes:duration>25:30:45</itunes:duration>") {
			t.Error("Long duration should format correctly as HH:MM:SS")
		}
	})
}
