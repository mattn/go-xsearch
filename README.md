# go-xsearch

A Go library for searching text in X. This provides a simple way to search for patterns in X.

## Installation

To use `go-xsearch` in your Go project, install it with:

```bash
go get github.com/mattn/go-xsearch
```

Ensure you have Go installed and configured properly.

## Usage

The primary functionality of `go-xsearch` is provided by the `xsearch.Search` function, which searches for a given pattern in X.

```go
entries, err := xsearch.Search("ぬるぽ")
if err != nil {
    log.Fatal(err)
}

for _, entry := range entries {
    println(entry.ID)
}
```

## reply-bot

The `reply-bot` is a bot designed to respond with a fixed message to statements matching a specified pattern.

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
