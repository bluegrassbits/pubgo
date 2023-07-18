package main

import (
	"io"
	"log"

	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown/ast"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	htmlFormatter = html.New(html.Standalone(false), html.TabWidth(2))
	if htmlFormatter == nil {
		log.Println("couldn't create html formatter")
	}
	styleName := "dracula"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		log.Printf("didn't find style '%s'", styleName)
	}
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}

// an actual rendering of Paragraph is more complicated
func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender(toc bool) (*mdhtml.Renderer, *parser.Parser) {
	var flags mdhtml.Flags

	if toc {
		flags = mdhtml.TOC
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags | flags,
		RenderNodeHook: myRenderHook,
	}
	return mdhtml.NewRenderer(opts), p
}
