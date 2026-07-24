package nodes

import (
	"context"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"golang.org/x/net/html"

	"christiangeorgelucas/html-markdown-tools/axiom"
	gen "christiangeorgelucas/html-markdown-tools/gen"
)

// skipSubtree names elements whose entire subtree (including all descendant
// text nodes) is invisible/non-content and must be dropped, not just
// unwrapped: script/style bodies aren't prose, and head metadata isn't
// rendered by a browser either.
var skipSubtree = map[string]bool{
	"script":   true,
	"style":    true,
	"noscript": true,
	"head":     true,
	"template": true,
}

// Reduce an HTML document or fragment to its visible plain text: every
// <script>/<style>/<noscript>/<head> subtree is dropped, and every other
// text node is kept. Block-level elements (per html-to-markdown's own
// IsInlineElement classification, e.g. <p>, <div>, <li>, <h1>) are separated
// by a newline in the output so adjacent blocks don't run together; inline
// elements (<span>, <a>, <strong>, ...) are not. This is a pure DOM
// traversal — no markup semantics beyond block/inline are applied, so it
// does not attempt to identify a "main content" region the way
// ExtractMainContentAsMarkdown does.
func StripToPlainText(ctx context.Context, ax axiom.Context, input *gen.PlainTextQuery) (*gen.PlainTextResult, error) {
	if input == nil {
		return &gen.PlainTextResult{Error: "request is required"}, nil
	}
	doc, err := html.Parse(strings.NewReader(input.Html))
	if err != nil {
		return &gen.PlainTextResult{Error: "parse: " + err.Error()}, nil
	}

	var sb strings.Builder
	extractText(doc, &sb)
	text := sb.String()

	if input.CollapseWhitespace {
		text = strings.Join(strings.Fields(text), " ")
	}

	return &gen.PlainTextResult{Text: text}, nil
}

// extractText appends n's visible text (and its children's) to sb,
// bracketing any block-level element's subtree with newlines — BOTH before
// and after its children — so a block is separated from surrounding text on
// both sides, not just from its following sibling. Without the leading
// newline, text immediately preceding a nested block (e.g. "A" in
// "<div>A<div>B</div></div>") would run directly into the block's own text
// ("AB") with no separator at all.
func extractText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.ElementNode && skipSubtree[n.Data] {
		return
	}
	isBlock := n.Type == html.ElementNode && !md.IsInlineElement(n.Data)

	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
	}
	if isBlock {
		sb.WriteString("\n")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, sb)
	}

	if isBlock {
		sb.WriteString("\n")
	}
}
