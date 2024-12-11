package pamphlet

import (
	"archive/zip"
	"io"
)

type ZipReaderCloser interface {
	Close() error
	Files() []*zip.File
}

type ZipReader struct {
	reader *zip.Reader
}

type ZipReadCloser struct {
	reader *zip.ReadCloser
}

func (z *ZipReader) Close() error {
	return nil
}

func (z *ZipReadCloser) Close() error {
	return z.reader.Close()
}

func (z *ZipReader) Files() []*zip.File {
	return z.reader.File
}

func (z *ZipReadCloser) Files() []*zip.File {
	return z.reader.File
}

type ZipFile struct {
	file *zip.File
}

// Open opens the file in zip for reading
func (z *ZipFile) Open() (io.ReadCloser, error) {
	return z.file.Open()
}

// GetRawContent returns the raw content of the file
func (c *ZipFile) GetRawContent() ([]byte, error) {

	rc, err := c.Open()
	if err != nil {
		return nil, err
	}

	defer func(rc io.ReadCloser) {
		err = rc.Close()
		if err != nil {

		}
	}(rc)

	content, err := io.ReadAll(rc)

	if err != nil {
		return nil, err
	}

	return content, nil
}
