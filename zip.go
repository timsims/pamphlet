package pamphlet

import (
	"archive/zip"
	"io"
)

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
