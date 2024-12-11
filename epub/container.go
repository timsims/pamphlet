package epub

import (
	"encoding/xml"
)

// Container is the struct for the container.xml file
type Container struct {
	XMLName  xml.Name `xml:"container"`
	RootFile RootFile `xml:"rootfiles>rootfile"`
}

// RootFile is the struct for the root file
type RootFile struct {
	FullPath string `xml:"full-path,attr"`
	Media    string `xml:"media-type,attr"`
}
