package main

// Static portfolio content. Keep these short — pages render summaries,
// and enter opens a detail view. Blog bodies stay plain text for now;
// markdown files can replace them later.

type Profile struct {
	Name     string
	Role     string
	Email    string
	Website  string
	GitHub   string
	SSH      string
	Bullets  []string
}

func profile() Profile {
	return Profile{
		Name:    "Ashutosh Bind",
		Role:    "Indie hacker",
		Email:   "dev@ashutoshbind.com",
		Website: "ashutoshbind.com",
		GitHub:  "github.com/Ashutoshbind15",
		SSH:     "ssh ash@tuis.ashutoshbind.com",
		Bullets: []string{
			"Building around chess analysis, local LLMs, and design.",
			"Open sourcing tools and the journey along the way.",
		},
	}
}

type Project struct {
	ID          string
	Name        string
	Summary     string
	Detail      string
	GitHub      string
	Site        string
	Tech        []string
}

func projects() []Project {
	return []Project{
		{
			ID:      "scribblesvg",
			Name:    "ScribbleSVG",
			Summary: "Hand-drawn SVG diagrams — editor plus a drop-in renderer.",
			Detail:  "React editor and a renderer you can embed. Packages on npm under @scribblesvg/*.",
			GitHub:  "https://github.com/Ashutoshbind15/ScribbleSVG",
			Tech:    []string{"React", "TypeScript", "Node.js"},
		},
		{
			ID:      "varipane",
			Name:    "VariPane",
			Summary: "Compare prompts and models in parallel agent sandboxes.",
			Detail:  "Spin sandboxes, compare results, preview each through a proxied Vite URL.",
			GitHub:  "https://github.com/Ashutoshbind15/ai-ui-comparator-cloud",
			Site:    "https://varipane.com",
			Tech:    []string{"React", "TypeScript", "Node.js", "Vite", "Docker"},
		},
		{
			ID:      "sliceseeker",
			Name:    "SliceSeeker",
			Summary: "Search video with fused multimodal, transcript, and frame embeddings.",
			Detail:  "Segment and chunk video, then rank hits with RRF across embedding kinds.",
			GitHub:  "https://github.com/Ashutoshbind15/demo-search-ai",
			Tech:    []string{"React", "TypeScript", "PostgreSQL", "Gemini", "Docker"},
		},
		{
			ID:      "wafercms",
			Name:    "WaferCMS",
			Summary: "Self-hostable CMS with diagrams and rich text as first-class fields.",
			Detail:  "MIT-licensed. Admin agents help with collection work.",
			GitHub:  "https://github.com/Ashutoshbind15/WaferCMS",
			Tech:    []string{"React", "TypeScript", "Node.js", "Docker"},
		},
		{
			ID:      "termchess",
			Name:    "Term Chess",
			Summary: "SSH TUI chess — timed multiplayer or untimed bot games.",
			Detail:  "Wish + Bubble Tea server. Clocks, fingerprints, Postgres history.",
			GitHub:  "https://github.com/Ashutoshbind15/tern-chess",
			Site:    "https://termchess.ashutoshbind.com",
			Tech:    []string{"Go", "PostgreSQL", "Docker"},
		},
	}
}

func projectByID(id string) (Project, bool) {
	for _, p := range projects() {
		if p.ID == id {
			return p, true
		}
	}
	return Project{}, false
}

type Experience struct {
	ID         string
	Role       string
	Org        string
	Period     string
	Current    bool
	Summary    string
	Highlights []string
	Tech       []string
}

func experiences() []Experience {
	return []Experience{
		{
			ID:      "indie",
			Role:    "Indie hacking",
			Period:  "Jul 2025 – Present",
			Current: true,
			Summary: "Shipping products and prototypes — frameworks, chess analysis, local LLMs.",
			Highlights: []string{
				"Open-source tools and developer packages.",
				"SaaS experiments around chess, LLMs, and UI sandboxes.",
			},
			Tech: []string{"Go", "TypeScript", "Electron", "Open source"},
		},
		{
			ID:      "cisco-6m",
			Role:    "Internship",
			Org:     "Cisco Systems",
			Period:  "Jan 2025 – Jun 2025",
			Summary: "System health on Cisco Catalyst Center.",
			Highlights: []string{
				"Networking features in Java.",
				"Faster, more reliable unit tests.",
			},
			Tech: []string{"Java", "Microservices", "Testing"},
		},
		{
			ID:      "cisco-2m",
			Role:    "Internship",
			Org:     "Cisco Systems",
			Period:  "May 2024 – Jul 2024",
			Summary: "Cisco ISE team.",
			Highlights: []string{
				"APIs and dashboards for system metrics.",
			},
			Tech: []string{"Java", "Bash", "APIs"},
		},
	}
}

func experienceByID(id string) (Experience, bool) {
	for _, e := range experiences() {
		if e.ID == id {
			return e, true
		}
	}
	return Experience{}, false
}

type StackCategory struct {
	Name  string
	Items []string
}

func stackCategories() []StackCategory {
	return []StackCategory{
		{Name: "Frontend", Items: []string{"React", "Next.js", "TypeScript", "Tailwind", "Vite", "Electron"}},
		{Name: "Backend", Items: []string{"Node.js", "Go", "Java"}},
		{Name: "Data", Items: []string{"PostgreSQL", "Drizzle"}},
		{Name: "Ops", Items: []string{"Docker", "Linux", "Bash"}},
		{Name: "AI", Items: []string{"Gemini", "Ollama"}},
	}
}

type Blog struct {
	ID      string
	Title   string
	Date    string
	Summary string
	Body    string
}

func blogs() []Blog {
	return []Blog{
		{
			ID:      "tui-portfolio",
			Title:   "A terminal portfolio",
			Date:    "2026-08-24",
			Summary: "Same site, different surface — a TUI you can SSH into.",
			Body: "This is a placeholder. Posts will live as markdown here first,\n" +
				"then get adapted into the web portfolio.\n\n" +
				"For now the blogs page is just a list you can open.",
		},
		{
			ID:      "notes",
			Title:   "Notes",
			Date:    "2026-08-24",
			Summary: "Scratch space for short writeups.",
			Body:    "Nothing here yet. Open this item to see how a post will read.",
		},
	}
}

func blogByID(id string) (Blog, bool) {
	for _, b := range blogs() {
		if b.ID == id {
			return b, true
		}
	}
	return Blog{}, false
}
