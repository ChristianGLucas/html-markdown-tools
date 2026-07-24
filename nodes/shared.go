package nodes

import (
	"fmt"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"

	gen "christiangeorgelucas/html-markdown-tools/gen"
)

// oneOf validates that v is either empty (caller wants this package's
// default) or exactly one of allowed. Returns a structured error naming the
// field, the allowed values, and what was actually given otherwise — this is
// what makes an unrecognized option value a deterministic INVALID_ARGUMENT
// rather than the underlying library silently logging a warning and
// proceeding with a half-valid config (which is what
// JohannesKaufmann/html-to-markdown's own NewConverter does).
func oneOf(field, v string, allowed ...string) error {
	if v == "" {
		return nil
	}
	for _, a := range allowed {
		if v == a {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %v, got %q", field, allowed, v)
}

// buildConverter validates a MarkdownOptions request and constructs a
// configured html-to-markdown Converter (v1 API — chosen specifically
// because v2's LinkStyle option is not yet implemented, see README).
//
// Defaults intentionally differ from the wrapped library's own zero-value
// defaults in exactly one place: CodeBlockStyle defaults to "fenced" here,
// not the library's "indented" — fenced code blocks (with the language tag
// preserved) are what virtually every modern Markdown consumer (GitHub,
// CommonMark viewers, LLM prompts) expects, whereas "indented" is the
// historical Markdown.pl default and silently drops the language tag.
// Every other field's default matches the library's own default exactly.
func buildConverter(o *gen.MarkdownOptions) (*md.Converter, error) {
	if o == nil {
		o = &gen.MarkdownOptions{}
	}

	if err := oneOf("heading_style", o.HeadingStyle, "atx", "setext"); err != nil {
		return nil, err
	}
	if err := oneOf("bullet_list_marker", o.BulletListMarker, "-", "+", "*"); err != nil {
		return nil, err
	}
	if err := oneOf("em_delimiter", o.EmDelimiter, "*", "_"); err != nil {
		return nil, err
	}
	if err := oneOf("strong_delimiter", o.StrongDelimiter, "**", "__"); err != nil {
		return nil, err
	}
	if err := oneOf("code_block_style", o.CodeBlockStyle, "fenced", "indented"); err != nil {
		return nil, err
	}
	if err := oneOf("code_block_fence", o.CodeBlockFence, "```", "~~~"); err != nil {
		return nil, err
	}
	if err := oneOf("link_style", o.LinkStyle, "inline", "referenced"); err != nil {
		return nil, err
	}
	if err := oneOf("link_reference_style", o.LinkReferenceStyle, "full", "collapsed", "shortcut"); err != nil {
		return nil, err
	}
	opts := &md.Options{
		HeadingStyle:       o.HeadingStyle,
		BulletListMarker:   o.BulletListMarker,
		EmDelimiter:        o.EmDelimiter,
		StrongDelimiter:    o.StrongDelimiter,
		LinkReferenceStyle: o.LinkReferenceStyle,
	}
	// CodeBlockStyle: our default ("fenced") differs from the library's
	// ("indented"), so — unlike the other fields above — we must set it
	// explicitly rather than leaving "" for the library to fill in.
	if o.CodeBlockStyle == "" {
		opts.CodeBlockStyle = "fenced"
	} else {
		opts.CodeBlockStyle = o.CodeBlockStyle
	}
	opts.Fence = o.CodeBlockFence

	switch o.LinkStyle {
	case "referenced":
		opts.LinkStyle = "referenced"
	default:
		// "" or "inline" both map to the library's "inlined" — our public
		// vocabulary says "inline", the library's internal string is
		// "inlined"; this is the one place the two differ.
		opts.LinkStyle = "inlined"
	}

	conv := md.NewConverter("", true, opts)

	var plugins []md.Plugin
	if !o.DisableTables {
		plugins = append(plugins, plugin.Table())
	}
	if !o.DisableStrikethrough {
		plugins = append(plugins, plugin.Strikethrough(""))
	}
	if !o.DisableTaskLists {
		plugins = append(plugins, plugin.TaskListItems())
	}
	if len(plugins) > 0 {
		conv.Use(plugins...)
	}
	if len(o.KeepHtmlTags) > 0 {
		conv.Keep(o.KeepHtmlTags...)
	}
	if len(o.RemoveHtmlTags) > 0 {
		conv.Remove(o.RemoveHtmlTags...)
	}

	return conv, nil
}
