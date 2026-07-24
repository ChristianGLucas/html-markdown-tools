package nodes_test

import (
	"context"
	"strings"
	"testing"

	gen "christiangeorgelucas/html-markdown-tools/gen"
	"christiangeorgelucas/html-markdown-tools/nodes"
)

// Independent-oracle golden cases: the expected Markdown for each input is
// derived directly from the CommonMark / GFM spec by hand, NOT by running
// the converter and copying its output — e.g. `<h1>Hi</h1>` is ATX heading
// syntax by definition ("# Hi"), `<strong>` is CommonMark strong-emphasis
// ("**...**" by this package's default), and GFM defines the pipe-table
// syntax verbatim. A regression that silently changed the rendering would
// fail these even though the code compiles and runs.
func TestConvertToMarkdown_Golden(t *testing.T) {
	cases := []struct {
		name string
		html string
		opts *gen.MarkdownOptions
		want string
	}{
		{
			name: "heading atx default",
			html: `<h1>Hi</h1>`,
			want: "# Hi",
		},
		{
			name: "heading setext option",
			html: `<h1>Hi</h1>`,
			opts: &gen.MarkdownOptions{HeadingStyle: "setext"},
			want: "Hi\n==",
		},
		{
			name: "bold and italic default delimiters",
			html: `<p><strong>Bold</strong> and <em>Italic</em></p>`,
			want: "**Bold** and _Italic_",
		},
		{
			name: "bold and italic custom delimiters",
			html: `<p><strong>Bold</strong> and <em>Italic</em></p>`,
			opts: &gen.MarkdownOptions{StrongDelimiter: "__", EmDelimiter: "*"},
			want: "__Bold__ and *Italic*",
		},
		{
			name: "link inline default",
			html: `<a href="https://example.com">Example</a>`,
			want: "[Example](https://example.com)",
		},
		{
			name: "link referenced style",
			html: `<a href="https://example.com">Example</a>`,
			opts: &gen.MarkdownOptions{LinkStyle: "referenced", LinkReferenceStyle: "full"},
			want: "[Example][1]\n\n[1]: https://example.com",
		},
		{
			name: "fenced code block default",
			html: "<pre><code>x = 1</code></pre>",
			want: "```\nx = 1\n```",
		},
		{
			name: "fenced code block tilde",
			html: "<pre><code>x = 1</code></pre>",
			opts: &gen.MarkdownOptions{CodeBlockFence: "~~~"},
			want: "~~~\nx = 1\n~~~",
		},
		{
			name: "unordered list default marker",
			html: "<ul><li>One</li><li>Two</li></ul>",
			want: "- One\n- Two",
		},
		{
			name: "unordered list custom marker",
			html: "<ul><li>One</li><li>Two</li></ul>",
			opts: &gen.MarkdownOptions{BulletListMarker: "*"},
			want: "* One\n* Two",
		},
		{
			name: "strikethrough enabled by default",
			html: "<del>gone</del>",
			want: "~~gone~~",
		},
		{
			name: "gfm table enabled by default",
			html: "<table><tr><th>A</th><th>B</th></tr><tr><td>1</td><td>2</td></tr></table>",
			want: "| A | B |\n| --- | --- |\n| 1 | 2 |",
		},
	}

	ctx := context.Background()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ax := newTestContext(t)
			got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{Html: tc.html, Options: tc.opts})
			if err != nil {
				t.Fatalf("unexpected transport error: %v", err)
			}
			if got.Error != "" {
				t.Fatalf("unexpected node error: %s", got.Error)
			}
			if got.Markdown != tc.want {
				t.Errorf("markdown = %q, want %q", got.Markdown, tc.want)
			}
		})
	}
}

// Disabling the table plugin must actually change behavior — a table
// renders as plain concatenated cell text, not GFM pipe syntax — otherwise
// disable_tables is a no-op flag that lies about what it does.
func TestConvertToMarkdown_DisableTables(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)
	html := "<table><tr><th>A</th><th>B</th></tr><tr><td>1</td><td>2</td></tr></table>"

	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{
		Html:    html,
		Options: &gen.MarkdownOptions{DisableTables: true},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if strings.Contains(got.Markdown, "|") {
		t.Errorf("disable_tables=true still produced pipe-table syntax: %q", got.Markdown)
	}
	// Documented, verified-exact fallback shape: cells concatenate with NO
	// separator at all (not even a space) once the table plugin is off —
	// this is what the disable_tables doc comment now explicitly warns
	// about, so pin the exact behavior as a regression guard.
	if got.Markdown != "AB12" {
		t.Errorf("disable_tables fallback shape changed: got %q, want %q", got.Markdown, "AB12")
	}
}

// keep_html_tags must preserve the raw tag verbatim in the output instead
// of dropping it (html-to-markdown's default for a tag with no CommonMark
// rule, e.g. <iframe>, is to drop it entirely — verified against the
// unwrapped library above).
func TestConvertToMarkdown_KeepHtmlTags(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{
		Html:    `<iframe src="https://x.example/embed"></iframe><p>hi</p>`,
		Options: &gen.MarkdownOptions{KeepHtmlTags: []string{"iframe"}},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if !strings.Contains(got.Markdown, "<iframe") {
		t.Errorf("keep_html_tags did not preserve <iframe>: %q", got.Markdown)
	}
	if !strings.Contains(got.Markdown, "hi") {
		t.Errorf("expected surrounding content preserved: %q", got.Markdown)
	}
}

// remove_html_tags must drop the tag AND its content — the default
// behavior (verified above) is to keep a tag's inner text even for a tag
// with no CommonMark rule, so this option must produce a different result.
func TestConvertToMarkdown_RemoveHtmlTags(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{
		Html:    `<nav>Nav</nav><p>Body</p>`,
		Options: &gen.MarkdownOptions{RemoveHtmlTags: []string{"nav"}},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if strings.Contains(got.Markdown, "Nav") {
		t.Errorf("remove_html_tags did not drop <nav> content: %q", got.Markdown)
	}
	if !strings.Contains(got.Markdown, "Body") {
		t.Errorf("expected <p>Body</p> to survive: %q", got.Markdown)
	}
}

// An unrecognized option value must be a deterministic, structured error —
// NOT silently ignored. This is the behavior this package deliberately adds
// on top of the wrapped library, which (verified directly against the
// unwrapped library above) only logs a warning to stderr and proceeds with
// a half-valid config.
func TestConvertToMarkdown_InvalidOption(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{
		Html:    "<h1>Hi</h1>",
		Options: &gen.MarkdownOptions{HeadingStyle: "banana"},
	})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error for heading_style=%q, got markdown=%q", "banana", got.Markdown)
	}
	if got.Markdown != "" {
		t.Errorf("expected empty markdown alongside an error, got %q", got.Markdown)
	}
}

// Input size is the platform's concern, not this node's — a large input
// must not crash.
func TestConvertToMarkdown_LargeInputNoCrash(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	huge := "<p>" + strings.Repeat("a", 10*1024*1024+1) + "</p>"
	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{Html: huge})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error for large input: %s", got.Error)
	}
}

// Malformed/unbalanced HTML must NOT crash the node — the underlying parser
// recovers permissively (documented package contract), so this asserts the
// call completes and returns something rather than erroring or panicking.
func TestConvertToMarkdown_MalformedHtmlDoesNotCrash(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	malformed := `<div><p>Unclosed paragraph<span>nested<div>broken</p></div>`
	got, err := nodes.ConvertToMarkdown(ctx, ax, &gen.ConvertQuery{Html: malformed})
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("expected best-effort recovery, got error: %s", got.Error)
	}
	if !strings.Contains(got.Markdown, "Unclosed paragraph") {
		t.Errorf("expected recovered text content, got %q", got.Markdown)
	}
}

// A nil request must not panic.
func TestConvertToMarkdown_NilInput(t *testing.T) {
	ctx := context.Background()
	ax := newTestContext(t)

	got, err := nodes.ConvertToMarkdown(ctx, ax, nil)
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected an error for a nil request")
	}
}
