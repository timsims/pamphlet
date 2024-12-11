package pamphlet

import (
	"archive/zip"
	"bytes"
	"github.com/timsims/pamphlet/epub"
	"io"
	"os"
)

// NewParser
// Deprecated: use Open instead
func NewParser(filePath string) (*Parser, error) {
	return Open(filePath)
}

// Open open epub by file path
func Open(path string) (*Parser, error) {
	result, err := zip.OpenReader(path)
	if err != nil {
		return nil, ErrOpenEpub
	}
	parser := &Parser{
		zipReader:     &ZipReadCloser{result},
		toc:           make(map[string]epub.NavPoint),
		manifestFiles: make(map[string]ManifestFile),
	}

	err = parser.parse()

	if err != nil {
		return nil, err
	}

	return parser, nil
}

// OpenFile open epub by file pointer
// use this function if you have the epub file in memory, eg: from http request
func OpenFile(file *os.File) (*Parser, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrOpenEpub
	}
	return OpenBytes(fileBytes)
}

// OpenBytes open epub by byte slice
// use this function if you have the epub file in memory and you don't have file pointer
func OpenBytes(data []byte) (*Parser, error) {
	result, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {
		return nil, ErrOpenEpub
	}
	parser := &Parser{
		zipReader:     &ZipReader{result},
		toc:           make(map[string]epub.NavPoint),
		manifestFiles: make(map[string]ManifestFile),
	}

	err = parser.parse()

	if err != nil {
		return nil, err
	}

	return parser, nil
}
