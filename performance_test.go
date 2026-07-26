package podcasts

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
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

// Benchmark using buffer pool vs regular allocation
func BenchmarkBufferPool(b *testing.B) {
	podcast := createTestPodcast(50)
	feed, err := podcast.Feed()
	if err != nil {
		b.Fatal(err)
	}

	b.Run("RegularAlloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			if err := feed.Write(&buf); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("BufferPool", func(b *testing.B) {
		pool := GetBufferPool()
		for i := 0; i < b.N; i++ {
			buf := pool.Get().(*bytes.Buffer)
			if err := feed.Write(buf); err != nil {
				b.Fatal(err)
			}
			buf.Reset()
			pool.Put(buf)
		}
	})
}

// Benchmark WriteWithOptions with different option combinations
func BenchmarkWriteWithOptions(b *testing.B) {
	podcast := createTestPodcast(50)
	feed, err := podcast.Feed()
	if err != nil {
		b.Fatal(err)
	}

	optionCombinations := map[string]WriteOptions{
		"NoOptions":     {},
		"UsePool":       {UsePool: true},
		"BufferSize":    {BufferSize: 8192},
		"PoolAndBuffer": {UsePool: true, BufferSize: 8192},
	}

	for name, opts := range optionCombinations {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				if err := feed.WriteWithOptions(&buf, opts); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark StreamWrite performance
func BenchmarkStreamWrite(b *testing.B) {
	podcast := createTestPodcast(100)
	feed, err := podcast.Feed()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := feed.StreamWrite(&buf); err != nil {
			b.Fatal(err)
		}
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
			Title:       fmt.Sprintf("Episode %d: Performance Testing", episodeNum),
			GUID:        fmt.Sprintf("https://example.com/episode-%d", episodeNum),
			PubDate:     NewPubDate(time.Date(2024, 1, (episodeNum%28)+1, 12, 0, 0, 0, time.UTC)),
			Duration:    NewDuration(time.Minute * time.Duration(20+(episodeNum%40))),
			Description: &CDATAText{Value: fmt.Sprintf("This is episode %d of our performance testing podcast. It contains detailed information about performance optimization.", episodeNum)},
			Enclosure: &Enclosure{
				URL:    fmt.Sprintf("https://example.com/episodes/episode-%d.mp3", episodeNum),
				Length: fmt.Sprintf("%d", 1000000+(episodeNum*50000)),
				Type:   "audio/mpeg",
			},
		})
	}

	return podcast
}

// StringPool provides a pool of reusable strings.Builder objects
type StringPool struct {
	pool sync.Pool
}

// NewStringPool creates a new string pool
func NewStringPool() *StringPool {
	return &StringPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

// Get retrieves a strings.Builder from the pool
func (sp *StringPool) Get() *strings.Builder {
	return sp.pool.Get().(*strings.Builder)
}

// Put returns a strings.Builder to the pool after resetting it
func (sp *StringPool) Put(sb *strings.Builder) {
	sb.Reset()
	sp.pool.Put(sb)
}

// BufferPool provides a pool of reusable bytes.Buffer objects
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new buffer pool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Get retrieves a bytes.Buffer from the pool
func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Put returns a bytes.Buffer to the pool after resetting it
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	bp.pool.Put(buf)
}

// Test performance-oriented functionality
func TestPerformanceMethods(t *testing.T) {
	podcast := &Podcast{
		Title:       "Performance Methods Test",
		Description: "Testing performance methods",
		Language:    "en",
		Link:        "https://example.com/perf",
		Copyright:   "2024",
	}

	t.Run("AddItemWithCapacity", func(t *testing.T) {
		item := &Item{
			Title:   "Capacity Test Episode",
			GUID:    "https://example.com/perf/capacity",
			PubDate: NewPubDate(time.Now()),
		}

		// Test adding with capacity hint
		podcast.AddItemWithCapacity(item, 100)

		if podcast.GetItemCount() != 1 {
			t.Errorf("Expected 1 item, got %d", podcast.GetItemCount())
		}

		items := podcast.GetItems()
		if len(items) != 1 {
			t.Errorf("Expected 1 item in copy, got %d", len(items))
		}

		if items[0].Title != "Capacity Test Episode" {
			t.Error("Item title doesn't match")
		}
	})

	t.Run("GetItemsSlice", func(t *testing.T) {
		directSlice := podcast.GetItemsSlice()
		if len(directSlice) != 1 {
			t.Errorf("Expected 1 item in direct slice, got %d", len(directSlice))
		}
	})

	t.Run("GetItemCount", func(t *testing.T) {
		count := podcast.GetItemCount()
		expectedCount := 1
		if count != expectedCount {
			t.Errorf("Expected %d items, got %d", expectedCount, count)
		}

		// Add another item to test count increase
		podcast.AddItem(&Item{
			Title:   "Second Episode",
			GUID:    "https://example.com/perf/second",
			PubDate: NewPubDate(time.Now()),
		})

		count = podcast.GetItemCount()
		expectedCount = 2
		if count != expectedCount {
			t.Errorf("Expected %d items after adding second, got %d", expectedCount, count)
		}
	})
}
