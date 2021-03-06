package main

import (
	"os"

	"./generator"
	"github.com/codegangsta/cli"
)

const APP_VER = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "gfm2html"
	app.Email = "xophiix@gmail.com"
	app.Usage = "Github generator html pages from markdown wiki"
	app.Version = APP_VER
	app.Action = generator.GenerateDoc
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input, i",
			Usage: "Directory with markdown files",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "documentation",
			Usage: "Directory for output files",
		},
		cli.StringFlag{
			Name:  "template, t",
			Value: "templates/documentation.tpl",
			Usage: "Template for generated documentation",
		},
		cli.StringFlag{
			Name:  "tocOutput, f",
			Value: "_toc.js",
			Usage: "OutputPath to toc data file",
		},
		cli.StringFlag{
			Name:  "optionFile, p",
			Usage: "Generate option config json file",
		},
	}
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "V",
			Email: "support@leanlabs.io",
		},
		cli.Author{
			Name:  "cnam",
			Email: "cnam812@gmail.com",
		},
		cli.Author{
			Name:  "xophiix",
			Email: "xophiix@gmail.com",
		},
	}
	app.Run(os.Args)
}
