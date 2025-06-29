# Configuration Options

This document describes the various ways to configure the podcasts library,
including feed metadata options and performance optimizations.

## Feed Configuration Options

Feed options are functions that modify the RSS feed during creation. They are passed to the `podcast.Feed()` method:

```go
feed, err := podcast.Feed(
    Author("John Doe"),
    Explicit,
    Owner("John Doe", "john@example.com"),
)
```

### Metadata Options

#### `Author(name string)`

Sets the `itunes:author` field for the podcast feed.

**When to use:** When you want to specify the primary author/creator of the podcast.

**Example:**

```go
feed, err := podcast.Feed(Author("Jane Smith"))
```

**Downsides:** None - this is purely metadata.

#### `Subtitle(subtitle string)`

Sets the `itunes:subtitle` field for a brief description.

**When to use:** When you need a short, descriptive subtitle that appears in podcast directories.

**Example:**

```go
feed, err := podcast.Feed(Subtitle("Weekly tech discussions"))
```

**Downsides:** Should be kept brief (recommended under 255 characters).

#### `Summary(summary string)`

Sets the `itunes:summary` field with a longer description (supports HTML via CDATA).

**When to use:** When you need a detailed description of your podcast with HTML formatting.

**Example:**

```go
feed, err := podcast.Feed(Summary(`<p>A detailed podcast about <strong>technology</strong> trends.</p>`))
```

**Downsides:** HTML content increases feed size slightly.

#### `Owner(name, email string)`

Sets the `itunes:owner` information for the podcast.

**When to use:** Always recommended - this identifies the podcast owner in iTunes/Apple Podcasts.

**Example:**

```go
feed, err := podcast.Feed(Owner("Jane Smith", "jane@example.com"))
```

**Downsides:** Email address will be visible in the RSS feed.

#### `Image(url string)`

Sets the podcast artwork URL.

**When to use:** Always recommended - podcast directories require artwork.

**Example:**

```go
feed, err := podcast.Feed(Image("https://example.com/podcast-art.jpg"))
```

**Downsides:**

- URL must be absolute (validated at runtime)
- Large images may slow feed parsing
- Recommended: 1400x1400 to 3000x3000 pixels, JPEG/PNG

### Content Control Options

#### `Block`

Sets `itunes:block` to "yes", preventing the podcast from appearing in iTunes.

**When to use:** When you want to remove your podcast from iTunes/Apple Podcasts directory.

**Example:**

```go
feed, err := podcast.Feed(Block)
```

**Downsides:** Podcast won't be discoverable through Apple Podcasts.

#### `Explicit`

Marks the podcast as containing explicit content.

**When to use:** When your podcast contains adult language or mature themes.

**Example:**

```go
feed, err := podcast.Feed(Explicit)
```

**Downsides:** May limit discoverability or require age verification.

#### `Complete`

Marks the podcast as complete (no more episodes will be added).

**When to use:** When you're ending your podcast permanently.

**Example:**

```go
feed, err := podcast.Feed(Complete)
```

**Downsides:** Podcast clients may stop checking for new episodes.

#### `NewFeedURL(url string)`

Redirects podcast clients to a new feed URL.

**When to use:** When migrating your podcast to a new hosting provider or URL.

**Example:**

```go
feed, err := podcast.Feed(NewFeedURL("https://newhost.com/my-podcast/feed.xml"))
```

**Downsides:**

- URL must be absolute (validated at runtime)
- Some clients may not support feed redirection immediately

## Performance Options

Performance options are configured through the `WriteOptions` struct when generating XML:

### `WriteOptions` Structure

```go
type WriteOptions struct {
    BufferSize int  // Initial buffer size in bytes
    UsePool    bool // Enable buffer pooling
}
```

### Standard XML Generation

```go
// Basic usage - no performance optimizations
xmlString, err := feed.XML()
err = feed.Write(writer)
```

**When to use:** For small podcasts (< 50 episodes) or when memory usage isn't a concern.

**Downsides:** May allocate more memory for large feeds.

### Buffered Writing

```go
opts := WriteOptions{BufferSize: 8192}
err = feed.WriteWithOptions(writer, opts)
xmlString, err := feed.XMLWithOptions(opts)
```

**When to use:**

- Large podcasts (100+ episodes)
- When you know approximate feed size
- Network writing where buffering improves performance

**Benefits:**

- Reduces memory allocations
- Improves write performance for large feeds
- Reduces system calls

**Downsides:**

- Uses more memory upfront
- Buffer size estimation required for optimal performance

### Pooled Buffers

```go
opts := WriteOptions{UsePool: true}
err = feed.WriteWithOptions(writer, opts)
```

**When to use:**

- High-frequency feed generation (multiple feeds per second)
- Server applications generating many feeds
- Memory-constrained environments

**Benefits:**

- Reuses buffer objects across operations
- Reduces garbage collection pressure
- Lower memory allocation overhead

**Downsides:**

- Slightly more complex memory management
- Buffers remain allocated in pool between uses

### Combined Optimization

```go
opts := WriteOptions{
    BufferSize: 16384,
    UsePool:    true,
}
err = feed.WriteWithOptions(writer, opts)
```

**When to use:** High-performance applications with large feeds and frequent generation.

**Benefits:** Maximum performance optimization combining both techniques.

**Downsides:** Highest memory usage per operation.

### Streaming for Large Feeds

```go
err = feed.StreamWrite(writer)
```

**When to use:**

- Very large podcasts (500+ episodes)
- Memory-constrained environments
- When feed size approaches available memory

**Benefits:**

- Constant memory usage regardless of feed size
- Suitable for extremely large feeds
- Lower peak memory consumption

**Downsides:**

- Slightly slower than buffered approaches
- Cannot be used to generate XML strings (only direct writing)
- Less control over output formatting

## Performance Comparison

| Method                          | Memory Usage | Speed   | Best For                         |
| ------------------------------- | ------------ | ------- | -------------------------------- |
| `XML()` / `Write()`             | Medium       | Fast    | Small feeds (< 50 episodes)      |
| `WriteWithOptions` (BufferSize) | High         | Fastest | Large feeds, known size          |
| `WriteWithOptions` (UsePool)    | Low          | Fast    | High frequency generation        |
| `WriteWithOptions` (Both)       | Medium-High  | Fastest | High performance + large feeds   |
| `StreamWrite`                   | Lowest       | Medium  | Very large feeds (500+ episodes) |

## Buffer Pool Access

For advanced use cases, you can access the global buffer pools directly:

```go
import "sync"

// Get the global buffer pool
bufferPool := GetBufferPool()
buf := bufferPool.Get().(*bytes.Buffer)
defer func() {
    buf.Reset()
    bufferPool.Put(buf)
}()

// Get the global string builder pool
stringPool := GetStringBuilderPool()
sb := stringPool.Get().(*strings.Builder)
defer func() {
    sb.Reset()
    stringPool.Put(sb)
}()
```

**When to use:** When integrating with existing pooling strategies or building custom optimizations.

**Downsides:** Requires manual memory management and proper cleanup.
