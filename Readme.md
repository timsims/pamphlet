# Pamphlet

Pamphlet is a Go library for parsing EPUB files. It extracts metadata, table of contents, and other relevant information from EPUB files.

## Features

- Parse EPUB metadata
- Extract table of contents
- Retrieve chapters and manifest items

## Installation

To install the library, use `go get`:

```sh
go get github.com/timsims/pamphlet
```

## Usage

Here is an example of how to use the Pamphlet library:

```go
package main

import (
    "fmt"
    "log"
    "github.com/timsims/pamphlet"
)

func main() {
    parser, err := pamphlet.NewParser("path/to/your.epub")
    if err != nil {
        log.Fatal(err)
    }

    book := parser.GetBook()
    fmt.Printf("Title: %s\n", book.Title)
    fmt.Printf("Author: %s\n", book.Author)
    
    chapters := book.Chapters
    for i, chapter := range chapters {
        fmt.Printf("Chapter %d: %s\n", i+1, chapter.Title)
        fmt.Printf("Content: %s\n", chapter.GetContent())
    }
    
    manifestItems := book.ManifestItems
    for i, item := range manifestItems {
        fmt.Printf("Item Media Type: %s\n", i+1, item.MediaType)
        fmt.Printf("Item %d: %s\n", i+1, item.Href)
    }
}
```

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## Project Structure

- `parser.go`: Contains the main logic for parsing EPUB files.
- `parser_test.go`: Contains tests for the parser.
- `epub/`: Contains the EPUB-related structs and types.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
```