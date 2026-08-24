---
title: A terminal portfolio
date: 2026-08-24
summary: Bubble Tea’s loop, Lip Gloss layout, a 30fps splash, and Glamour for posts.
---

# A terminal portfolio

Same copy as the site, different runtime. A Bubble Tea v2 app that
paints an alt-screen, optionally hosted over SSH with Wish. The
interesting part is not the sections — it is how the frame gets on
the glass.

## The loop

Bubble Tea is the Elm architecture in a terminal: `Init`, `Update`,
`View`. The program is a message pump. Keys, mouse, resizes, and
your own ticks arrive as `tea.Msg`. `Update` returns a new model and
a `Cmd` — a function that will produce a later message. `View`
returns a string (v2 wraps that in `tea.View` for mouse mode, title,
and terminal colors). Nothing draws itself. You recompute the whole
frame from state.

The root model owns every page. On each message it either handles
something global (quit, tab, theme, nav click) or forwards to the
active child. Splash frames are ignored once you leave the intro so
the 30fps ticker cannot leak into the portfolio. Page models look
the same: `Activate` when you land, `SetSize` when the window
changes, `applyTheme` when the palette swaps.

Two transports, one model. Locally it is `tea.NewProgram`. Over SSH,
Wish’s Bubble Tea middleware constructs a fresh `appModel` per
session, so each connection is an isolated loop.

```go
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.ctx.width, m.ctx.height = msg.Width, msg.Height
        m.setSizes()
        return m, nil
    case splashFrameMsg:
        if !m.showSplash {
            return m, nil
        }
        var cmd tea.Cmd
        m.splash, cmd = m.splash.Update(msg)
        return m, cmd
    }
    // keys, mouse, then the active page…
}
```

## Layout and styling

Lip Gloss is the CSS. Styles are functions that read a global
palette — accent, text, muted, border, background — and return a
fresh `lipgloss.Style`. There is no long-lived style cache. Theme
swap mutates the hex vars, then every widget restyles: lists,
tree, help, markdown, scrollbar.

The chrome is three vertical bands. Header and footer are full
width. Between them, content sits beside a right sidebar. Widths
are derived: sidebar frame (expanded 12 cols of inner nav, or 2
when collapsed, plus padding and a left border), content gets the
rest minus a small left pad. `JoinHorizontal` / `JoinVertical`
compose the strings; `Width` / `Height` / `MaxHeight` force the
frame to the window so a short page does not leave a hole.

Click targets are not CSS boxes. bubblezone marks spans in the
rendered string (`nav-home`, `brand`, `theme-cycle`) and `Scan`s
the output. A mouse click is hit-tested against those IDs. Drag
is handled separately: if the pointer moved, it is a selection,
not a click.

Palettes are named (Pumice, Phosphor, Ember, Harbor, Orchid,
Parchment). `]` cycles them. Parchment is the light one — help and
list delegates flip their light/dark defaults from `Theme.Light`.
The `tea.View` also sets the terminal background/foreground from
the same hex so the alt-screen matches the Lip Gloss fill.

## Animating the splash

The intro is its own loop. `tea.Tick` at 30fps sends
`splashFrameMsg`. Each frame advances `t`, then a Harmonica spring
(`FPS(30)`, frequency 7, damping 0.55) pulls the name banner from
a starting offset toward `y = 0`. The hint line fades in as the
spring settles, then pulses with a sine.

The background is not a sprite. Each terminal cell is a 2×4
braille pixel (U+2800 plus eight bit flags). A surface is three
sines — swell, chop, ripple — sampled at 2× the column count. A
ribbon of dots hangs below the crest; empty cells get sparse
sparkles. Crest cells take `colorAccent`; deeper ones mix toward
muted and faint. Runs of the same color flush through Lip Gloss
so the row is not one style per cell.

Banner art is ASCII, picked by terminal size (slant, doom, small,
or just the name). It is stamped onto the wave grid, then the
whole grid is joined into a string. Leave with enter, space, or a
click that was not a drag. Clicking `>_ name` in the header
resets the spring and starts the ticker again.

Home has a quieter tick: once a second, IST clock, not 30fps.

## Markdown

Posts are `//go:embed blogs/*.md`. Front matter is a few `key:
value` lines — title, date, summary — parsed with a scanner, not
a YAML library. The body is the rest.

Open a post and Glamour renders it into ANSI for the current
viewport width. The style starts from Glamour’s dark config, then
every token is remapped onto the active palette: headings to
accent/bright, code to `colorCode` on `colorCodeBg`, links to
accent, rules to border. Heading prefixes are stripped so you see
**Why**, not `## Why`. Chroma colors for fenced blocks follow the
same map. Theme change re-renders the open viewport; wrap width
follows `SetSize`.

The list is a Bubbles `list`; the post is a `viewport` with a
custom gutter scrollbar. The thumb (`▐`) is `height² / total`
rows, pinned to the bottom when you are there. The content block
is padded to a rectangle first so the gutter stays on the right
edge instead of clinging to the last glyph of a short line.

## The rest of the frame

Mouse mode is cell-motion so drag-select works. Selection walks
the last painted frame, uses ANSI-aware widths, inverts the range
with accent-on-background, and `tea.SetClipboard`s the plain
text. `j`/`k` and arrows move lists; enter opens a detail
viewport; esc pops back. Stack is a Bubbles tree with the same
palette restyle path.

None of this is a framework. It is one model, a handful of child
models, and strings that happen to be a UI.

```bash
ssh "ash@tuis.ashutoshbind.com"
```
