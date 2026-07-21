package nodes

import (
	"context"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
	"golang.org/x/net/html"

	"christiangeorgelucas/html-markdown-tools/axiom"
	gen "christiangeorgelucas/html-markdown-tools/gen"
)

// Cheaply check whether a page has the shape Readability considers a
// plausible article — enough text-bearing <p>/<pre>/<article> elements that
// aren't flagged as unlikely candidates by their class/id — WITHOUT paying
// for the full extraction ExtractMainContentAsMarkdown does. Useful as a
// pre-flight gate in a pipeline that only wants to run the heavier
// extraction on pages likely to have something worth extracting.
func IsReaderable(ctx context.Context, ax axiom.Context, input *gen.ReaderableQuery) (*gen.ReaderableResult, error) {
	if input == nil {
		return &gen.ReaderableResult{Error: "request is required"}, nil
	}
	if err := checkHTMLSize(input.Html); err != nil {
		return &gen.ReaderableResult{Error: err.Error()}, nil
	}

	doc, err := html.Parse(strings.NewReader(input.Html))
	if err != nil {
		return &gen.ReaderableResult{Error: "parse: " + err.Error()}, nil
	}

	ps := readability.NewParser()
	return &gen.ReaderableResult{IsReaderable: ps.CheckDocument(doc)}, nil
}
