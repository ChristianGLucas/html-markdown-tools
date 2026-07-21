package nodes

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	readability "codeberg.org/readeck/go-readability/v2"

	"christiangeorgelucas/html-markdown-tools/axiom"
	gen "christiangeorgelucas/html-markdown-tools/gen"
)

// Isolate the main article content of an HTML page — stripping navigation,
// ads, sidebars, and other boilerplate — using a Go port of Mozilla's
// Readability (Firefox Reader Mode) algorithm, then render the isolated
// content to Markdown through the same engine and options as
// ConvertToMarkdown. Also returns the page metadata Readability recovers
// along the way (title, byline, excerpt, site name, lead image, language,
// published time) and the isolated content as plain text. base_url, if
// given, is used only to resolve the extracted content's relative
// links/images to absolute URLs — a pure string join, never a network
// fetch. extracted=false (with markdown/text_content empty) is the correct,
// non-error result only for a page with essentially no body content at all
// (an empty <body>); a low-signal page (a bare login form) still typically
// extracts something and reports extracted=true with a small `length` —
// check `length`, not just `extracted`, to judge substance.
func ExtractMainContentAsMarkdown(ctx context.Context, ax axiom.Context, input *gen.MainContentQuery) (*gen.MainContentResult, error) {
	if input == nil {
		return &gen.MainContentResult{Error: "request is required"}, nil
	}
	if err := checkHTMLSize(input.Html); err != nil {
		return &gen.MainContentResult{Error: err.Error()}, nil
	}

	var pageURL *url.URL
	if input.BaseUrl != "" {
		u, err := url.Parse(input.BaseUrl)
		if err != nil {
			return &gen.MainContentResult{Error: "base_url: " + err.Error()}, nil
		}
		pageURL = u
	}

	conv, err := buildConverter(input.Options)
	if err != nil {
		return &gen.MainContentResult{Error: err.Error()}, nil
	}

	ps := readability.NewParser()
	ps.MaxElemsToParse = maxReadabilityElems

	article, err := ps.Parse(strings.NewReader(input.Html), pageURL)
	if err != nil {
		return &gen.MainContentResult{Error: "extract: " + err.Error()}, nil
	}

	result := &gen.MainContentResult{
		Title:    article.Title(),
		Byline:   article.Byline(),
		Excerpt:  article.Excerpt(),
		SiteName: article.SiteName(),
		ImageUrl: article.ImageURL(),
		Language: article.Language(),
	}
	if pt, err := article.PublishedTime(); err == nil {
		result.PublishedTime = pt.Format(time.RFC3339)
	}

	if article.Node == nil {
		// Readability found no article-shaped content. Not an error — some
		// pages genuinely have none.
		return result, nil
	}
	result.Extracted = true

	var textBuf bytes.Buffer
	if err := article.RenderText(&textBuf); err != nil {
		result.Error = "render text: " + err.Error()
		return result, nil
	}
	result.TextContent = textBuf.String()
	result.Length = int32(utf8.RuneCountInString(result.TextContent))

	var htmlBuf bytes.Buffer
	if err := article.RenderHTML(&htmlBuf); err != nil {
		result.Error = "render html: " + err.Error()
		return result, nil
	}

	markdown, err := conv.ConvertString(htmlBuf.String())
	if err != nil {
		result.Error = "convert: " + err.Error()
		return result, nil
	}
	result.Markdown = markdown

	return result, nil
}
