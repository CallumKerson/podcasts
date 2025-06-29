package podcasts

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type testItem struct {
	title             string
	guid              string
	pubDate           time.Time
	pubDateStr        string
	enclosureURL      string
	enclosureLength   string
	enclosureType     string
	duration          time.Duration
	durationStr       string
	descriptionStr    string
	encodedContentStr string
}

var (
	validItems = []testItem{
		{
			title:           "Item 1",
			guid:            "http://www.example-podcast.com/my-podcast/1/episode",
			pubDate:         time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			pubDateStr:      "Thu, 01 Jan 2015 00:00:00 +0000",
			enclosureURL:    "http://www.example-podcast.com/my-podcast/1/episode-one",
			enclosureLength: "1234",
			enclosureType:   "MP3",
		},
		{
			title:           "Item 2",
			guid:            "http://www.example-podcast.com/my-podcast/2/episode",
			pubDate:         time.Date(2015, time.January, 2, 0, 0, 0, 0, time.UTC),
			pubDateStr:      "Fri, 02 Jan 2015 00:00:00 +0000",
			enclosureURL:    "http://www.example-podcast.com/my-podcast/2/episode-two",
			enclosureLength: "56445",
			enclosureType:   "WAV",
			duration:        time.Second * 94,
			durationStr:     "1:34",
		},
		{
			title:             "Item 3",
			guid:              "http://www.example-podcast.com/my-podcast/3/episode",
			pubDate:           time.Date(2015, time.January, 3, 0, 0, 0, 0, time.UTC),
			pubDateStr:        "Thu, 01 Jan 2015 00:00:00 +0000",
			enclosureURL:      "http://www.example-podcast.com/my-podcast/3/episode-three",
			enclosureLength:   "1234",
			enclosureType:     "MP3",
			descriptionStr:    "A short description of the podcast episode",
			encodedContentStr: "<h1>Item 3</h1><p>A <em>longer</em> description of the podcast, specifically designed for embedded HTML.</p>",
		},
	}
)

func TestContainsXmlHeader(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := `<?xml version="1.0" encoding="UTF-8"?>`
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsRssElement(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := `<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">`
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsChannelElement(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := `<channel>`
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	want = `</channel>`
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsTitleElement(t *testing.T) {
	podcast := &Podcast{Title: "my title"}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<title>%v</title>", podcast.Title)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsDescriptionElement(t *testing.T) {
	podcast := &Podcast{Description: "my description"}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<description>%v</description>", podcast.Description)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsLanguageElement(t *testing.T) {
	podcast := &Podcast{Language: "en"}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<language>%v</language>", podcast.Language)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsLinkElement(t *testing.T) {
	podcast := &Podcast{Link: "https://example.com"}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<link>%v</link>", podcast.Link)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsCopyrightElement(t *testing.T) {
	podcast := &Podcast{Copyright: "MIT"}
	data, err := getPodcastXML(podcast)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<copyright>%v</copyright>", podcast.Copyright)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsBlockElement(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast, Block)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:block>%v</itunes:block>", ValueYes)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsExplicitElement(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast, Explicit)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:explicit>%v</itunes:explicit>", ValueYes)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsCompleteElement(t *testing.T) {
	podcast := &Podcast{}
	data, err := getPodcastXML(podcast, Complete)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:complete>%v</itunes:complete>", ValueYes)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsAuthorElement(t *testing.T) {
	podcast := &Podcast{}
	author := testAuthor
	data, err := getPodcastXML(podcast, Author(author))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:author>%v</itunes:author>", author)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsNewFeedURLElement(t *testing.T) {
	podcast := &Podcast{}
	url := "http://localhost/my-test-url"
	data, err := getPodcastXML(podcast, NewFeedURL(url))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:new-feed-url>%v</itunes:new-feed-url>", url)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsSubtitleElement(t *testing.T) {
	podcast := &Podcast{}
	subtitle := testSubtitle
	data, err := getPodcastXML(podcast, Subtitle(subtitle))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:subtitle>%v</itunes:subtitle>", subtitle)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsSummaryElement(t *testing.T) {
	podcast := &Podcast{}
	summary := `Test Summary with <a href="http://www.example.com/">link</a>`
	data, err := getPodcastXML(podcast, Summary(summary))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf("<itunes:summary><![CDATA[%v]]></itunes:summary>", summary)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsOwnerElement(t *testing.T) {
	podcast := &Podcast{}
	name := "Test Name"
	email := "test@name.com"
	data, err := getPodcastXML(podcast, Owner(name, email))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := "<itunes:owner>"
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	want = fmt.Sprintf("<itunes:name>%v</itunes:name>", name)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	want = fmt.Sprintf("<itunes:email>%v</itunes:email>", email)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	want = "</itunes:owner>"
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsImageElement(t *testing.T) {
	podcast := &Podcast{}
	image := "http://localhost/myimage.jpg"
	data, err := getPodcastXML(podcast, Image(image))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := fmt.Sprintf(`<itunes:image href="%v"></itunes:image>`, image)
	if !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestContainsItemElements(t *testing.T) {
	podcast := setupPodcast()
	feed, err := podcast.Feed()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	data, err := feed.XML()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	for i := range validItems {
		item := &validItems[i]
		t.Run(item.title, func(t *testing.T) {
			validatePodcastItem(t, data, item)
		})
	}
}

func TestPodcastFeedWrite(t *testing.T) {
	podcast := &Podcast{}
	feed, err := podcast.Feed()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	var b bytes.Buffer
	err = feed.Write(&b)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func getPodcastXML(p *Podcast, options ...func(f *Feed) error) (string, error) {
	feed, err := p.Feed(options...)
	if err != nil {
		return "", err
	}
	return feed.XML()
}

func setupPodcast() *Podcast {
	podcast := &Podcast{}
	for i := range validItems {
		item := &validItems[i]

		var description *CDATAText
		if item.descriptionStr != "" {
			description = &CDATAText{Value: item.descriptionStr}
		}

		var encodedContent *CDATAText
		if item.encodedContentStr != "" {
			encodedContent = &CDATAText{Value: item.encodedContentStr}
		}

		podcast.AddItem(&Item{
			Title:          item.title,
			GUID:           item.guid,
			PubDate:        NewPubDate(item.pubDate),
			Duration:       NewDuration(item.duration),
			Description:    description,
			ContentEncoded: encodedContent,

			Enclosure: &Enclosure{
				URL:    item.enclosureURL,
				Length: item.enclosureLength,
				Type:   item.enclosureType,
			},
		})
	}
	return podcast
}

func validatePodcastItem(t *testing.T, data string, item *testItem) {
	if want := "<item>"; !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	if want := fmt.Sprintf("<title>%v</title>", item.title); !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	if want := fmt.Sprintf("<guid>%v</guid>", item.guid); !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	if want := fmt.Sprintf("<pubDate>%v</pubDate>", item.pubDateStr); !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	if want := fmt.Sprintf(`<enclosure url="%v" length="%v" type="%v"></enclosure>`, item.enclosureURL, item.enclosureLength, item.enclosureType); !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
	if item.durationStr != "" {
		if want := fmt.Sprintf("<itunes:duration>%v</itunes:duration>", item.durationStr); !strings.Contains(data, want) {
			t.Errorf("expected %v to contain %v", data, want)
		}
	}
	if item.descriptionStr != "" {
		if want := fmt.Sprintf("<description><![CDATA[%v]]></description>", item.descriptionStr); !strings.Contains(data, want) {
			t.Errorf("expected %v to contain %v", data, want)
		}
	}
	if item.encodedContentStr != "" {
		if want := fmt.Sprintf("<content:encoded><![CDATA[%v]]></content:encoded>", item.encodedContentStr); !strings.Contains(data, want) {
			t.Errorf("expected %v to contain %v", data, want)
		}
	}
	if want := "</item>"; !strings.Contains(data, want) {
		t.Errorf("expected %v to contain %v", data, want)
	}
}

func TestPodcastOptionsWithError(t *testing.T) {
	podcast := &Podcast{
		Title:       "Options Error Test",
		Description: "Testing options with errors",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	errorOption := func(f *Feed) error {
		return ErrInvalidURL
	}

	_, err := podcast.Feed(errorOption)
	if err == nil {
		t.Error("Expected error from failing option")
	}
}

func TestPodcastURLValidation(t *testing.T) {
	invalidURLs := []string{
		"",
		"not-a-url",
		"relative/path",
		"://invalid",
	}

	for _, url := range invalidURLs {
		t.Run("InvalidURL_"+url, func(t *testing.T) {
			podcast := &Podcast{
				Title:       "URL Test",
				Description: "Testing URL validation",
				Language:    "en",
				Link:        "https://example.com",
				Copyright:   "2024",
			}

			_, err := podcast.Feed(NewFeedURL(url))
			if err == nil {
				t.Errorf("Expected error for invalid URL: %s", url)
			}
		})
	}
}

func TestPodcastImageValidation(t *testing.T) {
	invalidURLs := []string{
		"",
		"not-a-url",
		"relative/path",
		"://invalid",
	}

	for _, url := range invalidURLs {
		t.Run("InvalidImageURL_"+url, func(t *testing.T) {
			podcast := &Podcast{
				Title:       "Image Test",
				Description: "Testing image validation",
				Language:    "en",
				Link:        "https://example.com",
				Copyright:   "2024",
			}

			_, err := podcast.Feed(Image(url))
			if err == nil {
				t.Errorf("Expected error for invalid image URL: %s", url)
			}
		})
	}
}

func TestPodcastMethods(t *testing.T) {
	podcast := &Podcast{
		Title:       "Methods Test",
		Description: "Testing podcast methods",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	// Test initial state
	if podcast.GetItemCount() != 0 {
		t.Error("Initial item count should be 0")
	}

	// Add items
	item1 := &Item{
		Title:   "Item 1",
		GUID:    "https://example.com/1",
		PubDate: NewPubDate(time.Now()),
	}
	item2 := &Item{
		Title:   "Item 2",
		GUID:    "https://example.com/2",
		PubDate: NewPubDate(time.Now()),
	}

	podcast.AddItem(item1)
	podcast.AddItem(item2)

	if podcast.GetItemCount() != 2 {
		t.Errorf("Expected 2 items, got %d", podcast.GetItemCount())
	}

	// Test GetItems (safe copy)
	items := podcast.GetItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 items from GetItems, got %d", len(items))
	}

	// Modify the returned slice - should not affect original
	items[0] = nil

	// Test GetItemsSlice (direct reference)
	directItems := podcast.GetItemsSlice()
	if len(directItems) != 2 {
		t.Errorf("Expected 2 items from GetItemsSlice, got %d", len(directItems))
	}

	// Original items should be unchanged
	if directItems[0] == nil {
		t.Error("Direct items slice should not be affected by copy modification")
	}
	if directItems[0].Title != "Item 1" {
		t.Error("First item title should be unchanged")
	}

	// Test AddItemWithCapacity
	item3 := &Item{
		Title:   "Item 3",
		GUID:    "https://example.com/3",
		PubDate: NewPubDate(time.Now()),
	}
	podcast.AddItemWithCapacity(item3, 100)

	if podcast.GetItemCount() != 3 {
		t.Error("Should have 3 items after AddItemWithCapacity")
	}

	// Test adding with capacity less than current length
	item4 := &Item{
		Title:   "Item 4",
		GUID:    "https://example.com/4",
		PubDate: NewPubDate(time.Now()),
	}
	podcast.AddItemWithCapacity(item4, 1) // Less than current length

	if podcast.GetItemCount() != 4 {
		t.Errorf("Expected 4 items, got %d", podcast.GetItemCount())
	}
}
