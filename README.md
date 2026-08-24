# tui-pf

Terminal portfolio for [Ashutosh Bind](https://ashutoshbind.com). Same sections as the site — Home, Projects, Experience, Stack, Blogs — as a Bubble Tea app you can run locally or over SSH.

Navigation and chrome follow [Term Chess](https://github.com/Ashutoshbind15/tern-chess): a root model routes pages, Tab opens a page list, and lists open into a short detail view instead of dumping the full copy on the page.

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
| `tab` | Page list |
| `1`–`5` | Home / Projects / Experience / Stack / Blogs |
| `enter` | Open the selected list item |
| `esc` | Back (or close the page list) |
| `q` / `ctrl+c` | Quit |
| click | Header tabs (bubblezone) |

On Stack, `enter` / `h` / `l` open and close the new Bubbles tree.

## Stack

Charm v2 (`charm.land/*`), matching Term Chess libraries at current releases:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 — app model, declarative `tea.View`
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2 — layout and theme
- [Bubbles](https://github.com/charmbracelet/bubbles) v2.2 — `list`, `viewport`, `help`, and the new `tree`
- [Wish](https://github.com/charmbracelet/wish) v2 — optional SSH host
- [bubblezone](https://github.com/lrstanley/bubblezone) v2 — clickable header nav

Blogs are plain text in `content.go` for now. Drop markdown into `blogs/` later (Glamour) and adapt those posts into the web portfolio.

## Layout

```text
main.go          local TUI + optional Wish server
app.go           root router, header, footer
menu.go          Tab page picker
home.go          profile
projects.go      list → openable detail
experience.go    list → openable detail
stack.go         tree of tool groups
blogs.go         list → simple text post
content.go       short static copy
styles.go        shared Lip Gloss
```
