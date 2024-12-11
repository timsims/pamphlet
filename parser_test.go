package pamphlet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	path = "./fixture"
)

func TestParseNormalEpub(t *testing.T) {
	parser, err := NewParser(fmt.Sprintf("%s/normal.epub", path))
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	testParseBookContent(t, parser.GetBook())
}

func TestParseWrongToCAttributeEpubFile(t *testing.T) {
	parser, err := NewParser(fmt.Sprintf("%s/wrong_toc_attribute.epub", path))
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	testParseBookContent(t, parser.GetBook())
}

func TestParseMissingEpubFile(t *testing.T) {
	_, err := NewParser(fmt.Sprintf("%s/not_existing.epub", path))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), FailToOpenFile)
}

func TestParseInvalidEpub(t *testing.T) {
	_, err := NewParser(fmt.Sprintf("%s/wrong_mimetype.epub", path))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), NotEpub)
}

func testParseBookContent(t *testing.T, book *Book) {
	assert.NotNil(t, book)
	assert.Equal(t, "normal book", book.Title)
	assert.Equal(t, "timsims", book.Author)

	chapters := book.Chapters
	assert.NotEmpty(t, chapters)
}
