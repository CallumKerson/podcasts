package podcasts

import (
	"errors"
	"net/url"
)

var (
	// ErrInvalidURL represents a error returned for invalid url.
	ErrInvalidURL = errors.New("podcasts: invalid url")

	// ErrInvalidImage represents a error returned for invalid image.
	ErrInvalidImage = errors.New("podcasts: invalid image")

	// ErrInvalidCategory represents a error returned for invalid category.
	ErrInvalidCategory = errors.New("podcasts: invalid category")
)

const (
	// ValueYes represents positive value used in XML feed.
	ValueYes = "yes"
)

// Author sets itunes:author of given feed.
func Author(author string) func(f *Feed) error {
	return func(f *Feed) error {
		f.Channel.Author = author
		return nil
	}
}

// Block enables itunes:block of given feed.
func Block(f *Feed) error {
	f.Channel.Block = ValueYes
	return nil
}

// Explicit enables itunes:explicit of given feed.
func Explicit(f *Feed) error {
	f.Channel.Explicit = ValueYes
	return nil
}

// Complete enables itunes:complete of given feed.
func Complete(f *Feed) error {
	f.Channel.Complete = ValueYes
	return nil
}

// NewFeedURL sets itunes:new-feed-url of given feed.
func NewFeedURL(newURL string) func(feed *Feed) error {
	return func(feed *Feed) error {
		u, err := url.Parse(newURL)
		if err != nil {
			return err
		}
		if !u.IsAbs() {
			return ErrInvalidURL
		}
		feed.Channel.NewFeedURL = newURL
		return nil
	}
}

// Subtitle sets itunes:subtitle of given feed.
func Subtitle(subtitle string) func(feed *Feed) error {
	return func(feed *Feed) error {
		feed.Channel.Subtitle = subtitle
		return nil
	}
}

// Summary sets itunes:summary of given feed.
func Summary(summary string) func(f *Feed) error {
	return func(f *Feed) error {
		f.Channel.Summary = &CDATAText{summary}
		return nil
	}
}

// Owner sets itunes:owner of given feed.
func Owner(name, email string) func(feed *Feed) error {
	return func(feed *Feed) error {
		feed.Channel.Owner = &ItunesOwner{
			Name:  name,
			Email: email,
		}
		return nil
	}
}

// Category adds an itunes:category to the given feed, nesting any
// subcategories inside it. A feed may carry more than one category, so each
// call appends rather than replacing what is already there.
//
// The names are not checked against Apple's category list, which changes
// without notice; a feed submitted to a directory needs to use the names that
// directory publishes.
func Category(name string, subcategories ...string) func(feed *Feed) error {
	return func(feed *Feed) error {
		if name == "" {
			return ErrInvalidCategory
		}
		category := &ItunesCategory{Text: name}
		for _, subcategory := range subcategories {
			if subcategory == "" {
				return ErrInvalidCategory
			}
			category.Categories = append(category.Categories, &ItunesCategory{Text: subcategory})
		}
		feed.Channel.Categories = append(feed.Channel.Categories, category)
		return nil
	}
}

// Image sets itunes:image of given feed.
func Image(href string) func(feed *Feed) error {
	return func(feed *Feed) error {
		u, err := url.Parse(href)
		if err != nil {
			return err
		}
		if !u.IsAbs() {
			return ErrInvalidImage
		}
		feed.Channel.Image = &ItunesImage{
			Href: href,
		}
		return nil
	}
}
