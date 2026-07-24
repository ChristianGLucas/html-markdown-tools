package nodes_test

import (
	"context"
	"strings"
	"testing"

	gen "christiangeorgelucas/html-markdown-tools/gen"
	"christiangeorgelucas/html-markdown-tools/nodes"
)

// Golden cases with collapse_whitespace=true, where the expected text is
// unambiguous and independently derivable from the input by hand.
func TestStripToPlainText_Golden(t *testing.T) {
	cases := []struct {
		name string
		html string
		want string
	}{
		{
			name: "script and style dropped entirely",
			html: `<script>alert(1)</script><style>.x{color:red}</style><p>Hello</p>`,
			want: "Hello",
		},
		{
			name: "sibling paragraphs separated by a space",
			html: `<p>Hello</p><p>World</p>`,
			want: "Hello World",
		},
		{
			name: "inline spans do not force a separator beyond the source whitespace",
			html: `<span>Hello</span> <span>World</span>`,
			want: "Hello World",
		},
		{
			name: "text before a nested block gets separated from the block's own text",
			html: `<div>A<div>B</div>C</div>`,
			want: "A B C",
		},
		{
			name: "plain text with markup tags stripped",
			html: `<p>Some <strong>bold</strong> and <em>italic</em> text.</p>`,
			want: "Some bold and italic text.",
		},
	}

	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ax := newTestContext(t)
			got, err := nodes.StripToPlainText(ctx, ax, &gen.PlainTextQuery{Html: tc.html, CollapseWhitespace: true})
			if err != nil {
				t.Fatalf("unexpected transport error: %v", err)
			}
			if got.Error != "" {
				t.Fatalf("unexpected node error: %s", got.Error)
			}
			if got.Text != tc.want {
				t.Errorf("text = %q, want %q", got.Text, tc.want)
			}
		})
	}
}

// Without collapse_whitespace, exact formatting is not part of the
// contract, but adjacent block-level content must still be separated by
// SOME whitespace — asserting the words never fuse together is the
// meaningful, format-independent claim.
func TestStripToPlainText_NoCollapseStillSeparatesBlocks(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.StripToPlainText(ctx, ax, &gen.PlainTextQuery{
		Html:               `<p>Hello</p><p>World</p>`,
		CollapseWhitespace: false,
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if strings.Contains(got.Text, "HelloWorld") {
		t.Errorf("blocks ran together with no separator: %q", got.Text)
	}
	if !strings.Contains(got.Text, "Hello") || !strings.Contains(got.Text, "World") {
		t.Errorf("expected both words present: %q", got.Text)
	}
}

// Input size is the platform's concern, not this node's — a large input
// must not crash.
func TestStripToPlainText_LargeInputNoCrash(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	huge := "<p>" + strings.Repeat("a", 10*1024*1024+1) + "</p>"
	got, err := nodes.StripToPlainText(ctx, ax, &gen.PlainTextQuery{Html: huge})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error for large input: %s", got.Error)
	}
}

func TestStripToPlainText_NilInput(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.StripToPlainText(ctx, ax, nil)
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected an error for a nil request")
	}
}
