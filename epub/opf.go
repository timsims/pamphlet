package epub

import "encoding/xml"

// Opf is the struct for the opf file
type Opf struct {
	XMLName  xml.Name `xml:"package"`
	MetaData MetaData `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type MetaData struct {
	XMLName     xml.Name `xml:"metadata"`
	Title       string   `xml:"title"`
	Creator     Creator  `xml:"creator"`
	Language    string   `xml:"language"`
	Identifier  string   `xml:"identifier"`
	Publisher   string   `xml:"publisher"`
	Description string   `xml:"description"`
	Subject     string   `xml:"subject"`
	Date        string   `xml:"date"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID    string `xml:"id,attr"`
	Href  string `xml:"href,attr"`
	Media string `xml:"media-type,attr"`
}

type Spine struct {
	ItemRefs []ItemRef `xml:"itemref"`
	Toc      string    `xml:"toc,attr"`
}

type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}

type Creator struct {
	Data string `xml:",chardata"`
}
