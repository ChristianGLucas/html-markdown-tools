package nodes

import (
	"context"

	"christiangeorgelucas/html-markdown-tools/axiom"
	gen "christiangeorgelucas/html-markdown-tools/gen"
)

// Convert an HTML document or fragment to Markdown using
// JohannesKaufmann/html-to-markdown, with full control over heading style
// (ATX/setext), bullet and emphasis markers, code-block fencing, link style
// (inline/referenced), and GFM table/strikethrough/task-list rendering. Pure
// string transformation — malformed HTML is parsed permissively (the same
// tolerant recovery a browser applies) rather than erroring, and no network
// request is ever made even when the input contains remote-looking URLs. An
// unrecognized MarkdownOptions value (e.g. heading_style="foo") returns a
// structured error instead of silently falling back to a default.
func ConvertToMarkdown(ctx context.Context, ax axiom.Context, input *gen.ConvertQuery) (*gen.ConvertResult, error) {
	if input == nil {
		return &gen.ConvertResult{Error: "request is required"}, nil
	}
	conv, err := buildConverter(input.Options)
	if err != nil {
		return &gen.ConvertResult{Error: err.Error()}, nil
	}

	markdown, err := conv.ConvertString(input.Html)
	if err != nil {
		return &gen.ConvertResult{Error: "convert: " + err.Error()}, nil
	}

	return &gen.ConvertResult{Markdown: markdown}, nil
}
