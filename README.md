# html-markdown-tools

Composable [Axiom](https://axiomide.com) nodes for HTML-to-Markdown conversion —
render HTML to configurable Markdown, strip HTML to plain text, and isolate +
convert a page's main article content — wrapping the MIT-licensed
[JohannesKaufmann/html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown)
(v1 line — see below) and
[go-readability](https://codeberg.org/readeck/go-readability) (a maintained
fork of the Go port of Mozilla's Readability) libraries.

Built for the Axiom marketplace (handle `christiangeorgelucas`).

This is the Markdown-conversion counterpart to
[`markdown-tools`](https://github.com/ChristianGLucas/markdown-tools) (which
renders Markdown to HTML), and a flow-neighbor of
[`html-sanitize-tools`](https://github.com/ChristianGLucas/html-sanitize-tools)
and [`html-tools`](https://github.com/ChristianGLucas/html-tools): the `html`
input field matches `html-sanitize-tools`' `SanitizeResult.html` output
field, and the `markdown` output field matches `markdown-tools`' `markdown`
input field, so `html-sanitize-tools -> html-markdown-tools -> markdown-tools`
composes in a flow with no adapter. (`html-tools` shares the domain but its
node outputs are structured extraction results, not a top-level `html`
string, so wiring from it needs an explicit field selection.)

## Nodes

- **ConvertToMarkdown** — HTML → Markdown, with full control over heading
  style (ATX/setext), bullet and emphasis markers, code-block fencing
  (indented/fenced, backtick/tilde), link style (inline/referenced, with
  full/collapsed/shortcut reference numbering), and GFM
  table/strikethrough/task-list rendering (each independently toggleable).
  `keep_html_tags`/`remove_html_tags` let a caller preserve specific tags as
  raw HTML or drop a tag and its content entirely (e.g. `nav`, `script`)
  before conversion.
- **StripToPlainText** — HTML → plain text. Drops
  `<script>`/`<style>`/`<noscript>`/`<head>` entirely; separates block-level
  elements with a newline; an optional `collapse_whitespace` flag normalizes
  runs of whitespace to single spaces. Unlike ExtractMainContentAsMarkdown,
  this is a pure DOM traversal — no boilerplate detection.
- **ExtractMainContentAsMarkdown** — HTML page → isolated article content, as
  both Markdown and plain text, plus the metadata Readability recovers along
  the way (title, byline, excerpt, site name, lead image, language, published
  time). `base_url`, if given, resolves the extracted content's relative
  links/images to absolute URLs — a pure string join, never a fetch.
- **IsReaderable** — a cheap pre-flight heuristic: does this page look like
  it has real article content, without paying for the full extraction pass?

Every node is stateless, deterministic for a fixed (`html`, `options`) pair,
and offline: **no node in this package ever makes a network request**, even
when the input HTML references remote URLs. Input is capped at 10 MiB;
oversized input returns a structured error before any parsing happens.
Malformed/unbalanced HTML does not error — both wrapped libraries recover
from broken markup the same permissive way a browser does, and this package
documents that as its contract rather than promising strict validation. An
unrecognized `MarkdownOptions` value (e.g. `heading_style: "setext!"`) DOES
return a structured error — a deliberate improvement over the wrapped
library's own behavior of silently logging a warning and proceeding with a
half-valid config.

## Why v1, not v2, of html-to-markdown

JohannesKaufmann/html-to-markdown's `main` branch is a from-scratch v2
rewrite; its `LinkStyle` option (inline vs. referenced links) exists in the
`Options` struct but its setter is commented out in the v2 source pending a
render-logic rewrite. Since configurable link style is a documented option
here, `ConvertToMarkdown` wraps the frozen, stable v1 line (`v1.6.0`)
instead, where link style is fully implemented. See
`THIRD_PARTY_NOTICES.md` for the verification detail.

## License

MIT. See `THIRD_PARTY_NOTICES.md` for the full list of wrapped
dependencies (all MIT/BSD-2-Clause/BSD-3-Clause/Apache-2.0 — no copyleft
anywhere in the deployed build closure) and their license texts.
