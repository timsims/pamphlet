package pamphlet

import (
	"errors"
)

const (
	HTMLMimetype = "application/xhtml+xml"
	SVGMimetype  = "image/svg+xml"
)

type Book struct {
	Title         string
	Author        string
	Language      string
	Identifier    string
	Publisher     string
	Description   string
	Subject       string
	Date          string
	Chapters      []Chapter
	ManifestItems []ManifestItem
}

type Chapter struct {
	ZipFile
	// ID is the idref attribute in the spine tag
	ID        string
	Title     string
	Href      string
	MediaType string
	HasToc    bool
	// order of the chapter in the book, based on the toc file's playerOrder attribute
	Order int
}

func (c *Chapter) GetContent() (string, error) {
	if c.MediaType != HTMLMimetype && c.MediaType != SVGMimetype {
		return "", errors.New("unsupported media type")
	}
	content, err := c.GetRawContent()

	if err != nil {
		return "", err
	}

	return string(content), nil
}

type ManifestItem struct {
	ZipFile
	ID        string
	Href      string
	RealPath  string
	MediaType string
}
