package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/sanitized_anchor_name"
	"golang.org/x/net/html"
)

type MdParser struct {
}

type renderer struct {
	*blackfriday.Html
}

func NewMdParser() *MdParser {
	return &MdParser{}
}

// removeStuf trims spaces, removes new lines and code tag from a string
func removeStuf(s string) string {
	res := strings.Replace(s, "\n", "", -1)
	res = strings.Replace(res, "<code>", "", -1)
	res = strings.Replace(res, "</code>", "", -1)
	res = strings.TrimSpace(res)

	return res
}

func (prs *MdParser) Parse(d []byte) (string, string) {
	renderer := &renderer{Html: blackfriday.HtmlRenderer(0, "", "").(*blackfriday.Html)}

	// Parser extensions for GitHub Flavored Markdown.
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	//extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK

	unsanitized := blackfriday.Markdown(d, renderer, extensions)

	// GitHub Flavored Markdown-like sanitization policy.
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span", "code")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("checked", "disabled").Matching(regexp.MustCompile(`^$`)).OnElements("input")
	p.AllowDataURIImages()

	html := string(p.SanitizeBytes(unsanitized))
	html = strings.Replace(html, "README.md", "index.html", -1)
	html = strings.Replace(html, ".md", ".html", -1)

	// extract title from 'h1' tag
	re := `<h1>\s*<a.*>\s*.*\s*</a>\s*(?P<name>.*?)\s*</h1>`
	r := regexp.MustCompile(re)
	title := ""

	groups := make(map[string]string)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		// fill map for groups
		for i, name := range r.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			fmt.Println("match:", i, match[i])
			groups[name] = removeStuf(match[i])
		}

		title = removeStuf(groups["name"])
		break
	}

	// auto generate toc
	tocMd := grabToc(html)
	unsanitized = blackfriday.Markdown([]byte(tocMd), renderer, extensions)
	tocHtml := string(p.SanitizeBytes(unsanitized))

	if strings.Index(html, "[TOC]") >= 0 {
		html = strings.Replace(html, "[TOC]", tocHtml, 1)
	} else {
		html = tocHtml + html
	}

	fmt.Println("tocMd=", tocMd)

	return html, title
}

// Escapes special characters
func escapeSpecChars(s string) string {
	specChar := []string{"\\", "`", "*", "_", "{", "}", "#", "+", "-", ".", "!"}
	res := s

	for _, c := range specChar {
		res = strings.Replace(res, c, "\\"+c, -1)
	}
	return res
}

func grabToc(html string) string {
	re := `(?si)<h(?P<num>[1-6])>\s*` +
		`<a\s*name="[^"]*"\s*class="anchor"\s*` +
		`href="(?P<href>[^"]*)"[^>]*>\s*` +
		`<span[^<*]*</span>\s*</a>\s*(?P<name>.*?)</h`

	r := regexp.MustCompile(re)

	toc := ""
	groups := make(map[string]string)
	for _, match := range r.FindAllStringSubmatch(html, -1) {
		// fill map for groups
		for i, name := range r.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			groups[name] = removeStuf(match[i])
		}

		// format result
		n, _ := strconv.Atoi(groups["num"])
		if n > 1 {
			link := groups["href"]
			toc_item := strings.Repeat("  ", n-1) + "* " +
				"[" + escapeSpecChars(removeStuf(groups["name"])) + "]" +
				"(" + link + ")"

			toc = toc + toc_item + "\n"
		}
	}

	return toc
}

// GitHub Flavored Markdown header with clickable and hidden anchor.
func (_ *renderer) Header(out *bytes.Buffer, text func() bool, level int, _ string) {
	marker := out.Len()
	doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}

	textHtml := out.String()[marker:]
	out.Truncate(marker)

	// Extract text content of the header.
	var textContent string
	if node, err := html.Parse(strings.NewReader(textHtml)); err == nil {
		textContent = extractText(node)
	} else {
		// Failed to parse HTML (probably can never happen), so just use the whole thing.
		textContent = html.UnescapeString(textHtml)
	}
	anchorName := sanitized_anchor_name.Create(textContent)

	out.WriteString(fmt.Sprintf(`<h%d><a name="%s" class="anchor" href="#%s" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>`, level, anchorName, anchorName))
	out.WriteString(textHtml)
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

// extractText returns the recursive concatenation of the text content of an html node.
func extractText(n *html.Node) string {
	var out string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			out += c.Data
		} else {
			out += extractText(c)
		}
	}
	return out
}

// TODO: Clean up and improve this code.
// GitHub Flavored Markdown fenced code block with highlighting.
func (_ *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	doubleSpace(out)

	// parse out the language name
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		out.WriteString(`<pre><code class="`)
		attrEscape(out, []byte(elt))
		lang = elt
		out.WriteString(`">`)
		count++
		break
	}

	if count == 0 {
		out.WriteString("<pre><code>")
	}

	attrEscape(out, text)

	if count == 0 {
		out.WriteString("</code></pre>\n")
	} else {
		out.WriteString("</code></pre>\n")
	}
}

// Task List support.
func (r *renderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	switch {
	case bytes.HasPrefix(text, []byte("[ ] ")):
		text = append([]byte(`<input type="checkbox" disabled="">`), text[3:]...)
	case bytes.HasPrefix(text, []byte("[x] ")) || bytes.HasPrefix(text, []byte("[X] ")):
		text = append([]byte(`<input type="checkbox" checked="" disabled="">`), text[3:]...)
	}
	r.Html.ListItem(out, text, flags)
}

// Unexported blackfriday helpers.

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}
