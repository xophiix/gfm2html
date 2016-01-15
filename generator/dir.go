package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

// Dir represents file directory
type Dir struct {
	name        string
	dir         []*Dir
	mdDir       string
	htmlDir     string
	longDirName string
	pages       []*Page
	sidebar     template.HTML
	template    string
	static      []*StaticFile
}

// NewDir returns new dir
func NewDir(name, md, html, t, longDirName string) (*Dir, error) {
	fmt.Println("new dir", html, longDirName)

	return &Dir{
		name:        name,
		mdDir:       md,
		htmlDir:     html,
		template:    t,
		longDirName: longDirName,
	}, nil
}

// read reading all child directory and pages from dir
func (d *Dir) read() error {
	fmt.Printf("Read dir: %s\n", d.mdDir)
	osd, err := os.Open(d.mdDir)
	defer osd.Close()

	fi, err := osd.Readdir(-1)

	if err != nil {
		return err
	}

	for _, f := range fi {
		if f.Mode().IsDir() {
			dir, err := NewDir(f.Name(),
				getPath(d.mdDir, f.Name()),
				d.htmlDir,
				d.template,
				getUnderscorePath(d.longDirName, f.Name()),
			)
			if err == nil {
				dir.read()
				d.addDir(dir)
			}
		}
		if f.Mode().IsRegular() {
			page, err := d.NewPage(f)
			if err == nil {
				d.addPage(page)
			} else {
				st, err := d.NewStatic(f)
				if err == nil {
					d.addStatic(st)
				}
			}
		}
	}

	return nil
}

// write writes content to html directory
func (d *Dir) write(parent *Dir) error {
	err := os.MkdirAll(d.htmlDir, 0775)
	if err != nil {
		return err
	}
	sd, err := NewSidebar(getPath(d.mdDir, "_Sidebar.md"))

	if err == nil {
		fmt.Printf("Create new sidebar \n\t%s\n", d.mdDir)
		d.sidebar = sd
	} else {
		fmt.Printf("Sidebar not found \n\t%s\n", d.mdDir)
		d.sidebar = parent.sidebar
	}

	for _, p := range d.pages {
		err := p.save(d)
		if err != nil {
			return err
		}
	}

	for _, dir := range d.dir {
		err := dir.write(d)
		if err != nil {
			return err
		}
	}

	for _, st := range d.static {
		err := st.write(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func getUnderscorePath(parentName, name string) string {
	if parentName == "" {
		return name
	} else {
		return parentName + "_" + name
	}
}

// addPage adding new page to current dir
func (d *Dir) addPage(p *Page) {
	d.pages = append(d.pages, p)
}

// addDir adding new child directory
func (d *Dir) addDir(dir *Dir) {
	d.dir = append(d.dir, dir)
}

// addStating adding new static file to direcotory
func (d *Dir) addStatic(s *StaticFile) {
	d.static = append(d.static, s)
}

//  getPath returns concat string for current dir path
func getPath(c, f string) string {
	if c != "" {
		return fmt.Sprintf("%s%s%s", c, string(os.PathSeparator), f)
	} else {
		return f
	}
}

type TocData struct {
	Link     string
	Title    string
	Children []*TocData
}

func (d *Dir) GenerateToc(outputPath string) error {
	toc := &TocData{"toc", "toc", make([]*TocData, 0)}
	generateToc(d, toc)

	bytes, err := json.Marshal(toc)
	if err != nil {
		return err
	}

	tocStr := fmt.Sprintf("toc=%s", bytes)
	return ioutil.WriteFile(outputPath, []byte(tocStr), 0644)
}

func generateToc(d *Dir, toc *TocData) {
	for _, p := range d.pages {
		pageToc := &TocData{p.Url, p.Title, nil}
		toc.Children = append(toc.Children, pageToc)
	}

	for _, dir := range d.dir {
		dirToc := &TocData{"null", dir.name, make([]*TocData, 0)}
		toc.Children = append(toc.Children, dirToc)

		generateToc(dir, dirToc)
	}
}
