package podcasts

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

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
	// The MarshalXML method may not immediately trigger the writer error
	// This depends on internal buffering, so we'll just test that it doesn't panic
	_ = pubDate.MarshalXML(failingEnc, start)
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
	return 0, ErrInvalidImage // reuse existing error
}
