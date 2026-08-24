package main

// Static portfolio content, kept in sync with the website home page.
// Home renders these as a single scrollable page with expandable rows;
// the other pages reuse the same data as lists, details, and a tree.

type Profile struct {
	Name     string
	Role     string
	Email    string
	Website  string
	GitHub   string
	SSH      string
	Location string
	Bullets  []string
	Contact  []string
}

func profile() Profile {
	return Profile{
		Name:     "Ashutosh Bind",
		Role:     "Indie hacker",
		Email:    "dev@ashutoshbind.com",
		Website:  "ashutoshbind.com",
		GitHub:   "github.com/Ashutoshbind15",
		SSH:      "ssh ash@tuis.ashutoshbind.com",
		Location: "Navsari, India",
		Bullets: []string{
			"I create end-to-end web and AI apps with TS, React and other tools",
			"Building VaryPane, SSVG, SliceSeeker and more",
		},
		Contact: []string{
			"Available for freelance, contract-based, or other interesting work",
			"I'll try to reply within 24 hours",
		},
	}
}

type Project struct {
	ID      string
	Name    string
	Kind    string
	Summary string
	Detail  string
	GitHub  string
	Site    string
	SSH     string
	Tech    []string
}

func projects() []Project {
	return []Project{
		{
			ID:      "varipane",
			Name:    "VaryPane",
			Kind:    "SaaS",
			Summary: "Parallel UI-generation agents in the cloud",
			Detail:  "Cloud-hosted parallel UI agents with starter templates, prompt libraries, and streaming previews.",
			GitHub:  "https://github.com/Ashutoshbind15/ai-ui-comparator-cloud",
			Site:    "https://varipane.com",
			Tech:    []string{"React", "TypeScript", "Node.js", "Vite", "Docker"},
		},
		{
			ID:      "scribblesvg",
			Name:    "ScribbleSVG",
			Kind:    "OSS",
			Summary: "Hand-drawn diagramming toolkit",
			Detail:  "TypeScript core and React utilities for hand-drawn SVG diagrams on an infinite canvas.",
			GitHub:  "https://github.com/Ashutoshbind15/ScribbleSVG",
			Tech:    []string{"React", "TypeScript"},
		},
		{
			ID:      "sliceseeker",
			Name:    "SliceSeeker",
			Kind:    "OSS",
			Summary: "Self-hostable semantic search inside long-form video",
			Detail:  "Ask in plain text, get back exact moments from long-form video via fused transcript and frame embeddings.",
			GitHub:  "https://github.com/Ashutoshbind15/demo-search-ai",
			Tech:    []string{"React", "TypeScript", "PostgreSQL", "Gemini", "Docker"},
		},
		{
			ID:      "termchess",
			Name:    "TermChess",
			Kind:    "Service",
			Summary: "Chess over SSH, in your terminal",
			Detail:  "Real-time multiplayer chess with time controls, or solo games against rated bots: all over SSH.",
			GitHub:  "https://github.com/Ashutoshbind15/tern-chess",
			Site:    "https://termchess.ashutoshbind.com",
			SSH:     "ssh -p 58303 tchess.ashutoshbind.com",
			Tech:    []string{"Go", "PostgreSQL", "Docker"},
		},
		{
			ID:      "wafercms",
			Name:    "WaferCMS",
			Kind:    "OSS",
			Summary: "Self-hostable CMS for complex content",
			Detail:  "Schema-driven CMS with rich text, diagram fields, and opt-in AI agents: self-hostable.",
			GitHub:  "https://github.com/Ashutoshbind15/WaferCMS",
			Tech:    []string{"React", "TypeScript", "Node.js", "Docker"},
		},
		{
			ID:      "aiuicomparator",
			Name:    "AI UI Comparator",
			Kind:    "OSS",
			Summary: "Local side-by-side AI UI generation",
			Detail:  "Compare how AI models generate UI locally: isolated sandboxes, reusable prompts, side-by-side results.",
			GitHub:  "https://github.com/Ashutoshbind15/ai-ui-comparator",
			Tech:    []string{"React", "TypeScript", "Electron"},
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
			Role:    "Indie Hacking",
			Period:  "Jul 2025 – Present",
			Current: true,
			Summary: "Product design and end-to-end software development",
			Highlights: []string{
				"Open-source tools and developer packages.",
				"SaaS experiments around chess, LLMs, and UI sandboxes.",
			},
			Tech: []string{
				"TypeScript", "React", "Next.js", "Go", "PostgreSQL",
				"Docker", "DevOps", "AI Agents", "Local LLMs", "Open Source",
			},
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

type Education struct {
	Degree string
	School string
	Period string
	Note   string
}

func education() Education {
	return Education{
		Degree: "B.Tech in Computer Science and Engineering",
		School: "National Institute of Technology Patna",
		Period: "Dec 2021 – Jun 2025",
		Note:   "CGPA 9.44",
	}
}

type StackCategory struct {
	Name  string
	Items []string
}

func stackCategories() []StackCategory {
	return []StackCategory{
		{Name: "Frontend", Items: []string{"React", "Next.js", "Tailwind CSS", "shadcn/ui"}},
		{Name: "Backend & Databases", Items: []string{"Node.js", "TypeScript", "Go", "Java", "Electron", "PostgreSQL", "MongoDB", "Drizzle"}},
		{Name: "Cloud & Ops", Items: []string{"AWS", "Railway", "Docker", "Caddy", "Linux", "Bash"}},
		{Name: "Testing", Items: []string{"JUnit", "Vitest"}},
		{Name: "Workflow & AI", Items: []string{"Git", "GitHub", "pnpm", "Cursor", "Gemini", "Qwen", "OpenCode"}},
	}
}

func stackToolCount() int {
	n := 0
	for _, cat := range stackCategories() {
		n += len(cat.Items)
	}
	return n
}
