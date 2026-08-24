package main

import (
	"strings"
	"sync"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	gstyles "charm.land/glamour/v2/styles"
	"github.com/alecthomas/chroma/v2"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
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
	// Glamour registers custom chroma palettes once under the name "charm"
	// and never updates them. Point it at a per-theme style instead.
	s.CodeBlock.Theme = registerMarkdownChroma()
	s.CodeBlock.Chroma = nil

	return s
}

var chromaMu sync.Mutex

func chromaThemeName() string {
	return "tui-pf-" + currentTheme().ID
}

func registerMarkdownChroma() string {
	name := chromaThemeName()
	c := markdownChroma()
	chromaMu.Lock()
	chromastyles.Register(chroma.MustNewStyle(name, chroma.StyleEntries{
		chroma.Background:          chromaEntry(c.Background),
		chroma.Text:                chromaEntry(c.Text),
		chroma.Error:               chromaEntry(c.Error),
		chroma.Comment:             chromaEntry(c.Comment),
		chroma.CommentPreproc:      chromaEntry(c.CommentPreproc),
		chroma.Keyword:             chromaEntry(c.Keyword),
		chroma.KeywordReserved:     chromaEntry(c.KeywordReserved),
		chroma.KeywordNamespace:    chromaEntry(c.KeywordNamespace),
		chroma.KeywordType:         chromaEntry(c.KeywordType),
		chroma.Operator:            chromaEntry(c.Operator),
		chroma.Punctuation:         chromaEntry(c.Punctuation),
		chroma.Name:                chromaEntry(c.Name),
		chroma.NameBuiltin:         chromaEntry(c.NameBuiltin),
		chroma.NameTag:             chromaEntry(c.NameTag),
		chroma.NameAttribute:       chromaEntry(c.NameAttribute),
		chroma.NameClass:           chromaEntry(c.NameClass),
		chroma.NameConstant:        chromaEntry(c.NameConstant),
		chroma.NameDecorator:       chromaEntry(c.NameDecorator),
		chroma.NameException:       chromaEntry(c.NameException),
		chroma.NameFunction:        chromaEntry(c.NameFunction),
		chroma.NameOther:           chromaEntry(c.NameOther),
		chroma.Literal:             chromaEntry(c.Literal),
		chroma.LiteralNumber:       chromaEntry(c.LiteralNumber),
		chroma.LiteralDate:         chromaEntry(c.LiteralDate),
		chroma.LiteralString:       chromaEntry(c.LiteralString),
		chroma.LiteralStringEscape: chromaEntry(c.LiteralStringEscape),
		chroma.GenericDeleted:      chromaEntry(c.GenericDeleted),
		chroma.GenericEmph:         chromaEntry(c.GenericEmph),
		chroma.GenericInserted:     chromaEntry(c.GenericInserted),
		chroma.GenericStrong:       chromaEntry(c.GenericStrong),
		chroma.GenericSubheading:   chromaEntry(c.GenericSubheading),
	}))
	chromaMu.Unlock()
	return name
}

func chromaEntry(p ansi.StylePrimitive) string {
	var parts []string
	if p.Color != nil {
		parts = append(parts, *p.Color)
	}
	if p.BackgroundColor != nil {
		parts = append(parts, "bg:"+*p.BackgroundColor)
	}
	if p.Italic != nil && *p.Italic {
		parts = append(parts, "italic")
	}
	if p.Bold != nil && *p.Bold {
		parts = append(parts, "bold")
	}
	if p.Underline != nil && *p.Underline {
		parts = append(parts, "underline")
	}
	return strings.Join(parts, " ")
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
