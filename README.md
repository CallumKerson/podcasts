# podcasts

[![Go Report Card](https://goreportcard.com/badge/github.com/CallumKerson/podcasts?style=flat-square)](https://goreportcard.com/report/github.com/CallumKerson/podcasts)
[![Go Reference](https://pkg.go.dev/badge/github.com/CallumKerson/podcasts.svg)](https://pkg.go.dev/github.com/CallumKerson/podcasts)
[![Release](https://img.shields.io/github/release/CallumKerson/podcasts.svg?style=flat-square)](https://github.com/CallumKerson/podcasts/releases/latest)

Podcast generator written in Go.

## Install

```bash
go get github.com/CallumKerson/podcasts
```

## Go Docs

[https://godoc.org/github.com/CallumKerson/podcasts](https://godoc.org/github.com/CallumKerson/podcasts)

## Example usage

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/CallumKerson/podcasts"
)

func main() {
	// initialize the podcast
	p := &podcasts.Podcast{
		Title:       "My podcast",
		Description: "This is my very simple podcast.",
		Language:    "EN",
		Link:        "http://www.example-podcast.com/my-podcast",
		Copyright:   "2015 My podcast copyright",
	}

	// add first podcast item
	p.AddItem(&podcasts.Item{
		Title:    "Episode 1",
		GUID:     "http://www.example-podcast.com/my-podcast/1/episode-one",
		PubDate:  podcasts.NewPubDate(time.Now()),
		Duration: podcasts.NewDuration(time.Second * 230),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/1/episode.mp3",
			Length: "12312",
			Type:   "MP3",
		},
	})

	// add second podcast item
	p.AddItem(&podcasts.Item{
		Title:   "Episode 2",
		GUID:    "http://www.example-podcast.com/my-podcast/2/episode-two",
		PubDate: podcasts.NewPubDate(time.Now()),
		Duration: podcasts.NewDuration(time.Second * 320),
		Enclosure: &podcasts.Enclosure{
			URL:    "http://www.example-podcast.com/my-podcast/2/episode.mp3",
			Length: "46732",
			Type:   "MP3",
		},
	})

	// get podcast feed, you can pass options to customize it
	feed, err := p.Feed(
		podcasts.Author("Author Name"),
		podcasts.Block,
		podcasts.Explicit,
		podcasts.Complete,
		podcasts.NewFeedURL("http://www.example-podcast.com/new-feed-url"),
		podcasts.Subtitle("This is my very simple podcast subtitle."),
		podcasts.Summary("This is my very simple podcast summary."),
		podcasts.Owner("Podcast Owner", "owner@example-podcast.com"),
		podcasts.Image("http://www.example-podcast.com/my-podcast.jpg"),
	)

	// handle error
	if err != nil {
		log.Fatal(err)
	}

	// finally write the xml to any io.Writer
	feed.Write(os.Stdout)
}
```

Which gives us this XML output:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" version="2.0">
  <channel>
    <title>My podcast</title>
    <link>http://www.example-podcast.com/my-podcast</link>
    <copyright>2015 My podcast copyright</copyright>
    <language>EN</language>
    <description>This is my very simple podcast.</description>
    <itunes:author>Author Name</itunes:author>
    <itunes:block>yes</itunes:block>
    <itunes:explicit>yes</itunes:explicit>
    <itunes:complete>yes</itunes:complete>
    <itunes:new-feed-url>http://www.example-podcast.com/new-feed-url</itunes:new-feed-url>
    <itunes:subtitle>This is my very simple podcast subtitle.</itunes:subtitle>
    <itunes:summary><![CDATA[This is my very simple podcast summary.]]></itunes:summary>
    <itunes:owner>
      <itunes:name>Podcast Owner</itunes:name>
      <itunes:email>owner@example-podcast.com</itunes:email>
    </itunes:owner>
    <itunes:image href="http://www.example-podcast.com/my-podcast.jpg"></itunes:image>
    <item>
      <title>Episode 1</title>
      <guid>http://www.example-podcast.com/my-podcast/1/episode-one</guid>
      <pubDate>Tue, 10 Nov 2009 23:00:00 +0000</pubDate>
      <itunes:duration>3:50</itunes:duration>
      <enclosure url="http://www.example-podcast.com/my-podcast/1/episode.mp3" length="12312" type="MP3"></enclosure>
    </item>
    <item>
      <title>Episode 2</title>
      <guid>http://www.example-podcast.com/my-podcast/2/episode-two</guid>
      <pubDate>Tue, 10 Nov 2009 23:00:00 +0000</pubDate>
      <itunes:duration>5:20</itunes:duration>
      <enclosure url="http://www.example-podcast.com/my-podcast/2/episode.mp3" length="46732" type="MP3"></enclosure>
    </item>
  </channel>
</rss>
```

## Options

For further options and configuration, please see the [options docs](./docs/options.md).
