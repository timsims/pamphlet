package pamphlet

import (
	"errors"
	"strings"
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

func (b *Book) GetChaptersSize() int {
	size := 0
	for _, c := range b.Chapters {
		if c.HasToc {
			size++
		}
	}
	return size
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

	contentStr := string(content)
	contentStr = removeTags(contentStr, "script")
	contentStr = removeTags(contentStr, "style")
	contentStr = removeEventHandlers(contentStr)
	contentStr = strings.ReplaceAll(contentStr, "\u0000", "\n")

	return contentStr, nil
}

func removeTags(str, tag string) string {
	for {
		start := strings.Index(strings.ToLower(str), "<"+tag)
		if start == -1 {
			break
		}
		end := strings.Index(strings.ToLower(str[start:]), "</"+tag+">")
		if end == -1 {
			break
		}
		str = str[:start] + str[start+end+len(tag)+3:]
	}
	return str
}

func removeEventHandlers(str string) string {
	return strings.ReplaceAll(str, "on", "skip-on")
}

type ManifestItem struct {
	ZipFile
	ID        string
	Href      string
	RealPath  string
	MediaType string
}
