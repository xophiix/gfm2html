# Github flavored Markdown to html page generator

### Overview

This repo is just an personal customization based on [md2html](https://github.com/cnam/md2html)
I want to adapt this tool to fit my project's requirement, and finally It does almost same thing md2html does except:
* all output html files in same folder level to make sure dependent resource path stay consistent(e.g. css or js path in html template)
* automatically generate sidebar from markdown folder hierachy(require custom js function) instead of read from external md file, subfolder display name of hierachy can be replace with alias in generation config.(e.g. subfolder name is `A` and you want to display as `Artist`)
* automatically generate table of contents prepended on each page

### Installation

build with source
```
git clone https://github.com/xophiix/gfm2html
cd gfm2html
go build
./gfm2html
```

### Usage

`gfm2html -i markdown -o html -t res/template.tpl -f res/js/toc.js -p res/generate_option.json`

**WHERE:**

- **-i or --input** Directory with markdown files
- **-o or --output** Directory for output generated html files
- **-t or --template** Template for generated documentation
- **-f or --tocFile** (Optional)output toc file which containing markdown folder hierachy info in js format, used to generate sidebar on page load
- **-p or --optionFile** (Optional)json file containing generation options(folder alias, etc.)

`-p` and `-s` option is removed from origin version.

You can enter sample folder and execute build.bat/build.sh for a glance of what's done, provided with a html template and a generation option file example.

### Templates

Your must be create html template with variables. [example](https://github.com/xophiix/gfm2html/blob/master/sample/res/documentation.tpl)

**Available Variables**

- **{{ .Title }}** Page title similarly filename
- **{{ .Body}}** Generated html body from markdown
