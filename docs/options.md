# Configuration Options

This document describes the various ways to configure the podcasts library,
including feed metadata options and performance optimizations.

## Feed Configuration Options

Feed options are functions that modify the RSS feed during creation.
They are passed to the `podcast.Feed()` method:

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

## Writing the Feed

There are two ways to serialise a configured feed.

### `feed.XML() (string, error)`

Marshals the feed and returns it as a string.

```go
xmlString, err := feed.XML()
```

**When to use:** When you need the feed in memory, for example to hand to a template or store it.

**Downsides:** Holds the whole feed in memory as a string.

### `feed.Write(w io.Writer) error`

Marshals the feed straight to any `io.Writer`, such as an `http.ResponseWriter` or a file.

```go
err := feed.Write(w)
```

**When to use:** Whenever you are writing the feed somewhere rather than inspecting it.
This is the cheaper of the two, as no intermediate string is built.

**Downsides:** None; prefer this unless you specifically need a string.

### Buffering

`Write` issues a number of small writes as it encodes.
If the destination is expensive to write to, such as a file or a socket, wrap it:

```go
bw := bufio.NewWriter(w)
if err := feed.Write(bw); err != nil {
    return err
}
return bw.Flush()
```

The library deliberately does not provide its own buffering or pooling options.
Earlier versions did, and they measured no faster than a `bufio.Writer` while adding public API that could not be changed later.
