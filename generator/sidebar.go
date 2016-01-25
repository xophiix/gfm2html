package generator

import (
	"html/template"
	"io/ioutil"

	"../parser"
)

// NewSidebar created new sidebar
func NewSidebar(dir string) (template.HTML, error) {
	s, err := generateSidebar(dir)

	if err != nil {
		return "", err
	}

	return s, nil
}

func generateSidebar(mdSidebar string) (template.HTML, error) {
	var sidebar template.HTML
	prs, err := parser.New(mdSidebar)
	if err != nil {
		return "", err
	}

	file, err := ioutil.ReadFile(mdSidebar)

	if err != nil {
		return "", err
	}

	html, _ := prs.Parse(file)
	sidebar = template.HTML(html)

	return sidebar, nil
}
