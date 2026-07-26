package podcasts

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	itunesXMLNS  = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	contentXMLNS = "http://purl.org/rss/1.0/modules/content/"
	rssVersion   = "2.0"
	rfc2822      = "Mon, 02 Jan 2006 15:04:05 -0700"
)

// NewPubDate returns a new PubDate.
func NewPubDate(d time.Time) *PubDate {
	return &PubDate{d}
}

// PubDate represents pubDate attribute of given podcast item.
type PubDate struct {
	time.Time
}

// MarshalXML marshalls pubdate using the rfc2822 time format.
func (p PubDate) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if err := encoder.EncodeToken(start); err != nil {
		return err
	}
	if err := encoder.EncodeToken(xml.CharData(p.Format(rfc2822))); err != nil {
		return err
	}
	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

// NewDuration returns a new Duration.
func NewDuration(d time.Duration) *Duration {
	return &Duration{d}
}

// Duration represents itunes:duration attribute of given podcast item.
type Duration struct {
	time.Duration
}

// MarshalXML marshalls duration using HH:MM:SS, H:MM:SS, MM:SS, M:SS formats.
func (d Duration) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if err := encoder.EncodeToken(start); err != nil {
		return err
	}
	if err := encoder.EncodeToken(xml.CharData(formatDuration(d.Duration))); err != nil {
		return err
	}
	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

// formatDuration formats duration in these formats: HH:MM:SS, H:MM:SS, MM:SS, M:SS.
func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	hours := total / 3600
	total %= 3600
	minutes := total / 60
	total %= 60

	var builder strings.Builder
	if hours > 0 {
		builder.WriteString(strconv.Itoa(hours) + ":")
	}
	if hours > 0 && minutes < 10 {
		builder.WriteString("0")
	}
	builder.WriteString(strconv.Itoa(minutes) + ":")
	if total < 10 {
		builder.WriteString("0")
	}
	builder.WriteString(strconv.Itoa(total))
	return builder.String()
}

// ItunesOwner represents the itunes:owner of given channel.
type ItunesOwner struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

// ItunesImage represents the itunes:image of given item or channel.
type ItunesImage struct {
	XMLName xml.Name `xml:"itunes:image"`
	Href    string   `xml:"href,attr"`
}

// ItunesCategory represents itunes:category of given channel.
type ItunesCategory struct {
	XMLName    xml.Name `xml:"itunes:category"`
	Text       string   `xml:"text,attr"`
	Categories []*ItunesCategory
}

// Enclosure represents audio or video file of given item.
type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr,omitempty"`
	Type    string   `xml:"type,attr"`
}

// CDATAText represents some content that may contain
// embedded HTML such as <a href="...">...</a> links.
type CDATAText struct {
	Value string `xml:",cdata"`
}

// Item represents item of given channel.
type Item struct {
	XMLName         xml.Name   `xml:"item"`
	Title           string     `xml:"title"`
	GUID            string     `xml:"guid"`
	PubDate         *PubDate   `xml:"pubDate"`
	Description     *CDATAText `xml:"description,omitempty"`
	ContentEncoded  *CDATAText `xml:"content:encoded,omitempty"`
	Author          string     `xml:"itunes:author,omitempty"`
	Block           string     `xml:"itunes:block,omitempty"`
	Duration        *Duration  `xml:"itunes:duration,omitempty"`
	Explicit        string     `xml:"itunes:explicit,omitempty"`
	ClosedCaptioned string     `xml:"itunes:isClosedCaptioned,omitempty"`
	Order           int        `xml:"itunes:order,omitempty"`
	Subtitle        string     `xml:"itunes:subtitle,omitempty"`
	Summary         *CDATAText `xml:"itunes:summary,omitempty"`
	Enclosure       *Enclosure
	Image           *ItunesImage
}

// Channel represents a RSS channel for given podcast.
type Channel struct {
	XMLName     xml.Name   `xml:"channel"`
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Copyright   string     `xml:"copyright"`
	Language    string     `xml:"language"`
	Description string     `xml:"description"`
	Author      string     `xml:"itunes:author,omitempty"`
	Block       string     `xml:"itunes:block,omitempty"`
	Explicit    string     `xml:"itunes:explicit,omitempty"`
	Complete    string     `xml:"itunes:complete,omitempty"`
	NewFeedURL  string     `xml:"itunes:new-feed-url,omitempty"`
	Subtitle    string     `xml:"itunes:subtitle,omitempty"`
	Summary     *CDATAText `xml:"itunes:summary,omitempty"`
	Owner       *ItunesOwner
	Image       *ItunesImage
	Items       []*Item
	Categories  []*ItunesCategory
}

// Feed wraps the given RSS channel.
type Feed struct {
	XMLName      xml.Name `xml:"rss"`
	ItunesXMLNS  string   `xml:"xmlns:itunes,attr"`
	ContentXMLNS string   `xml:"xmlns:content,attr"`
	Version      string   `xml:"version,attr"`
	Channel      *Channel
}

// SetOptions sets options of given feed.
func (f *Feed) SetOptions(options ...func(f *Feed) error) error {
	for _, opt := range options {
		if err := opt(f); err != nil {
			return err
		}
	}
	return nil
}

// XML marshalls feed to XML string.
func (f *Feed) XML() (string, error) {
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Write writes marshalled XML to the given writer.
func (f *Feed) Write(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(f)
}

// Global pools for performance optimization
var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	stringBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

// WriteOptions defines options for XML writing
type WriteOptions struct {
	// BufferSize sets the initial buffer size for large feeds
	BufferSize int
	// UsePool enables buffer pooling for reduced allocations
	UsePool bool
}

// WriteWithOptions writes marshalled XML with performance options
func (f *Feed) WriteWithOptions(writer io.Writer, opts WriteOptions) error {
	var buf *bytes.Buffer

	// Get buffer if pooling is enabled or buffer size is specified
	if opts.UsePool {
		buf = bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)
	} else if opts.BufferSize > 0 {
		buf = &bytes.Buffer{}
		buf.Grow(opts.BufferSize)
	}

	finalWriter := writer
	if buf != nil {
		finalWriter = buf
	}

	// Write XML header
	if _, err := finalWriter.Write([]byte(xml.Header)); err != nil {
		return err
	}

	// Encode XML
	enc := xml.NewEncoder(finalWriter)
	enc.Indent("", "  ")
	if err := enc.Encode(f); err != nil {
		return err
	}

	// Copy buffer to final writer if buffering was used
	if buf != nil {
		if _, err := buf.WriteTo(writer); err != nil {
			return err
		}
	}

	return nil
}

// writeOpeningTags writes the RSS opening tags for streaming
func (f *Feed) writeOpeningTags(writer io.Writer) error {
	if _, err := writer.Write([]byte(xml.Header)); err != nil {
		return err
	}

	// Write opening tags manually for streaming
	if _, err := writer.Write([]byte(`<rss xmlns:itunes="`)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(f.ItunesXMLNS)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(`" xmlns:content="`)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(f.ContentXMLNS)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(`" version="`)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(f.Version)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(`">`)); err != nil {
		return err
	}
	if _, err := writer.Write([]byte("\n  <channel>\n")); err != nil {
		return err
	}
	return nil
}

// writeBasicChannelMetadata writes basic channel elements
func writeBasicChannelMetadata(enc *xml.Encoder, channel *Channel) error {
	if err := writeElement(enc, "title", channel.Title); err != nil {
		return err
	}
	if err := writeElement(enc, "link", channel.Link); err != nil {
		return err
	}
	if err := writeElement(enc, "description", channel.Description); err != nil {
		return err
	}
	if err := writeElement(enc, "language", channel.Language); err != nil {
		return err
	}
	if err := writeElement(enc, "copyright", channel.Copyright); err != nil {
		return err
	}
	return nil
}

// writeItunesChannelMetadata writes iTunes-specific channel elements
func writeItunesChannelMetadata(enc *xml.Encoder, channel *Channel) error {
	if channel.Author != "" {
		if err := writeElement(enc, "itunes:author", channel.Author); err != nil {
			return err
		}
	}
	if channel.Block != "" {
		if err := writeElement(enc, "itunes:block", channel.Block); err != nil {
			return err
		}
	}
	if channel.Explicit != "" {
		if err := writeElement(enc, "itunes:explicit", channel.Explicit); err != nil {
			return err
		}
	}
	if channel.Complete != "" {
		if err := writeElement(enc, "itunes:complete", channel.Complete); err != nil {
			return err
		}
	}
	if channel.NewFeedURL != "" {
		if err := writeElement(enc, "itunes:new-feed-url", channel.NewFeedURL); err != nil {
			return err
		}
	}
	if channel.Subtitle != "" {
		if err := writeElement(enc, "itunes:subtitle", channel.Subtitle); err != nil {
			return err
		}
	}
	return nil
}

// writeComplexChannelElements writes complex iTunes elements
func writeComplexChannelElements(enc *xml.Encoder, channel *Channel) error {
	// Write summary
	if channel.Summary != nil {
		summaryElement := struct {
			XMLName xml.Name `xml:"itunes:summary"`
			Value   string   `xml:",cdata"`
		}{
			XMLName: xml.Name{Local: "itunes:summary"},
			Value:   channel.Summary.Value,
		}
		if err := enc.Encode(summaryElement); err != nil {
			return err
		}
	}

	// Write owner
	if channel.Owner != nil {
		if err := enc.Encode(channel.Owner); err != nil {
			return err
		}
	}

	// Write image
	if channel.Image != nil {
		if err := enc.Encode(channel.Image); err != nil {
			return err
		}
	}

	// Write categories
	for _, category := range channel.Categories {
		if err := enc.Encode(category); err != nil {
			return err
		}
	}

	return nil
}

// writeClosingTags writes the RSS closing tags
func writeClosingTags(w io.Writer) error {
	if _, err := w.Write([]byte("\n  </channel>\n</rss>\n")); err != nil {
		return err
	}
	return nil
}

// StreamWrite writes XML in streaming fashion for large feeds
func (f *Feed) StreamWrite(writer io.Writer) error {
	// Write opening tags
	if err := f.writeOpeningTags(writer); err != nil {
		return err
	}

	// Prepare channel without items for metadata writing
	tempChannel := *f.Channel
	tempChannel.Items = nil

	enc := xml.NewEncoder(writer)
	enc.Indent("    ", "  ")

	// Write all channel metadata using helper functions
	if err := writeBasicChannelMetadata(enc, &tempChannel); err != nil {
		return err
	}
	if err := writeItunesChannelMetadata(enc, &tempChannel); err != nil {
		return err
	}
	if err := writeComplexChannelElements(enc, &tempChannel); err != nil {
		return err
	}

	// Stream items one by one to reduce memory usage
	for _, item := range f.Channel.Items {
		if err := enc.Encode(item); err != nil {
			return err
		}
	}

	// Write closing tags
	return writeClosingTags(writer)
}

// writeElement is a helper function for streaming XML writing
func writeElement(enc *xml.Encoder, name, value string) error {
	if value == "" {
		return nil
	}
	return enc.Encode(struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}{
		XMLName: xml.Name{Local: name},
		Value:   value,
	})
}

// XMLWithOptions generates XML string with performance options
func (f *Feed) XMLWithOptions(opts WriteOptions) (string, error) {
	var buf *bytes.Buffer

	if opts.UsePool {
		buf = bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)
	} else {
		buf = &bytes.Buffer{}
		if opts.BufferSize > 0 {
			buf.Grow(opts.BufferSize)
		}
	}

	if err := f.WriteWithOptions(buf, opts); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetBufferPool returns the global buffer pool for external use
func GetBufferPool() *sync.Pool {
	return &bufferPool
}

// GetStringBuilderPool returns the global string builder pool for external use
func GetStringBuilderPool() *sync.Pool {
	return &stringBuilderPool
}
