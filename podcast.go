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

// AddItemWithCapacity adds an item with pre-allocated capacity hint
func (p *Podcast) AddItemWithCapacity(item *Item, expectedTotal int) {
	if cap(p.items) < expectedTotal {
		newItems := make([]*Item, len(p.items), expectedTotal)
		copy(newItems, p.items)
		p.items = newItems
	}
	p.items = append(p.items, item)
}

// GetItemCount returns the number of items in the podcast
func (p *Podcast) GetItemCount() int {
	return len(p.items)
}

// GetItems returns a copy of the items slice (safe for concurrent access)
func (p *Podcast) GetItems() []*Item {
	items := make([]*Item, len(p.items))
	copy(items, p.items)
	return items
}

// GetItemsSlice returns the items slice directly (unsafe for concurrent modification)
func (p *Podcast) GetItemsSlice() []*Item {
	return p.items
}
