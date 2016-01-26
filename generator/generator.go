package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/codegangsta/cli"
)

type GenerateOptions struct {
	FolderAliases map[string]string
}

var generateOptions *GenerateOptions

func (opts *GenerateOptions) parseGenerateOptions(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, opts)
}

// GenerateDoc generating new documentation
func GenerateDoc(c *cli.Context) {
	md := c.String("input")
	html := c.String("output")
	t := c.String("template")
	sidebar := c.String("sidebar")
	tocOutput := c.String("tocOutput")
	optionFile := c.String("optionFile")

	if md == "" {
		cli.ShowAppHelp(c)
		return
	}

	if optionFile != "" {
		generateOptions = &GenerateOptions{make(map[string]string)}
		err := generateOptions.parseGenerateOptions(optionFile)
		if err != nil {
			fmt.Println("Option file provided but parse failed: ", optionFile, err)
		}
	}

	fmt.Println("Begin generate")
	sb, err := NewSidebar(sidebar)
	if err != nil {
		fmt.Println("Sidebar not exists and will autogenerate from dir hierachy")
	}

	parent := &Dir{sidebar: sb}

	splits := strings.Split(md, "/")
	dir, err := NewDir(splits[len(splits)-1], md, html, t, "")
	if err != nil {
		fmt.Printf("Error read dir %s\n \t%s\n", dir.mdDir, err.Error())
	}
	err = dir.read()
	if err != nil {
		fmt.Printf("Error read dir %s\n \t%s\n", dir.mdDir, err.Error())
	}

	if tocOutput != "" {
		err = dir.GenerateToc(tocOutput)
		if err != nil {
			fmt.Printf("Error generate toc file %s: %s\n", tocOutput, err.Error())
		}
	}

	err = dir.write(parent)

	if err != nil {
		fmt.Printf("Error write dir %s\n", dir.htmlDir)
	}

	fmt.Println("End generate")
}
