package nodes_test

import (
	"context"
	"strings"
	"testing"

	gen "christiangeorgelucas/html-markdown-tools/gen"
	"christiangeorgelucas/html-markdown-tools/nodes"
)

// articleHTML is a full page with real article shape: metadata in <head>,
// nav/footer boilerplate, and three substantial paragraphs (well over
// Readability's ~500-char threshold) inside <article>, including a
// relative link and a relative image for base_url resolution.
const articleHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>My Great Article - Example Site</title>
<meta name="description" content="A short excerpt about the article.">
<meta property="og:site_name" content="Example Site">
<meta name="author" content="Jane Doe">
<meta property="article:published_time" content="2026-01-15T10:00:00Z">
</head>
<body>
<nav><a href="/">Home</a> <a href="/about">About</a></nav>
<article>
<h1>My Great Article</h1>
<p>This is the first paragraph of a genuinely long article about something interesting and important, written to exceed the minimum character threshold that Mozilla Readability requires before it will consider a block of text worth extracting as the main content of the page.</p>
<p>This is the second paragraph, continuing the article with more substantive prose so that the overall character count comfortably clears five hundred characters, giving the heuristic scoring function enough signal to prefer this content block over the navigation and footer boilerplate elsewhere on the page.</p>
<p>A third paragraph adds even more text, including a <a href="/related-page">relative link</a> and an <img src="/images/photo.jpg" alt="a photo"> image, both of which should have their URLs resolved to absolute form when a base URL is supplied to the extraction call.</p>
</article>
<footer>Copyright 2026 Example Site. All rights reserved.</footer>
</body>
</html>`

// notArticleHTML is a bare login form — the canonical "not an article" page.
const notArticleHTML = `<!DOCTYPE html>
<html><head><title>Login</title></head>
<body>
<form>
<label>Username</label><input type="text">
<label>Password</label><input type="password">
<button>Log in</button>
</form>
</body></html>`

const emptyBodyHTML = `<html><body></body></html>`

// unlabeledFormHTML has no text of its own at all — no <label>s, just bare
// inputs and a button — unlike notArticleHTML above (which has <label>
// text and so, verified separately, still extracts a small amount of
// text). This is the fixture for the "no text content at all" carve-out
// extracted's doc comment describes.
const unlabeledFormHTML = `<!DOCTYPE html>
<html><head><title>Login</title></head>
<body>
<form>
<input type="text">
<input type="password">
<button>Log in</button>
</form>
</body></html>`

// A realistic article page: markdown, metadata, and URL resolution against
// base_url must all be correct together. Values (title, byline, excerpt,
// site_name, published_time) are exactly what the HTML fixture's <head>
// metadata declares, independent of any Readability internals — this is
// the oracle.
func TestExtractMainContentAsMarkdown_ArticlePage(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{
		Html:    articleHTML,
		BaseUrl: "https://example.com/blog/my-article",
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected node error: %s", got.Error)
	}
	if !got.Extracted {
		t.Fatalf("expected extracted=true for a real article page")
	}
	if got.Title != "My Great Article - Example Site" {
		t.Errorf("title = %q", got.Title)
	}
	if got.Byline != "Jane Doe" {
		t.Errorf("byline = %q", got.Byline)
	}
	if got.Excerpt != "A short excerpt about the article." {
		t.Errorf("excerpt = %q", got.Excerpt)
	}
	if got.SiteName != "Example Site" {
		t.Errorf("site_name = %q", got.SiteName)
	}
	if got.Language != "en" {
		t.Errorf("language = %q", got.Language)
	}
	if got.PublishedTime != "2026-01-15T10:00:00Z" {
		t.Errorf("published_time = %q", got.PublishedTime)
	}
	if got.Length <= 0 {
		t.Errorf("expected positive length, got %d", got.Length)
	}
	if !strings.Contains(got.Markdown, "first paragraph") {
		t.Errorf("markdown missing article body: %q", got.Markdown)
	}
	if strings.Contains(got.Markdown, "Copyright 2026 Example Site") {
		t.Errorf("markdown leaked footer boilerplate: %q", got.Markdown)
	}
	if strings.Contains(got.Markdown, "Home") && strings.Contains(got.Markdown, "About") {
		t.Errorf("markdown leaked nav boilerplate: %q", got.Markdown)
	}
	// base_url resolution: the relative href/src in the fixture must come
	// back absolute.
	if !strings.Contains(got.Markdown, "https://example.com/related-page") {
		t.Errorf("relative link was not resolved against base_url: %q", got.Markdown)
	}
	if !strings.Contains(got.TextContent, "first paragraph") {
		t.Errorf("text_content missing article body: %q", got.TextContent)
	}
}

// A page with no real content (empty body) must report extracted=false with
// empty markdown/text_content and NO error — this is a normal outcome, not
// a failure.
func TestExtractMainContentAsMarkdown_EmptyBodyNotExtracted(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{Html: emptyBodyHTML})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected node error for an empty-body page: %s", got.Error)
	}
	if got.Extracted {
		t.Errorf("expected extracted=false for an empty body")
	}
	if got.Markdown != "" || got.TextContent != "" {
		t.Errorf("expected empty markdown/text_content, got markdown=%q text=%q", got.Markdown, got.TextContent)
	}
}

// A page whose only elements carry no text of their own (unlabeled inputs
// and a button, no prose anywhere) must report extracted=false — this is
// the precise claim the extracted field's doc comment makes, verified
// literally rather than assumed. (A form WITH label text, like
// notArticleHTML used in IsReaderable's tests, is a different case — it
// still extracts a small amount of text and is not asserted here.)
func TestExtractMainContentAsMarkdown_UnlabeledFormNotExtracted(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{Html: unlabeledFormHTML})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected node error: %s", got.Error)
	}
	if got.Extracted {
		t.Errorf("expected extracted=false for a form with no text content of its own, got markdown=%q text=%q", got.Markdown, got.TextContent)
	}
}

// The Markdown output honors the same MarkdownOptions as ConvertToMarkdown
// — e.g. a custom heading style should show up in the rendered article.
func TestExtractMainContentAsMarkdown_HonorsMarkdownOptions(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{
		Html:    articleHTML,
		Options: &gen.MarkdownOptions{HeadingStyle: "setext"},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected node error: %s", got.Error)
	}
	if strings.Contains(got.Markdown, "# ") {
		t.Errorf("expected setext headings (no leading '# '), got %q", got.Markdown)
	}
}

func TestExtractMainContentAsMarkdown_BadBaseUrl(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{
		Html:    articleHTML,
		BaseUrl: "http://[::1]:namedport", // invalid: named port isn't numeric
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error for a malformed base_url")
	}
}

func TestExtractMainContentAsMarkdown_InvalidMarkdownOption(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{
		Html:    articleHTML,
		Options: &gen.MarkdownOptions{HeadingStyle: "banana"},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error for an invalid heading_style")
	}
}

// Input size is the platform's concern, not this node's — a large input
// must not crash, though it need not extract anything (no article-shaped
// content in a giant run of plain text).
func TestExtractMainContentAsMarkdown_LargeInputNoCrash(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	huge := "<p>" + strings.Repeat("a", 10*1024*1024+1) + "</p>"
	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, &gen.MainContentQuery{Html: huge})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error for large input: %s", got.Error)
	}
}

func TestExtractMainContentAsMarkdown_NilInput(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ExtractMainContentAsMarkdown(ctx, ax, nil)
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected an error for a nil request")
	}
}
