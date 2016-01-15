package generator

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

// GenerateDoc generating new documentation
func GenerateDoc(c *cli.Context) {
	md := c.String("input")
	html := c.String("output")
	t := c.String("template")
	sidebar := c.String("sidebar")
	tocOutput := c.String("tocOutput")

	if md == "" {
		cli.ShowAppHelp(c)
		return
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
