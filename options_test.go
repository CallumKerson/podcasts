package podcasts

import (
	"errors"
	"testing"
)

const (
	testAuthor   = "Test Author"
	testSubtitle = "Test Subtitle"
)

func TestAuthor(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	author := "john"
	if err := Author(author)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if author != feed.Channel.Author {
		t.Errorf("expected %v got %v", author, feed.Channel.Author)
	}
}

func TestBlock(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := Block(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if ValueYes != feed.Channel.Block {
		t.Errorf("expected %v got %v", ValueYes, feed.Channel.Block)
	}
}

func TestExplicit(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := Explicit(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if ValueYes != feed.Channel.Explicit {
		t.Errorf("expected %v got %v", ValueYes, feed.Channel.Explicit)
	}
}

func TestComplete(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := Complete(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if ValueYes != feed.Channel.Complete {
		t.Errorf("expected %v got %v", ValueYes, feed.Channel.Complete)
	}
}

func TestNewFeedURL(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	newURL := "http://example.com/test"
	if err := NewFeedURL(newURL)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if newURL != feed.Channel.NewFeedURL {
		t.Errorf("expected %v got %v", ValueYes, feed.Channel.NewFeedURL)
	}
}

func TestNewFeedURLInvalid(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}

	// Test invalid URL
	newURL := "invalid url"
	if err := NewFeedURL(newURL)(feed); err == nil {
		t.Error("expected error for invalid URL")
	}

	// Test relative URL (should be rejected)
	relativeURL := "/relative/path"
	if err := NewFeedURL(relativeURL)(feed); !errors.Is(err, ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}

	// Test URL with characters that make it unparseable
	invalidCharURL := "http://example.com/path\x00\x01"
	if err := NewFeedURL(invalidCharURL)(feed); err == nil {
		t.Error("expected error for URL with null characters")
	}
}

func TestSubtitle(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	subtitle := "this is subtitle"
	if err := Subtitle(subtitle)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if subtitle != feed.Channel.Subtitle {
		t.Errorf("expected %v got %v", subtitle, feed.Channel.Subtitle)
	}
}

func TestSummary(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	summary := `this is summary. <a href="http://example.com/more">more</a>`
	if err := Summary(summary)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if summary != feed.Channel.Summary.Value {
		t.Errorf("expected %v got %v", summary, feed.Channel.Summary.Value)
	}
}

func TestOwner(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	name := "anabelle"
	email := "test@test.com"
	if err := Owner(name, email)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if name != feed.Channel.Owner.Name {
		t.Errorf("expected %v got %v", name, feed.Channel.Owner.Name)
	}
	if email != feed.Channel.Owner.Email {
		t.Errorf("expected %v got %v", email, feed.Channel.Owner.Email)
	}
}

func TestImage(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	href := "http://example.com/test/image.jpg"
	if err := Image(href)(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if href != feed.Channel.Image.Href {
		t.Errorf("expected %v got %v", href, feed.Channel.Image.Href)
	}
}

func TestCategory(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := Category("Technology")(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(feed.Channel.Categories) != 1 {
		t.Fatalf("expected 1 category got %v", len(feed.Channel.Categories))
	}
	if want := "Technology"; feed.Channel.Categories[0].Text != want {
		t.Errorf("expected %v got %v", want, feed.Channel.Categories[0].Text)
	}
	if got := feed.Channel.Categories[0].Categories; len(got) != 0 {
		t.Errorf("expected no subcategories got %v", len(got))
	}
}

func TestCategoryWithSubcategories(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := Category("Society & Culture", "Documentary", "Personal Journals")(feed); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(feed.Channel.Categories) != 1 {
		t.Fatalf("expected 1 category got %v", len(feed.Channel.Categories))
	}

	subcategories := feed.Channel.Categories[0].Categories
	if len(subcategories) != 2 {
		t.Fatalf("expected 2 subcategories got %v", len(subcategories))
	}
	for i, want := range []string{"Documentary", "Personal Journals"} {
		if subcategories[i].Text != want {
			t.Errorf("expected %v got %v", want, subcategories[i].Text)
		}
	}
}

func TestCategoryAppends(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}
	if err := feed.SetOptions(Category("Technology"), Category("News", "Tech News")); err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(feed.Channel.Categories) != 2 {
		t.Fatalf("expected 2 categories got %v", len(feed.Channel.Categories))
	}
	for i, want := range []string{"Technology", "News"} {
		if feed.Channel.Categories[i].Text != want {
			t.Errorf("expected %v got %v", want, feed.Channel.Categories[i].Text)
		}
	}
}

func TestCategoryInvalid(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}

	if err := Category("")(feed); !errors.Is(err, ErrInvalidCategory) {
		t.Errorf("expected ErrInvalidCategory, got %v", err)
	}

	if err := Category("Technology", "")(feed); !errors.Is(err, ErrInvalidCategory) {
		t.Errorf("expected ErrInvalidCategory for empty subcategory, got %v", err)
	}

	if len(feed.Channel.Categories) != 0 {
		t.Errorf("expected no categories to be added, got %v", len(feed.Channel.Categories))
	}
}

func TestImageInvalid(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}

	// Test invalid URL
	href := "invalid img url"
	if err := Image(href)(feed); err == nil {
		t.Error("expected error for invalid image URL")
	}

	// Test relative URL (should be rejected)
	relativeHref := "/relative/image.jpg"
	if err := Image(relativeHref)(feed); !errors.Is(err, ErrInvalidImage) {
		t.Errorf("expected ErrInvalidImage, got %v", err)
	}

	// Test URL with characters that make it unparseable
	invalidCharHref := "http://example.com/image\x00\x01.jpg"
	if err := Image(invalidCharHref)(feed); err == nil {
		t.Error("expected error for image URL with null characters")
	}
}

// TestSetOptionsWithMultipleOptions tests applying multiple options
func TestSetOptionsWithMultipleOptions(t *testing.T) {
	feed := &Feed{
		Channel: &Channel{},
	}

	// Apply multiple options
	err := feed.SetOptions(
		Author(testAuthor),
		Block,
		Explicit,
		Subtitle(testSubtitle),
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if feed.Channel.Author != testAuthor {
		t.Errorf("expected author to be set")
	}
	if feed.Channel.Block != ValueYes {
		t.Errorf("expected block to be set")
	}
	if feed.Channel.Explicit != ValueYes {
		t.Errorf("expected explicit to be set")
	}
	if feed.Channel.Subtitle != testSubtitle {
		t.Errorf("expected subtitle to be set")
	}
}
