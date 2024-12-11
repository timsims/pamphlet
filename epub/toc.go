package epub

import "encoding/xml"

// TocContainer is the struct for the toc.ncx file
type TocContainer struct {
	XMLName xml.Name `xml:"ncx"`
	NavMap  NavMap   `xml:"navMap"`
}

type NavMap struct {
	XMLName   xml.Name   `xml:"navMap"`
	NavPoints []NavPoint `xml:"navPoint"`
}

type NavPoint struct {
	XMLName   xml.Name `xml:"navPoint"`
	ID        string   `xml:"id,attr"`
	PlayOrder int      `xml:"playOrder,attr"`
	NavLabel  NavLabel `xml:"navLabel"`
	Content   Content  `xml:"content"`
	// NavPoints can be nested
	NavPoints []NavPoint `xml:"navPoint"`
}

type Content struct {
	Src string `xml:"src,attr"`
}

type NavLabel struct {
	XMLName xml.Name `xml:"navLabel"`
	Text    string   `xml:"text"`
}
