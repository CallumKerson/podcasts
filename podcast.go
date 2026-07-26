package podcasts

// Podcast represents a web podcast.
type Podcast struct {
	Title       string
	Description string
	Link        string
	Language    string
	Copyright   string
	items       []*Item
}

// AddItem adds an item to the podcast.
func (p *Podcast) AddItem(item *Item) {
	p.items = append(p.items, item)
}

// Feed creates a new feed for current podcast.
func (p *Podcast) Feed(options ...func(f *Feed) error) (*Feed, error) {
	feed := &Feed{
		ItunesXMLNS:  itunesXMLNS,
		ContentXMLNS: contentXMLNS,
		Version:      rssVersion,
		Channel: &Channel{
			Title:       p.Title,
			Description: p.Description,
			Link:        p.Link,
			Copyright:   p.Copyright,
			Language:    p.Language,
			Items:       p.items,
		},
	}
	err := feed.SetOptions(options...)
	return feed, err
}

// Len returns the number of items in the podcast.
func (p *Podcast) Len() int {
	return len(p.items)
}

// Items returns a copy of the podcast items, so that callers cannot reach the
// slice the podcast builds its feed from. A Podcast is not safe for concurrent
// use: guard it yourself if one goroutine may call AddItem while another reads.
func (p *Podcast) Items() []*Item {
	items := make([]*Item, len(p.items))
	copy(items, p.items)
	return items
}
