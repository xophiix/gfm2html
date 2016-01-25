package generator

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"../parser"
)

//Page represent page for generate
type Page struct {
	Title    string
	FileName string
	Url      string
	Path     string
	Items    []*Page
	Body     template.HTML
	Sidebar  template.HTML
	Template string
	Seo      *Seo
}

type Seo struct {
	Title       string
	Description string
	Keywords    string
}

func getTitleByFileName(name string) string {
	return strings.Replace(strings.Replace(name, ".md", "", 1), "_", " ", -1)
}

// NewPage create new page
func (d *Dir) NewPage(f os.FileInfo) (*Page, error) {
	prs, err := parser.New(f.Name())
	if err != nil || strings.HasPrefix(f.Name(), "_") {
		return nil, errors.New(fmt.Sprintf("Not allowed file format %s\n", f.Name()))
	}

	cont, err := ioutil.ReadFile(getPath(d.mdDir, f.Name()))

	if err != nil {
		return nil, err
	}

	html, title := prs.Parse(cont)
	if title == "" {
		title = getTitleByFileName(f.Name())
	}

	p := &Page{}
	ext := path.Ext(f.Name())
	p.FileName = strings.Replace(f.Name(), ext, "", -1)
	p.Title = title

	fmt.Println(f.Name(), "'s title=", title)

	p.Seo = &Seo{
		Title:       "",
		Description: "",
		Keywords:    "",
	}
	p.Body = template.HTML(html)
	p.Path = getPagePath(p, d, false)
	p.Url = getPagePath(p, d, true)
	p.Template = d.template

	fmt.Printf("new page: %s, %s, %s\n", p.Title, p.FileName, p.Path)

	return p, nil
}

func getPagePath(page *Page, d *Dir, url bool) string {
	filename := page.FileName + ".html"
	prefix := ""

	if !url {
		prefix = d.htmlDir
		if prefix != "" {
			prefix += "/"
		}
	}

	if d.longDirName == "" {
		return prefix + filename
	} else {
		return prefix + d.longDirName + "_" + filename
	}
}

// save saving current page to filesystem
func (p *Page) save(d *Dir) error {
	p.Sidebar = d.sidebar
	p.Items = d.pages
	file, err := os.Create(p.Path)

	if err != nil {
		return err
	}

	fmt.Printf("Create new page: %s\n \tby link:%s\n", p.Title, p.Path)

	return p.render(file)
}

// render rendering current page template
func (p *Page) render(f *os.File) error {
	t, err := template.ParseFiles(p.Template)

	if err != nil {
		return err
	}

	return t.Execute(f, p)
}
