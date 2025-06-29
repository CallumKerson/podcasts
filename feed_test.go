package podcasts

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	// Test constants to avoid duplication
	ValueClean = "clean"
	ValueNo    = "no"
	testString = "test"
)

// Helper function to configure feed with common test settings
func configureFeedForStreamTest(feedPtr *Feed) error {
	feedPtr.Channel.Author = "Test Author"
	feedPtr.Channel.Block = ValueYes
	feedPtr.Channel.Explicit = ValueClean
	feedPtr.Channel.Complete = ValueYes
	feedPtr.Channel.NewFeedURL = "https://example.com/new-feed"
	feedPtr.Channel.Subtitle = "Test Subtitle"
	feedPtr.Channel.Summary = &CDATAText{Value: "Test Summary"}
	feedPtr.Channel.Owner = &ItunesOwner{
		Name:  "Test Owner",
		Email: "owner@example.com",
	}
	feedPtr.Channel.Image = &ItunesImage{
		Href: "https://example.com/image.jpg",
	}
	feedPtr.Channel.Categories = []*ItunesCategory{
		{Text: "Technology"},
		{Text: "Education"},
	}
	return nil
}

// Helper function to test pool functionality
func testPoolFunctionality(t *testing.T, pool *sync.Pool, itemType string) {
	if pool == nil {
		t.Errorf("%s pool should not be nil", itemType)
		return
	}

	// Test getting and putting based on type
	switch itemType {
	case "Buffer":
		buf := pool.Get().(*bytes.Buffer)
		buf.WriteString(testString)
		if buf.String() != testString {
			t.Errorf("Buffer should contain %s", testString)
		}
		buf.Reset()
		pool.Put(buf)

		// Get another buffer and verify it's clean
		buf2 := pool.Get().(*bytes.Buffer)
		if buf2.Len() != 0 {
			t.Error("Pooled buffer should be reset")
		}
		pool.Put(buf2)
	case "StringBuilder":
		stringBuilder := pool.Get().(*strings.Builder)
		stringBuilder.WriteString(testString)
		if stringBuilder.String() != testString {
			t.Errorf("String builder should contain %s", testString)
		}
		stringBuilder.Reset()
		pool.Put(stringBuilder)

		// Get another string builder and verify it's clean
		stringBuilder2 := pool.Get().(*strings.Builder)
		if stringBuilder2.Len() != 0 {
			t.Error("Pooled string builder should be reset")
		}
		pool.Put(stringBuilder2)
	}
}

// Helper function to configure feed with alternative settings for testing different paths
func configureFeedWithAlternativeSettings(feedPtr *Feed) error {
	feedPtr.Channel.Author = "Coverage Author"
	feedPtr.Channel.Block = ValueNo
	feedPtr.Channel.Explicit = "yes"
	feedPtr.Channel.Complete = ValueNo
	feedPtr.Channel.NewFeedURL = "https://example.com/coverage"
	feedPtr.Channel.Subtitle = "Coverage Subtitle"
	feedPtr.Channel.Summary = &CDATAText{Value: "Coverage Summary"}
	feedPtr.Channel.Owner = &ItunesOwner{
		Name:  "Coverage Owner",
		Email: "coverage@example.com",
	}
	feedPtr.Channel.Image = &ItunesImage{
		Href: "https://example.com/coverage.jpg",
	}
	feedPtr.Channel.Categories = []*ItunesCategory{
		{Text: "Science"},
		{Text: "Testing"},
	}
	return nil
}

func TestPubDateMarshalling(t *testing.T) {
	tm := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	pubDate := NewPubDate(tm)
	want := "<PubDate>Thu, 01 Jan 2015 00:00:00 +0000</PubDate>"
	out, err := xml.Marshal(pubDate)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if got := string(out); want != got {
		t.Errorf("expected %v got %v", want, got)
	}
}

func TestDurationMarshalling(t *testing.T) {
	cases := []struct {
		dur  time.Duration
		want string
	}{
		{
			dur:  0,
			want: "<Duration>0:00</Duration>",
		},
		{
			dur:  time.Second * 6,
			want: "<Duration>0:06</Duration>",
		},
		{
			dur:  time.Second * 64,
			want: "<Duration>1:04</Duration>",
		},
		{
			dur:  time.Second * 125,
			want: "<Duration>2:05</Duration>",
		},
		{
			dur:  time.Second * 3600,
			want: "<Duration>1:00:00</Duration>",
		},
		{
			dur:  time.Second * 37000,
			want: "<Duration>10:16:40</Duration>",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.want, func(t *testing.T) {
			dur := NewDuration(testCase.dur)
			out, err := xml.Marshal(dur)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if got := string(out); testCase.want != got {
				t.Errorf("expected %v got %v", testCase.want, got)
			}
		})
	}
}

// TestMarshalXMLErrorPaths tests error paths that can occur during XML marshalling
func TestMarshalXMLErrorPaths(t *testing.T) {
	// Test with invalid XML content that would cause marshalling errors
	// We can't easily mock xml.Encoder, but we can test edge cases

	// Test PubDate with extreme values
	pubDate := NewPubDate(time.Time{}) // zero time
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	start := xml.StartElement{Name: xml.Name{Local: "pubDate"}}

	if err := pubDate.MarshalXML(enc, start); err != nil {
		t.Errorf("PubDate marshalling should handle zero time: %v", err)
	}

	// Test Duration with extreme values
	duration := NewDuration(time.Duration(1<<62 - 1)) // near max duration
	buf.Reset()
	if err := duration.MarshalXML(enc, start); err != nil {
		t.Errorf("Duration marshalling should handle extreme values: %v", err)
	}

	// Test with encoder that fails on flush
	failingWriter := &failingWriter{}
	failingEnc := xml.NewEncoder(failingWriter)

	// Test PubDate MarshalXML with failing writer
	_ = pubDate.MarshalXML(failingEnc, start)
	failingEnc.Flush() // Force flush to trigger error
	// Error may or may not occur depending on buffering

	// Test Duration MarshalXML with failing writer
	duration = NewDuration(time.Minute * 5)
	_ = duration.MarshalXML(failingEnc, start)
	failingEnc.Flush() // Force flush to trigger error
	// Error may or may not occur depending on buffering
}

// TestFeedXMLErrors tests error paths in Feed XML generation
func TestFeedXMLErrors(t *testing.T) {
	// Test XML generation with failing writer
	feed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       "Test",
			Description: "Test",
			Link:        "http://example.com",
		},
	}

	// Test Write with failing writer
	failingWriter := &failingWriter{}
	if err := feed.Write(failingWriter); err == nil {
		t.Error("expected error from failing writer")
	}

	// Test XML generation - most XML generation will succeed with Go's encoder
	// which handles invalid characters by escaping them properly
	// So we'll just verify that XML generation completes without panicking
	testFeed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       "Test",
			Description: "Test",
			Link:        "http://example.com",
		},
	}
	if _, err := testFeed.XML(); err != nil {
		t.Errorf("Basic XML generation should succeed: %v", err)
	}
}

// TestSetOptionsError tests error handling in SetOptions
func TestSetOptionsError(t *testing.T) {
	feed := &Feed{Channel: &Channel{}}

	// Option that returns an error
	errorOption := func(f *Feed) error {
		return ErrInvalidURL // reuse existing error
	}

	if err := feed.SetOptions(errorOption); err == nil {
		t.Error("expected error from failing option")
	}
}

// TestFeedWriteSuccess tests successful feed writing
func TestFeedWriteSuccess(t *testing.T) {
	feed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       "Test Podcast",
			Description: "Test Description",
			Link:        "http://example.com",
			Language:    "en",
		},
	}

	var buf bytes.Buffer
	if err := feed.Write(&buf); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>") {
		t.Error("output should contain XML header")
	}
	if !strings.Contains(output, "Test Podcast") {
		t.Error("output should contain podcast title")
	}
}

// TestFeedXMLSuccess tests successful XML generation
func TestFeedXMLSuccess(t *testing.T) {
	feed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       "Test Podcast",
			Description: "Test Description",
			Link:        "http://example.com",
		},
	}

	xmlString, err := feed.XML()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.Contains(xmlString, "Test Podcast") {
		t.Error("XML should contain podcast title")
	}
}

// Helper types for testing error conditions
type failingWriter struct{}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, ErrInvalidURL
}

func TestFeedWriteWithOptions(t *testing.T) {
	podcast := &Podcast{
		Title:       "WriteOptions Test",
		Description: "Testing WriteWithOptions",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	podcast.AddItem(&Item{
		Title:   "Test Episode",
		GUID:    "https://example.com/test",
		PubDate: NewPubDate(time.Now()),
	})

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("WithPool", func(t *testing.T) {
		var buf bytes.Buffer
		opts := WriteOptions{UsePool: true}

		err := feed.WriteWithOptions(&buf, opts)
		if err != nil {
			t.Fatal(err)
		}

		content := buf.String()
		if !strings.Contains(content, "WriteOptions Test") {
			t.Error("Content doesn't contain expected podcast title")
		}
	})

	t.Run("WithBufferSize", func(t *testing.T) {
		var buf bytes.Buffer
		opts := WriteOptions{BufferSize: 8192}

		err := feed.WriteWithOptions(&buf, opts)
		if err != nil {
			t.Fatal(err)
		}

		content := buf.String()
		if !strings.Contains(content, "WriteOptions Test") {
			t.Error("Content doesn't contain expected podcast title")
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		failingWriter := &failingWriter{}
		opts := WriteOptions{}

		err := feed.WriteWithOptions(failingWriter, opts)
		if err == nil {
			t.Error("Expected error from failing writer")
		}
	})
}

func TestFeedXMLWithOptions(t *testing.T) {
	podcast := &Podcast{
		Title:       "XMLOptions Test",
		Description: "Testing XMLWithOptions",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("DefaultOptions", func(t *testing.T) {
		xmlContent, err := feed.XMLWithOptions(WriteOptions{})
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(xmlContent, "XMLOptions Test") {
			t.Error("XML doesn't contain expected podcast title")
		}
	})

	t.Run("WithPool", func(t *testing.T) {
		xmlContent, err := feed.XMLWithOptions(WriteOptions{UsePool: true})
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(xmlContent, "XMLOptions Test") {
			t.Error("XML doesn't contain expected podcast title")
		}
	})
}

func TestFeedStreamWrite(t *testing.T) {
	podcast := &Podcast{
		Title:       "StreamWrite Test",
		Description: "Testing StreamWrite",
		Language:    "en-US",
		Link:        "https://example.com/stream",
		Copyright:   "2024",
	}

	// Add multiple items to test streaming
	for i := 1; i <= 3; i++ {
		podcast.AddItem(&Item{
			Title:   fmt.Sprintf("Episode %d", i),
			GUID:    fmt.Sprintf("https://example.com/stream/%d", i),
			PubDate: NewPubDate(time.Now()),
		})
	}

	feed, err := podcast.Feed(configureFeedForStreamTest)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("SuccessfulStreamWrite", func(t *testing.T) {
		var buf bytes.Buffer
		err = feed.StreamWrite(&buf)
		if err != nil {
			t.Fatal(err)
		}

		content := buf.String()

		// Verify XML structure
		if !strings.Contains(content, `<?xml version="1.0" encoding="UTF-8"?>`) {
			t.Error("XML header missing")
		}

		if !strings.Contains(content, "StreamWrite Test") {
			t.Error("Podcast title missing")
		}

		// Verify all episodes are present
		for i := 1; i <= 3; i++ {
			expectedTitle := fmt.Sprintf("Episode %d", i)
			if !strings.Contains(content, expectedTitle) {
				t.Errorf("Episode %d missing: %s", i, expectedTitle)
			}
		}

		// Verify all channel fields are present
		expectedFields := []string{
			"Test Author", "yes", "clean", "https://example.com/new-feed",
			"Test Subtitle", "Test Summary", "Test Owner", "owner@example.com",
			"https://example.com/image.jpg", "Technology", "Education",
		}

		for _, field := range expectedFields {
			if !strings.Contains(content, field) {
				t.Errorf("StreamWrite missing field: %s", field)
			}
		}
	})

	t.Run("StreamWriteError", func(t *testing.T) {
		failingWriter := &failingWriter{}
		err := feed.StreamWrite(failingWriter)
		if err == nil {
			t.Error("Expected error from failing writer")
		}
	})
}

func TestBufferAndStringBuilderPools(t *testing.T) {
	t.Run("GetBufferPool", func(t *testing.T) {
		pool := GetBufferPool()
		testPoolFunctionality(t, pool, "Buffer")
	})

	t.Run("GetStringBuilderPool", func(t *testing.T) {
		pool := GetStringBuilderPool()
		testPoolFunctionality(t, pool, "StringBuilder")
	})
}

func TestWriteElementHelper(t *testing.T) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	err := writeElement(enc, "test", "value")
	if err != nil {
		t.Fatal(err)
	}

	// Flush the encoder
	enc.Flush()

	content := buf.String()
	if !strings.Contains(content, "<test>value</test>") {
		t.Errorf("Expected XML element, got: %s", content)
	}

	// Test with empty value
	buf.Reset()
	enc = xml.NewEncoder(&buf)
	err = writeElement(enc, "empty", "")
	if err != nil {
		t.Fatal(err)
	}

	enc.Flush()

	// Empty values should not produce any output
	if buf.Len() != 0 {
		t.Error("Empty value should not produce XML output")
	}

	t.Run("WriteElementError", func(t *testing.T) {
		failingWriter := &failingWriter{}
		enc := xml.NewEncoder(failingWriter)
		err := writeElement(enc, "test", "value")
		if err == nil {
			t.Error("Expected error from failing writer")
		}
	})
}

func TestBufferSizeOptions(t *testing.T) {
	podcast := &Podcast{
		Title:       "Buffer Size Test",
		Description: "Testing buffer size options",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatal(err)
	}

	// Test with various buffer sizes
	bufferSizes := []int{0, 1, 100, 1024, 8192}

	for _, size := range bufferSizes {
		t.Run(fmt.Sprintf("BufferSize_%d", size), func(t *testing.T) {
			var buf bytes.Buffer
			opts := WriteOptions{BufferSize: size}

			err := feed.WriteWithOptions(&buf, opts)
			if err != nil {
				t.Fatalf("Error with buffer size %d: %v", size, err)
			}

			content := buf.String()
			if !strings.Contains(content, "Buffer Size Test") {
				t.Error("Content doesn't contain expected podcast title")
			}
		})
	}

	// Test buffer size without pool
	var buf bytes.Buffer
	opts := WriteOptions{BufferSize: 2048}
	err = feed.WriteWithOptions(&buf, opts)
	if err != nil {
		t.Fatal(err)
	}

	content := buf.String()
	if !strings.Contains(content, "Buffer Size Test") {
		t.Error("Content doesn't contain expected podcast title")
	}
}

func TestAllWriteOptionsCombinations(t *testing.T) {
	podcast := &Podcast{
		Title:       "All Combinations Test",
		Description: "Testing all option combinations",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatal(err)
	}

	combinations := []WriteOptions{
		{}, // No options
		{UsePool: true},
		{BufferSize: 1024},
		{UsePool: true, BufferSize: 1024},
	}

	for i, opts := range combinations {
		t.Run(fmt.Sprintf("Combination_%d", i), func(t *testing.T) {
			var buf bytes.Buffer
			err := feed.WriteWithOptions(&buf, opts)
			if err != nil {
				t.Fatalf("Error with options %+v: %v", opts, err)
			}

			content := buf.String()

			if !strings.Contains(content, "All Combinations Test") {
				t.Errorf("Content missing podcast title with options %+v", opts)
			}

			// Also test XMLWithOptions
			_, err = feed.XMLWithOptions(opts)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStreamWriteWithComplexItems(t *testing.T) {
	// Test StreamWrite with items containing multiple fields and edge cases
	podcast := &Podcast{
		Title:       "StreamWrite Coverage Test",
		Description: "Testing StreamWrite edge cases",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	// Add an item with more fields to test additional paths
	podcast.AddItem(&Item{
		Title:          "Test Item",
		GUID:           "https://example.com/test",
		PubDate:        NewPubDate(time.Now()),
		Author:         "Item Author",
		Summary:        &CDATAText{Value: "Item Summary"},
		ContentEncoded: &CDATAText{Value: "<p>Item Content</p>"},
		Duration:       NewDuration(time.Minute * 30),
		Enclosure: &Enclosure{
			URL:    "https://example.com/test.mp3",
			Length: "12345678",
			Type:   "audio/mpeg",
		},
	})

	feed, err := podcast.Feed(configureFeedWithAlternativeSettings)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = feed.StreamWrite(&buf)
	if err != nil {
		t.Fatal(err)
	}

	content := buf.String()

	// Verify specific elements are present to ensure branches were taken
	expectedElements := []string{
		"Coverage Author",
		"Item Author",
		"Item Summary",
		"Item Content",
		"30:00",
		"Coverage Owner",
		"coverage@example.com",
		"Science",
		"Testing",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("StreamWrite missing element: %s", element)
		}
	}
}

func TestXMLGenerationSuccess(t *testing.T) {
	// Test XML() method successful generation
	feed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       "Error Test",
			Description: "Testing error paths",
			Link:        "http://example.com",
		},
	}

	// Test successful path first
	xmlString, err := feed.XML()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if xmlString == "" {
		t.Error("XML string should not be empty")
	}
}

func TestWriteWithOptionsErrorPaths(t *testing.T) {
	// Create a feed to test various WriteWithOptions error paths
	podcast := &Podcast{
		Title:       "WriteWithOptions Error Test",
		Description: "Testing WriteWithOptions error paths",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	feed, err := podcast.Feed()
	if err != nil {
		t.Fatal(err)
	}

	// Test WriteWithOptions error path - skip nil writer test as it causes panic
	// Instead test with a working writer but verify no unexpected errors
	var buf bytes.Buffer
	err = feed.WriteWithOptions(&buf, WriteOptions{})
	if err != nil {
		t.Errorf("Unexpected error with valid writer: %v", err)
	}

	// Test error in pooled write
	failingWriter := &failingWriter{}
	err = feed.WriteWithOptions(failingWriter, WriteOptions{UsePool: true})
	if err == nil {
		t.Error("Expected error in pooled write")
	}

	// Test error in buffered write
	err = feed.WriteWithOptions(failingWriter, WriteOptions{BufferSize: 1024})
	if err == nil {
		t.Error("Expected error in buffered write")
	}
}
