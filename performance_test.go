package podcasts

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

// Benchmark current XML generation performance
func BenchmarkFeedGeneration(b *testing.B) {
	podcast := createTestPodcast(100) // 100 episodes

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feed, err := podcast.Feed()
		if err != nil {
			b.Fatal(err)
		}
		_, err = feed.XML()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark XML generation with different feed sizes
func BenchmarkFeedGenerationSizes(b *testing.B) {
	sizes := []int{10, 50, 100, 500, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Episodes_%d", size), func(b *testing.B) {
			podcast := createTestPodcast(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				feed, err := podcast.Feed()
				if err != nil {
					b.Fatal(err)
				}
				_, err = feed.XML()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark XML writing to different writers
func BenchmarkFeedWrite(b *testing.B) {
	podcast := createTestPodcast(100)
	feed, err := podcast.Feed()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			if err := feed.Write(&buf); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Discard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := feed.Write(io.Discard); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark memory allocations during feed creation
func BenchmarkFeedCreationAllocs(b *testing.B) {
	podcast := createTestPodcast(100)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := podcast.Feed()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark AddItem performance
func BenchmarkAddItem(b *testing.B) {
	podcast := &Podcast{
		Title:       "Benchmark Podcast",
		Description: "Testing item addition performance",
		Language:    "en",
		Link:        "https://example.com",
		Copyright:   "2024",
	}

	item := &Item{
		Title:    "Test Episode",
		GUID:     "https://example.com/test",
		PubDate:  NewPubDate(time.Now()),
		Duration: NewDuration(time.Minute * 30),
		Enclosure: &Enclosure{
			URL:    "https://example.com/test.mp3",
			Length: "12345678",
			Type:   "audio/mpeg",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		podcast.AddItem(item)
	}
}

// Helper function to create test podcast with specified number of episodes
func createTestPodcast(episodeCount int) *Podcast {
	podcast := &Podcast{
		Title:       "Performance Test Podcast",
		Description: "A podcast for performance testing",
		Language:    "en-US",
		Link:        "https://example.com/podcast",
		Copyright:   "2024 Test Corp",
	}

	for episodeNum := 1; episodeNum <= episodeCount; episodeNum++ {
		podcast.AddItem(&Item{
			Title:    fmt.Sprintf("Episode %d: Performance Testing", episodeNum),
			GUID:     fmt.Sprintf("https://example.com/episode-%d", episodeNum),
			PubDate:  NewPubDate(time.Date(2024, 1, (episodeNum%28)+1, 12, 0, 0, 0, time.UTC)),
			Duration: NewDuration(time.Minute * time.Duration(20+(episodeNum%40))),
			Description: &CDATAText{
				Value: fmt.Sprintf(
					"This is episode %d of our performance testing podcast. It contains detailed information about performance optimization.",
					episodeNum,
				),
			},
			Enclosure: &Enclosure{
				URL:    fmt.Sprintf("https://example.com/episodes/episode-%d.mp3", episodeNum),
				Length: fmt.Sprintf("%d", 1000000+(episodeNum*50000)),
				Type:   "audio/mpeg",
			},
		})
	}

	return podcast
}
