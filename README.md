# tui-pf

Terminal portfolio for [Ashutosh Bind](https://ashutoshbind.com). Same sections as the site — Home, Projects, Experience, Stack, Blogs — as a Bubble Tea app you can run locally or over SSH.

A right sidebar switches pages (`[` collapses it to icons). The app opens on a full-screen intro; `enter` (or a click) goes into the portfolio, and clicking `>_ Ashutosh Bind` in the header brings the intro back. Lists open into a short detail view instead of dumping the full copy on the page. Blog posts are markdown.

## Run

```bash
go run .
```

SSH (Wish), default `127.0.0.1:23235`:

```bash
go run . ssh
ssh localhost -p 23235
```

`TUI_PF_ENV=production` binds `0.0.0.0`. `TUI_PF_PORT` overrides the port.

## Keys

| Key | Action |
| --- | --- |
| `tab` / `shift+tab` | Next / previous page |
| `[` | Collapse / expand nav |
| `enter` | Open the selected list item (or leave the intro) |
| click name | Return to the intro |
| `esc` | Back |
| `q` / `ctrl+c` | Quit |
| click | Sidebar pages, or `<`/`>` to toggle nav |
| drag | Select text and copy |

On Stack, `enter` / `h` / `l` open and close the tree.

## Stack

Charm v2 (`charm.land/*`):

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 — app model, declarative `tea.View`
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2 — layout and theme
- [Bubbles](https://github.com/charmbracelet/bubbles) v2.2 — `list`, `viewport`, `help`, and `tree`
- [Glamour](https://github.com/charmbracelet/glamour) v2 — blog markdown
- [Wish](https://github.com/charmbracelet/wish) v2 — optional SSH host
- [Harmonica](https://github.com/charmbracelet/harmonica) — spring animation on the intro
- [bubblezone](https://github.com/lrstanley/bubblezone) v2 — clickable sidebar

Drop posts into `blogs/*.md` (title / date / summary front matter). They are embedded at build time.

## Layout

```text
main.go          local TUI + optional Wish server
app.go           root router, header, footer, intro splash
splash.go        opening screen: name banner, braille waves, enter to continue
select.go        drag-to-select text and copy
sidebar.go       right-hand page nav
scrollbar.go     viewport scrollbar when content overflows
home.go          profile
projects.go      list → openable detail
experience.go    list → openable detail
stack.go         tree of tool groups
blogs.go         list → markdown post
markdown.go      Glamour renderer
content.go       short static copy
styles.go        shared Lip Gloss
blogs/           markdown posts
```
