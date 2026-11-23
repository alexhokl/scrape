package scraper

import (
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func convertTestHTML(t *testing.T, body string, skip ...string) string {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><div id="root">` + body + `</div></body></html>`))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	base, err := url.Parse("https://example.com/blog/post/")
	if err != nil {
		t.Fatalf("failed to parse base URL: %v", err)
	}
	return htmlToMarkdown(doc.Find("#root"), base, skip)
}

func TestHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "paragraphs",
			input: `<p>First   paragraph
with a line break in source.</p><p>Second.</p>`,
			expected: "First paragraph with a line break in source.\n\nSecond.\n",
		},
		{
			name:     "headings",
			input:    `<h1>One</h1><h2>Two</h2><h3>Three</h3><h4>Four</h4><h5>Five</h5><h6>Six</h6>`,
			expected: "# One\n\n## Two\n\n### Three\n\n#### Four\n\n##### Five\n\n###### Six\n",
		},
		{
			name:     "bold, italic and strikethrough",
			input:    `<p><strong>bold</strong> <b>b</b> <em>italic</em> <i>i</i> <del>gone</del></p>`,
			expected: "**bold** **b** *italic* *i* ~~gone~~\n",
		},
		{
			name:     "emphasis keeps surrounding spaces outside markers",
			input:    `<p>Hello<em> world </em>again</p>`,
			expected: "Hello *world* again\n",
		},
		{
			name:     "absolute link",
			input:    `<p>See <a href="https://go.dev/">Go</a>.</p>`,
			expected: "See [Go](https://go.dev/).\n",
		},
		{
			name:     "relative link resolved against base URL",
			input:    `<p><a href="/about">About</a> and <a href="other/">other</a></p>`,
			expected: "[About](https://example.com/about) and [other](https://example.com/blog/post/other/)\n",
		},
		{
			name:     "fragment link rendered as text",
			input:    `<p>Note<sup><a href="#fn:1">1</a></sup></p>`,
			expected: "Note1\n",
		},
		{
			name:     "link without text is dropped",
			input:    `<p>Text<a href="/x"> </a></p>`,
			expected: "Text\n",
		},
		{
			name:     "inline code",
			input:    `<p>Run <code>go test</code> and <code>a ` + "`" + `b</code></p>`,
			expected: "Run `go test` and ``a `b``\n",
		},
		{
			name:     "line break",
			input:    `<p>one<br>two</p>`,
			expected: "one\ntwo\n",
		},
		{
			name:     "unordered list",
			input:    `<ul><li>One</li><li><a href="/two">Two</a></li></ul>`,
			expected: "* One\n* [Two](https://example.com/two)\n",
		},
		{
			name:     "ordered list",
			input:    `<ol><li>First</li><li>Second</li></ol>`,
			expected: "1. First\n2. Second\n",
		},
		{
			name:     "nested list",
			input:    `<ul><li>Parent<ul><li>Child</li><li>Child 2</li></ul></li><li>Sibling</li></ul>`,
			expected: "* Parent\n  * Child\n  * Child 2\n* Sibling\n",
		},
		{
			name:     "list item with paragraphs",
			input:    `<ol><li><p>First line</p><p>Second line</p></li></ol>`,
			expected: "1. First line\n   Second line\n",
		},
		{
			name:     "code block with language class on code",
			input:    "<pre><code class=\"language-go\">func main() {\n\tfmt.Println(\"&lt;hi&gt;\")\n}\n</code></pre>",
			expected: "```go\nfunc main() {\n\tfmt.Println(\"<hi>\")\n}\n```\n",
		},
		{
			name:     "code block with lang class on pre",
			input:    "<pre class=\"lang-csharp\">var x = 1;</pre>",
			expected: "```csharp\nvar x = 1;\n```\n",
		},
		{
			name:     "code block without language",
			input:    "<pre><code>plain</code></pre>",
			expected: "```\nplain\n```\n",
		},
		{
			name:     "code block containing a fence",
			input:    "<pre><code>```\ninner\n```</code></pre>",
			expected: "````\n```\ninner\n```\n````\n",
		},
		{
			name:     "code block with highlighted spans",
			input:    `<pre class="language-sh"><code><span>echo</span> <span>hi</span></code></pre>`,
			expected: "```sh\necho hi\n```\n",
		},
		{
			name:     "blockquote",
			input:    `<blockquote><p>Quoted <em>text</em></p><p>Second</p></blockquote>`,
			expected: "> Quoted *text*\n>\n> Second\n",
		},
		{
			name:     "figure with caption",
			input:    `<figure><img src="/img/a.png" alt="An image"><figcaption>The caption</figcaption></figure>`,
			expected: "![An image](https://example.com/img/a.png)\n\n*The caption*\n",
		},
		{
			name:     "figure without image is rendered as blocks",
			input:    `<figure><pre><code>code</code></pre></figure>`,
			expected: "```\ncode\n```\n",
		},
		{
			name:     "lazy-loaded image uses data-src",
			input:    `<p><img src="data:image/gif;base64,AAAA" data-src="https://cdn.example.com/a.jpg" alt="Lazy"></p>`,
			expected: "![Lazy](https://cdn.example.com/a.jpg)\n",
		},
		{
			name:     "protocol-relative image",
			input:    `<div class="img"><img src="//cdn.example.com/b.png" alt="fig"></div>`,
			expected: "![fig](https://cdn.example.com/b.png)\n",
		},
		{
			name:     "image with only a data URI is dropped",
			input:    `<p><img src="data:image/gif;base64,AAAA"></p>`,
			expected: "",
		},
		{
			name:     "table",
			input:    `<table><thead><tr><th>Name</th><th>Value</th></tr></thead><tbody><tr><td>a|b</td><td><code>1</code></td></tr><tr><td>c</td></tr></tbody></table>`,
			expected: "| Name | Value |\n| --- | --- |\n| a\\|b | `1` |\n| c |  |\n",
		},
		{
			name:     "horizontal rule",
			input:    `<p>Above</p><hr><p>Below</p>`,
			expected: "Above\n\n---\n\nBelow\n",
		},
		{
			name:     "wrapper divs are flattened",
			input:    `<div class="paragraph"><p>One</p></div><section><div><p>Two</p></div></section>`,
			expected: "One\n\nTwo\n",
		},
		{
			name:     "loose inline content becomes a paragraph",
			input:    `Loose <b>text</b><p>Para</p>tail`,
			expected: "Loose **text**\n\nPara\n\ntail\n",
		},
		{
			name:     "inline element wrapping blocks keeps the blocks",
			input:    `<span><p>One</p><p>Two</p></span>`,
			expected: "One\n\nTwo\n",
		},
		{
			name:     "scripts, styles, iframes and forms are removed",
			input:    `<p>Kept</p><script>alert(1)</script><style>p{}</style><iframe src="x"></iframe><form><input></form><noscript>No JS</noscript>`,
			expected: "Kept\n",
		},
		{
			name:     "definition list",
			input:    `<dl><dt>Term</dt><dd>Definition</dd></dl>`,
			expected: "**Term**\n\nDefinition\n",
		},
		{
			name:     "empty input",
			input:    ``,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTestHTML(t, tt.input)
			if result != tt.expected {
				t.Errorf("htmlToMarkdown(%q)\n got: %q\nwant: %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHTMLToMarkdown_SkipSelectors(t *testing.T) {
	input := `<h1>Title</h1><p class="meta">Date</p><p>Body <span class="sr-only">(new window)</span>text</p><div class="share"><p>Share</p></div>`
	result := convertTestHTML(t, input, "h1", "h1 + p", "span.sr-only", "div.share")

	if result != "Body text\n" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestHTMLToMarkdown_DoesNotModifySource(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<div id="root"><p class="x">Gone</p><p>Kept</p></div>`))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	root := doc.Find("#root")
	_ = htmlToMarkdown(root, nil, []string{"p.x"})

	if root.Find("p.x").Length() != 1 {
		t.Error("expected the source document to keep skipped elements")
	}
}

func TestHTMLToMarkdown_NilBaseKeepsRelativeURLs(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<div id="root"><p><a href="/x">X</a></p></div>`))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	result := htmlToMarkdown(doc.Find("#root"), nil, nil)
	if result != "[X](/x)\n" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestCodeLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"language-go", "go"},
		{"hljs language-csharp", "csharp"},
		{"lang-c++", "c++"},
		{"prettyprint", ""},
		{"", ""},
		{"notlanguage-go", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if result := codeLanguage(tt.input); result != tt.expected {
				t.Errorf("codeLanguage(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWrapInline(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		marker   string
		expected string
	}{
		{"plain", "word", "**", "**word**"},
		{"leading space", " word", "*", " *word*"},
		{"trailing space", "word ", "*", "*word* "},
		{"only spaces", "  ", "*", "  "},
		{"empty", "", "*", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := wrapInline(tt.text, tt.marker); result != tt.expected {
				t.Errorf("wrapInline(%q, %q) = %q, want %q", tt.text, tt.marker, result, tt.expected)
			}
		})
	}
}
