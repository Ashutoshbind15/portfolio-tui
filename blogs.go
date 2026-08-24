package main

import (
	"bufio"
	"embed"
	"io/fs"
	"path"
	"sort"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

//go:embed blogs/*.md
var blogFS embed.FS

type Blog struct {
	ID      string
	Title   string
	Date    string
	Summary string
	Body    string
}

var cachedBlogs []Blog

func init() {
	cachedBlogs = loadBlogs()
}

func blogs() []Blog { return cachedBlogs }

func blogByID(id string) (Blog, bool) {
	for _, b := range blogs() {
		if b.ID == id {
			return b, true
		}
	}
	return Blog{}, false
}

func loadBlogs() []Blog {
	entries, err := fs.ReadDir(blogFS, "blogs")
	if err != nil {
		return nil
	}
	out := make([]Blog, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		data, err := blogFS.ReadFile(path.Join("blogs", entry.Name()))
		if err != nil {
			continue
		}
		out = append(out, parseBlogFile(entry.Name(), string(data)))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Date == out[j].Date {
			return out[i].Title < out[j].Title
		}
		return out[i].Date > out[j].Date
	})
	return out
}

func parseBlogFile(filename, raw string) Blog {
	id := strings.TrimSuffix(filename, ".md")
	b := Blog{ID: id, Title: id}

	body := raw
	if strings.HasPrefix(raw, "---\n") {
		rest := strings.TrimPrefix(raw, "---\n")
		end := strings.Index(rest, "\n---\n")
		if end >= 0 {
			applyFrontMatter(&b, rest[:end])
			body = rest[end+len("\n---\n"):]
		}
	}
	b.Body = strings.TrimSpace(body)
	if b.Title == id {
		if line := firstHeading(b.Body); line != "" {
			b.Title = line
		}
	}
	return b
}

func applyFrontMatter(b *Blog, matter string) {
	sc := bufio.NewScanner(strings.NewReader(matter))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "title":
			b.Title = val
		case "date":
			b.Date = val
		case "summary":
			b.Summary = val
		}
	}
}

func firstHeading(body string) string {
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

type blogItem struct {
	blog Blog
}

func (i blogItem) FilterValue() string { return i.blog.Title }
func (i blogItem) Title() string       { return i.blog.Title }
func (i blogItem) Description() string {
	if i.blog.Date == "" {
		return i.blog.Summary
	}
	if i.blog.Summary == "" {
		return i.blog.Date
	}
	return i.blog.Date + "  ·  " + i.blog.Summary
}

type blogsModel struct {
	ctx      *Context
	list     list.Model
	viewport viewport.Model
	openID   string
}

func newBlogsModel(ctx *Context) blogsModel {
	items := make([]list.Item, 0, len(blogs()))
	for _, b := range blogs() {
		items = append(items, blogItem{blog: b})
	}
	return blogsModel{
		ctx:      ctx,
		list:     newItemList("Blogs", items, 0, 0),
		viewport: viewport.New(),
	}
}

func (m blogsModel) Init() tea.Cmd { return nil }

func (m blogsModel) Activate() (blogsModel, tea.Cmd) { return m, nil }

func (m *blogsModel) SetSize(width, height int) {
	m.list.SetSize(width, height)
	sizeViewport(&m.viewport, width, height)
	if m.openID != "" {
		m.refreshDetail()
	}
}

func (m *blogsModel) refreshDetail() {
	b, ok := blogByID(m.openID)
	if !ok {
		return
	}
	m.viewport.SetContent(renderBlogDetail(b, m.viewport.Width()))
}

func (m blogsModel) Update(msg tea.Msg) (blogsModel, tea.Cmd) {
	if m.openID != "" {
		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "esc" {
			m.openID = ""
			return m, nil
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
		if it, ok := m.list.SelectedItem().(blogItem); ok {
			m.openID = it.blog.ID
			m.refreshDetail()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m blogsModel) View() string {
	if m.openID != "" {
		return viewWithScrollbar(m.viewport)
	}
	return m.list.View()
}

func renderBlogDetail(b Blog, width int) string {
	w := max(20, width-2)
	var sb strings.Builder
	if b.Date != "" {
		sb.WriteString(styleMuted().Render(b.Date))
		sb.WriteString("\n")
	}
	sb.WriteString(renderMarkdown(b.Body, w))
	sb.WriteString("\n\n")
	sb.WriteString(styleFaint().Render("esc back to the list"))
	return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render(sb.String())
}
