# Changelog

## [2.1.0](https://github.com/CallumKerson/podcasts/compare/v2.0.0...v2.1.0) (2026-07-26)


### Features

* add Category option ([#37](https://github.com/CallumKerson/podcasts/issues/37)) ([e3fa10c](https://github.com/CallumKerson/podcasts/commit/e3fa10c7c00758e1b5eca9cf753ebced0e705035))

## [2.0.0](https://github.com/CallumKerson/podcasts/compare/v1.1.0...v2.0.0) (2026-07-26)


### ⚠ BREAKING CHANGES

* the module path is now github.com/CallumKerson/podcasts/v2, so imports need the /v2 suffix.
* Podcast.GetItemCount is renamed to Podcast.Len, Podcast.GetItems is renamed to Podcast.Items, and Podcast.GetItemsSlice is removed in favour of Podcast.Items.
* Go 1.25 or newer is now required to build this module.
* Feed.StreamWrite, Feed.WriteWithOptions, Feed.XMLWithOptions, WriteOptions, GetBufferPool, GetStringBuilderPool and Podcast.AddItemWithCapacity are removed. Use Feed.Write or Feed.XML, wrapping the writer in a bufio.Writer where buffering is wanted, and Podcast.AddItem.

### Bug Fixes

* correct negative durations and channel element ordering ([#30](https://github.com/CallumKerson/podcasts/issues/30)) ([7dc4718](https://github.com/CallumKerson/podcasts/commit/7dc471806281d753e703f1abbfd638919cb2b575))


### Code Refactoring

* give the podcast accessors idiomatic names ([#34](https://github.com/CallumKerson/podcasts/issues/34)) ([b77bf10](https://github.com/CallumKerson/podcasts/commit/b77bf10ae3623a48042f4b9265d7ba47c1fd8b6b))
* remove the performance API ([#32](https://github.com/CallumKerson/podcasts/issues/32)) ([2b2b0a0](https://github.com/CallumKerson/podcasts/commit/2b2b0a0d61199d2bf4b80bb543ac4ea989da916e))


### Build System

* move the module to the /v2 import path ([#35](https://github.com/CallumKerson/podcasts/issues/35)) ([747311f](https://github.com/CallumKerson/podcasts/commit/747311f663ba40c095aec1fd11e171a76717d37a))
* require Go 1.25 ([#36](https://github.com/CallumKerson/podcasts/issues/36)) ([95246e3](https://github.com/CallumKerson/podcasts/commit/95246e3989e8cfa0b5445e04e18d62c93159e6d1))
