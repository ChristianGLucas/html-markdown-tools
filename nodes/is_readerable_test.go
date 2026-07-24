package nodes_test

import (
	"context"
	"strings"
	"testing"

	gen "christiangeorgelucas/html-markdown-tools/gen"
	"christiangeorgelucas/html-markdown-tools/nodes"
)

// A real article page (see articleHTML, shared with
// extract_main_content_as_markdown_test.go) must report readerable=true;
// a bare login form (notArticleHTML) must report false. This is the same
// oracle pair ExtractMainContentAsMarkdown's own tests use, letting the two
// nodes' behavior be checked for consistency: IsReaderable is documented as
// a cheap pre-flight for ExtractMainContentAsMarkdown, so it would be a
// real defect for one to say "readerable" and the other to disagree.
func TestIsReaderable_Golden(t *testing.T) {
	ctx := context.Background()

	t.Run("article page is readerable", func(t *testing.T) {
		ax := newTestContext(t)
		got, err := nodes.IsReaderable(ctx, ax, &gen.ReaderableQuery{Html: articleHTML})
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if got.Error != "" {
			t.Fatalf("unexpected node error: %s", got.Error)
		}
		if !got.IsReaderable {
			t.Errorf("expected a real article page to be readerable")
		}
	})

	t.Run("login form is not readerable", func(t *testing.T) {
		ax := newTestContext(t)
		got, err := nodes.IsReaderable(ctx, ax, &gen.ReaderableQuery{Html: notArticleHTML})
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if got.Error != "" {
			t.Fatalf("unexpected node error: %s", got.Error)
		}
		if got.IsReaderable {
			t.Errorf("expected a bare login form to NOT be readerable")
		}
	})
}

// Input size is the platform's concern, not this node's — a large input
// must not crash.
func TestIsReaderable_LargeInputNoCrash(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	huge := "<p>" + strings.Repeat("a", 10*1024*1024+1) + "</p>"
	got, err := nodes.IsReaderable(ctx, ax, &gen.ReaderableQuery{Html: huge})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error for large input: %s", got.Error)
	}
}

func TestIsReaderable_NilInput(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.IsReaderable(ctx, ax, nil)
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected an error for a nil request")
	}
}
