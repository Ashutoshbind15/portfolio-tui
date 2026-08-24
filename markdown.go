package main

import (
	"strings"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	gstyles "charm.land/glamour/v2/styles"
)

func renderMarkdown(src string, width int) string {
	if width < 20 {
		width = 20
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(markdownStyle()),
		glamour.WithWordWrap(width),
		glamour.WithChromaFormatter("terminal16m"),
	)
	if err != nil {
		return src
	}
	out, err := r.Render(src)
	if err != nil {
		return src
	}
	return strings.TrimRight(out, "\n")
}

func markdownStyle() ansi.StyleConfig {
	s := gstyles.DarkStyleConfig
	s.Document.Color = strPtr(colorText)
	s.Document.Margin = uintPtr(0)

	s.Heading.Color = strPtr(colorAccent)
	s.Heading.Bold = boolPtr(true)
	s.Heading.BlockSuffix = "\n"

	// Render as headings, not as markdown source (`## Why`).
	s.H1.Prefix = ""
	s.H1.Suffix = ""
	s.H1.BackgroundColor = nil
	s.H1.Color = strPtr(colorBright)
	s.H1.Bold = boolPtr(true)
	s.H2.Prefix = ""
	s.H3.Prefix = ""
	s.H4.Prefix = ""
	s.H5.Prefix = ""
	s.H6.Prefix = ""
	s.H6.Color = strPtr(colorMuted)
	s.H6.Bold = boolPtr(false)

	s.Link.Color = strPtr(colorAccent)
	s.LinkText.Color = strPtr(colorAccent)
	s.Image.Color = strPtr(colorAccent)
	s.ImageText.Color = strPtr(colorMuted)
	s.HorizontalRule.Color = strPtr(colorBorder)

	s.Code.Color = strPtr(colorCode)
	s.Code.BackgroundColor = strPtr(colorCodeBg)

	s.CodeBlock.Color = strPtr(colorCode)
	s.CodeBlock.BackgroundColor = strPtr(colorCodeBg)
	s.CodeBlock.Margin = uintPtr(1)
	s.CodeBlock.Chroma = markdownChroma()

	return s
}

func markdownChroma() *ansi.Chroma {
	bg := strPtr(colorCodeBg)
	return &ansi.Chroma{
		Background:          ansi.StylePrimitive{BackgroundColor: bg},
		Text:                ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		Error:               ansi.StylePrimitive{Color: strPtr(colorBright), BackgroundColor: bg},
		Comment:             ansi.StylePrimitive{Color: strPtr(colorFaint), BackgroundColor: bg},
		CommentPreproc:      ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		Keyword:             ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		KeywordReserved:     ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		KeywordNamespace:    ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		KeywordType:         ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		Operator:            ansi.StylePrimitive{Color: strPtr(colorMuted), BackgroundColor: bg},
		Punctuation:         ansi.StylePrimitive{Color: strPtr(colorMuted), BackgroundColor: bg},
		Name:                ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		NameBuiltin:         ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		NameTag:             ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		NameAttribute:       ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		NameClass:           ansi.StylePrimitive{Color: strPtr(colorBright), BackgroundColor: bg},
		NameConstant:        ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		NameDecorator:       ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		NameException:       ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		NameFunction:        ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		NameOther:           ansi.StylePrimitive{Color: strPtr(colorText), BackgroundColor: bg},
		Literal:             ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		LiteralNumber:       ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		LiteralDate:         ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		LiteralString:       ansi.StylePrimitive{Color: strPtr(colorCode), BackgroundColor: bg},
		LiteralStringEscape: ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		GenericDeleted:      ansi.StylePrimitive{Color: strPtr(colorWarn), BackgroundColor: bg},
		GenericEmph:         ansi.StylePrimitive{Italic: boolPtr(true), BackgroundColor: bg},
		GenericInserted:     ansi.StylePrimitive{Color: strPtr(colorAccent), BackgroundColor: bg},
		GenericStrong:       ansi.StylePrimitive{Bold: boolPtr(true), BackgroundColor: bg},
		GenericSubheading:   ansi.StylePrimitive{Color: strPtr(colorMuted), BackgroundColor: bg},
	}
}

func strPtr(s string) *string { return &s }
func uintPtr(u uint) *uint    { return &u }
func boolPtr(b bool) *bool    { return &b }
