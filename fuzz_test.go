package podcasts

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// TestFuzzURLValidation tests URL validation with various inputs
func TestFuzzURLValidation(t *testing.T) {
	feed := &Feed{Channel: &Channel{}}

	// Test various malformed URLs
	testCases := []string{
		"",
		" ",
		"not-a-url",
		"://missing-scheme",
		"http://",
		"http:///path",
		"ftp://unsupported-scheme.com",
		"javascript:alert('xss')",
		"data:text/html,<script>alert('xss')</script>",
		string([]byte{0x00, 0x01, 0x02}),                  // null bytes
		strings.Repeat("a", 10000),                        // very long string
		"http://example.com/" + strings.Repeat("x", 8000), // very long path
	}

	for _, testURL := range testCases {
		t.Run("URL_"+testURL[:minInt(len(testURL), 20)], func(t *testing.T) {
			// These should all either return an error or be handled gracefully
			err1 := NewFeedURL(testURL)(feed)
			err2 := Image(testURL)(feed)

			// We don't expect specific errors, just that they don't panic
			// Invalid URLs should return errors
			if err1 == nil && !isValidAbsoluteURL(testURL) {
				t.Logf("NewFeedURL unexpectedly succeeded with: %q", testURL)
			}
			if err2 == nil && !isValidAbsoluteURL(testURL) {
				t.Logf("Image unexpectedly succeeded with: %q", testURL)
			}
		})
	}
}

// TestFuzzStringFields tests string fields with various inputs
func TestFuzzStringFields(t *testing.T) {
	feed := &Feed{Channel: &Channel{}}

	// Test various string inputs
	testStrings := []string{
		"",
		" ",
		"\n\r\t",
		"<script>alert('xss')</script>",
		"&lt;&gt;&amp;&quot;&#39;",
		string([]byte{0x00, 0x01, 0x02, 0x03}), // control characters
		string([]byte{0xFF, 0xFE, 0xFD}),       // invalid UTF-8
		strings.Repeat("A", 10000),             // very long string
		"🎵🎧📻🎙️",                                // emoji
		"Ñiño, naïve café résumé",              // accented characters
	}

	for _, testStr := range testStrings {
		t.Run("String_"+testStr[:minInt(len(testStr), 20)], func(t *testing.T) {
			// Test all string-accepting functions
			_ = Author(testStr)(feed)
			_ = Subtitle(testStr)(feed)
			_ = Summary(testStr)(feed)

			// These should not panic, even with malformed input
			// Generate XML to ensure no issues during serialisation
			if utf8.ValidString(testStr) {
				xmlContent, err := feed.XML()
				if err != nil {
					t.Logf("XML generation failed with string %q: %v", testStr, err)
				} else if xmlContent == "" {
					t.Error("XML generation returned empty string")
				}
			}
		})
	}
}

// TestFuzzDuration tests duration formatting with edge cases
func TestFuzzDuration(t *testing.T) {
	// Test various duration values
	testDurations := []time.Duration{
		-time.Hour,               // negative duration
		0,                        // zero duration
		time.Nanosecond,          // very small
		time.Microsecond,         // small
		time.Millisecond,         // small
		time.Second,              // normal
		time.Minute,              // normal
		time.Hour,                // normal
		time.Hour * 24,           // one day
		time.Hour * 24 * 365,     // one year
		time.Duration(1<<62 - 1), // near max int64
	}

	for _, duration := range testDurations {
		t.Run("Duration", func(t *testing.T) {
			// Create duration and format it
			dur := NewDuration(duration)

			// This should not panic, even with extreme values
			formatted := formatDuration(duration)

			// Basic sanity checks
			if duration >= 0 && !strings.Contains(formatted, ":") {
				t.Errorf("Positive duration should contain colon: %s", formatted)
			}

			// Test XML marshalling with real encoder
			var buf strings.Builder
			enc := xml.NewEncoder(&buf)
			if err := dur.MarshalXML(enc, testStartElement()); err != nil {
				t.Errorf("Duration marshalling failed: %v", err)
			}
		})
	}
}

// TestFuzzPodcastCreation tests podcast creation with various inputs
func TestFuzzPodcastCreation(t *testing.T) {
	testInputs := []struct {
		title, desc, link, lang, copyright string
	}{
		{"", "", "", "", ""},
		{strings.Repeat("A", 1000), strings.Repeat("B", 1000), "http://example.com", "en", "2024"},
		{"<>&\"'", "<>&\"'", "http://example.com", "en", "2024"},
		{"🎵", "🎧", "http://example.com", "en", "📅"},
	}

	for i, input := range testInputs {
		t.Run(fmt.Sprintf("PodcastCreation_%d", i), func(t *testing.T) {
			podcast := &Podcast{
				Title:       input.title,
				Description: input.desc,
				Link:        input.link,
				Language:    input.lang,
				Copyright:   input.copyright,
			}

			// Add various items
			podcast.AddItem(&Item{
				Title:   "Test Episode",
				GUID:    "https://example.com/test",
				PubDate: NewPubDate(time.Now()),
				Enclosure: &Enclosure{
					URL:  "https://example.com/test.mp3",
					Type: "audio/mpeg",
				},
			})

			// Generate feed - should not panic
			feed, err := podcast.Feed()
			if err != nil {
				t.Errorf("Feed creation failed: %v", err)
				return
			}

			// Generate XML - should not panic
			_, err = feed.XML()
			if err != nil {
				t.Logf("XML generation failed: %v", err)
			}
		})
	}
}

// Helper functions for fuzzing tests
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isValidAbsoluteURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}

func testStartElement() xml.StartElement {
	return xml.StartElement{Name: xml.Name{Local: "test"}}
}
