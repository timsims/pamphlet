package pamphlet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	path = "./fixture"
)

func TestParseNormalEpub(t *testing.T) {
	parser, err := Open(fmt.Sprintf("%s/normal.epub", path))
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	testParseBookContent(t, parser.GetBook())
	assert.Nil(t, parser.Close())
}

func TestParseNormalEpubByBytes(t *testing.T) {
	filename := fmt.Sprintf("%s/normal.epub", path)
	file, err := os.Open(filename)
	assert.Nil(t, err)
	parser, err := OpenFile(file)
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	testParseBookContent(t, parser.GetBook())
	assert.Nil(t, parser.Close())
}

func TestParseWrongToCAttributeEpubFile(t *testing.T) {
	parser, err := Open(fmt.Sprintf("%s/wrong_toc_attribute.epub", path))
	assert.Nil(t, err)
	assert.NotNil(t, parser)

	testParseBookContent(t, parser.GetBook())
	assert.Nil(t, parser.Close())
}

func TestParseMissingEpubFile(t *testing.T) {
	_, err := Open(fmt.Sprintf("%s/not_existing.epub", path))
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrOpenEpub)
}

func TestParseInvalidEpub(t *testing.T) {
	_, err := Open(fmt.Sprintf("%s/wrong_mimetype.epub", path))
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrNotEpub)
}

func testParseBookContent(t *testing.T, book *Book) {
	assert.NotNil(t, book)
	assert.Equal(t, "normal book", book.Title)
	assert.Equal(t, "timsims", book.Author)

	chapters := book.Chapters
	assert.NotEmpty(t, chapters)
}
