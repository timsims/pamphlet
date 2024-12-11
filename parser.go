package pamphlet

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"pamphlet/epub"
	"strings"
)

const (
	// Mimetype of epub file
	Mimetype      = "application/epub+zip"
	NcxMediaType  = "application/x-dtbncx+xml"
	ContainerFile = "meta-inf/container.xml"
	NcxFileExt    = ".ncx"
)

type Parser struct {
	// current epub file pointer
	zipReader *zip.ReadCloser

	// root file of epub
	rootFile *epub.RootFile

	// directory of manifest files, usually same as the root file
	manifestPath string

	// package struct of epub, parse from xml of root file
	opf *epub.Opf

	// map of manifest files, key is the id attribute of the item
	manifestFiles map[string]ManifestFile

	// map of navPoint of chapters, it contains chapter's title, order and src
	//key is the src attribute of the navPoint
	toc map[string]epub.NavPoint

	book *Book
}

type ManifestFile struct {
	ID        string
	mediaType string
	Href      string
	file      *zip.File
}

func NewParser(path string) (*Parser, error) {
	result, err := zip.OpenReader(path)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("failed to open epub file")
	}
	parser := &Parser{
		zipReader:     result,
		toc:           make(map[string]epub.NavPoint),
		manifestFiles: make(map[string]ManifestFile),
	}

	err = parser.parse()

	if err != nil {
		return nil, err
	}

	return parser, nil
}

func (p *Parser) parse() error {
	if !p.checkMimeType() {
		return errors.New("not epub")
	}

	err := p.getRootFile()
	if err != nil {
		return err
	}

	err = p.parseRootFile()
	if err != nil {
		return err
	}

	p.cacheManifestFiles()
	p.parseToc()
	p.convertToBook()
	return nil
}

func (p *Parser) Close() error {
	return p.zipReader.Close()
}

// CheckMimeType check mimetype file of epub file
// check file mimetype exists
// check content of mimetype file is "application/epub+zip"
// return true if file is epub
func (p *Parser) checkMimeType() bool {
	isValid := false

	for _, f := range p.zipReader.File {
		if f.Name == "mimetype" {
			rc, err := f.Open()
			if err != nil {
				return isValid
			}
			mimetype, err := io.ReadAll(rc)
			if err != nil {
				return isValid
			}
			if strings.TrimSpace(string(mimetype)) == Mimetype {
				isValid = true
			}
			err = rc.Close()
			if err != nil {
				return isValid
			}
		}
	}
	return isValid
}

// getRootFile get root file of epub
// root file is the file that contains the metadata of the epub,
// META-INF/container.xml is the file that defines the root file's name
func (p *Parser) getRootFile() error {

	for _, f := range p.zipReader.File {
		if strings.ToLower(f.Name) == ContainerFile {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			var container epub.Container
			err = xml.Unmarshal(content, &container)
			if err != nil {
				return err
			}
			err = rc.Close()
			if err != nil {
				return err
			}

			if container.RootFile.FullPath == "" {
				return errors.New("root file not found")
			}

			p.rootFile = &container.RootFile
			return nil
		}
	}
	return errors.New("root file not found")
}

// parseRootFile parse root file of epub
// root file is the file that contains the metadata of the epub
// we parse xml and get the struct of Opf
// and also determine the path of manifest files
func (p *Parser) parseRootFile() error {
	for _, f := range p.zipReader.File {
		if f.Name == p.rootFile.FullPath {

			// get the manifest files path
			path := strings.Split(p.rootFile.FullPath, "/")[0]
			// in case the manifest files is not in the root directory
			if path != p.rootFile.FullPath {
				p.manifestPath = path
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			var pack epub.Opf
			err = xml.Unmarshal(content, &pack)
			if err != nil {
				return err
			}
			err = rc.Close()
			if err != nil {
				return err
			}
			p.opf = &pack
			return nil
		}
	}
	return errors.New("opf file not found")
}

// parseToc parse table of content of epub
// we can get the tile of each chapters from the toc file
// toc.ncx is the default toc file name
func (p *Parser) parseToc() {
	tocFile := p.getTocFile()
	if tocFile != nil {
		rc, err := tocFile.Open()
		if err != nil {
			return
		}
		content, err := io.ReadAll(rc)
		if err != nil {
			return
		}

		// parse toc file
		var container epub.TocContainer
		err = xml.Unmarshal(content, &container)
		if err != nil {
			log.Fatal(err)
		}
		err = rc.Close()
		if err != nil {
			log.Fatal(err)
		}

		// get the title of each chapter
		for _, navPoint := range container.NavMap.NavPoints {
			p.parseNavPoint(navPoint)
		}
	}
}

// getTocFile get toc file of epub
// due to the toc file can be named differently, we need to guess the toc file
func (p *Parser) getTocFile() *zip.File {
	var possibleTocFiles []*zip.File

	// guess toc file by media type
	for _, f := range p.manifestFiles {
		if f.mediaType == NcxMediaType {
			return f.file
		}

		// guess toc file by file extension
		if strings.HasSuffix(f.Href, NcxFileExt) {
			possibleTocFiles = append(possibleTocFiles, f.file)
		}
	}

	if len(possibleTocFiles) > 0 {
		return possibleTocFiles[0]
	}

	return nil
}

func (p *Parser) parseNavPoint(navPoint epub.NavPoint) {
	src := navPoint.Content.Src
	// some src will have # after the file name we need to remove it
	// eg: src="chapter1.html#chapter1" we need to remove the #chapter1
	if strings.Contains(src, "#") {
		src = strings.Split(src, "#")[0]
	}

	p.toc[src] = navPoint

	// parse nested nav points
	for _, np := range navPoint.NavPoints {
		p.parseNavPoint(np)
	}
}

// cacheManifestFiles associate the manifest files with the id attribute and zip.File, so we can easily access the file later
func (p *Parser) cacheManifestFiles() {
	p.manifestFiles = make(map[string]ManifestFile, len(p.zipReader.File))

	for _, item := range p.opf.Manifest.Items {
		for _, f := range p.zipReader.File {
			fileName := item.Href

			// if manifestPath is not empty, then the file is in a subdirectory
			if p.manifestPath != "" {
				fileName = fmt.Sprintf("%s/%s", p.manifestPath, item.Href)
			}

			if f.Name == fileName {
				p.manifestFiles[item.ID] = ManifestFile{
					ID:        item.ID,
					mediaType: item.Media,
					Href:      item.Href,
					file:      f,
				}
			}
		}
	}
}

// convertToBook convert the epub.Opf to Book struct
func (p *Parser) convertToBook() {
	meta := p.opf.MetaData
	p.book = &Book{
		meta.Title,
		meta.Creator.Data,
		meta.Language,
		meta.Identifier,
		meta.Publisher,
		meta.Description,
		meta.Subject,
		meta.Date,
		p.convertChapters(),
		p.convertManifest(),
	}
}

// convertChapters convert the epub.Spine to Chapter struct
func (p *Parser) convertChapters() []Chapter {
	chapters := make([]Chapter, 0)
	for _, item := range p.opf.Spine.ItemRefs {
		if f, ok := p.manifestFiles[item.IDRef]; ok {

			// not all chapters have reference in toc file,
			//so we need to check if the chapter has title
			chapterTitle := ""
			order := 0
			hasToc := false
			if navPoint, ok := p.toc[f.Href]; ok {
				chapterTitle = navPoint.NavLabel.Text
				order = navPoint.PlayOrder
				hasToc = true
			}

			chapter := Chapter{
				ID:        item.IDRef,
				Href:      f.Href,
				MediaType: f.mediaType,
				Title:     chapterTitle,
				HasToc:    hasToc,
				ZipFile:   ZipFile{f.file},
				Order:     order,
			}
			chapters = append(chapters, chapter)
		}
	}
	return chapters
}

// convertManifest convert the epub.Manifest  to ManifestItem struct
func (p *Parser) convertManifest() []ManifestItem {
	manifestItems := make([]ManifestItem, 0)
	for _, item := range p.opf.Manifest.Items {
		if f, ok := p.manifestFiles[item.ID]; ok {
			manifestItem := ManifestItem{
				ZipFile:   ZipFile{f.file},
				ID:        item.ID,
				Href:      item.Href,
				RealPath:  f.file.Name,
				MediaType: item.Media,
			}

			manifestItems = append(manifestItems, manifestItem)
		}
	}
	return manifestItems
}

func (p *Parser) Book() *Book {
	return p.book
}
