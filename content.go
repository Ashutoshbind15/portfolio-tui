package main

// Static portfolio content, kept in sync with portfolio-new.
// Home and the project list use the short card copy; opening a project
// uses the individual page text from the site (what / why / where /
// technicals / upcoming), without the gallery.

type Profile struct {
	Name     string
	Role     string
	Email    string
	Website  string
	GitHub   string
	X        string
	LinkedIn string
	SSH      string
	Location string
	Bullets  []string
	Contact  []string
}

type contactRow struct {
	Icon  string
	Value string
}

func profile() Profile {
	return Profile{
		Name:     "Ashutosh Bind",
		Role:     "Indie hacker",
		Email:    "dev@ashutoshbind.com",
		Website:  "ashutoshbind.com",
		GitHub:   "github.com/Ashutoshbind15",
		X:        "x.com/Ashutosh_Bind15",
		LinkedIn: "linkedin.com/in/ashutosh-bind-56806b22b",
		SSH:      "ssh ash@tuis.ashutoshbind.com",
		Location: "Navsari, India",
		Bullets: []string{
			"Solo developing apps around chess-analysis, local llms, design and much more",
			"Open Sourcing some of my projects, tools and journey",
		},
		Contact: []string{
			"Available for freelance, contract-based, or other interesting work",
			"I'll try to reply within 24 hours",
		},
	}
}

func contactRows(p Profile) []contactRow {
	rows := []contactRow{
		{"📧", p.Email},
		{"🐙", p.GitHub},
		{"𝕏", p.X},
		{"🔗", p.LinkedIn},
		{"🌐", p.Website},
		{">_", p.SSH},
	}
	out := rows[:0]
	for _, row := range rows {
		if row.Value != "" {
			out = append(out, row)
		}
	}
	return out
}

type NpmPackage struct {
	Name string
	URL  string
}

type ProjectLinkExtra struct {
	Label string
	URL   string
}

type Project struct {
	ID         string
	Name       string
	Kind       string
	Summary    string
	Tagline    string
	What       string
	Why        string
	Technicals string
	Upcoming   []string
	Detail     string
	GitHub     string
	Site       string
	SSH        string
	Packages   string
	Npm        []NpmPackage
	More       []ProjectLinkExtra
	Tech       []string
}

func projects() []Project {
	return []Project{
		{
			ID:         "scribblesvg",
			Name:       "ScribbleSVG",
			Kind:       "OSS",
			Summary:    "SVG diagramming toolkit with a hand-drawn look — React editor plus a renderer you can drop in anywhere.",
			Tagline:    "Hand-drawn diagramming toolkit",
			What:       "A featherlight diagramming toolkit that renders hand-drawn-style SVGs — a TypeScript core plus React utilities for an infinite canvas and a drop-in renderer.",
			Why:        "I plan and reason systems by drawing. I wanted something with fewer dependencies than the usual options, and a lighter base for a diagram editor built on existing browser SVG primitives.",
			Technicals: "Layered SVG on a simple state-driven canvas — undo/redo, multi-select, and custom SVG icon imports. Core ships with a single runtime dependency; the React layer peers on React and depends on the core.",
			Upcoming: []string{
				"Bring-your-own fonts",
				"Opt-in AI text-to-diagram agent on the playground",
				"Custom themes and canvas styles",
				"SVG pack imports and image-to-SVG tracing",
			},
			GitHub: "https://github.com/Ashutoshbind15/ScribbleSVG",
			Npm: []NpmPackage{
				{Name: "@scribblesvg/core", URL: "https://www.npmjs.com/package/@scribblesvg/core"},
				{Name: "@scribblesvg/react-utils", URL: "https://www.npmjs.com/package/@scribblesvg/react-utils"},
			},
			More: []ProjectLinkExtra{
				{Label: "play", URL: "https://scribblesvg.ashutoshbind.com"},
				{Label: "docs", URL: "https://scribblesvg-docs.ashutoshbind.com"},
			},
			Tech: []string{"React", "TypeScript", "Node.js"},
		},
		{
			ID:         "varipane",
			Name:       "VariPane",
			Kind:       "SaaS",
			Summary:    "Spin parallel agent sandboxes to compare prompts and models, then preview each result through a proxied Vite URL.",
			Tagline:    "Parallel UI-generation agents in the cloud",
			What:       "Parallel UI-generation agents in the cloud. Pick a starter template, keep a reusable prompt library, attach skills, and compare models side by side — results stream in real time. Runs pi-coding-agent inside every sandbox.",
			Why:        "Visual comparison — and picking the best bits out of generated UIs — is usually very helpful, especially for developers doing backend work, people new to design, or anyone who just doesn't have the time for a skill this demanding. The cloud version scales and improves that workflow for hosting and for end users.",
			Technicals: "BYOK for OpenRouter and Upstash Box. Parallel agent sandboxes per session, a dedicated preview gateway, labelled box ownership, and self-hosted secret storage via Infisical.",
			Upcoming: []string{
				"Screenshot-to-agent input — iterate from a wireframe or your current design",
				"GitHub repo connections and agent-loop iterations on one-shot designs",
				"Scripted environment setup",
			},
			GitHub: "https://github.com/Ashutoshbind15/ai-ui-comparator-cloud",
			Site:   "https://varipane.com",
			Tech:   []string{"React", "TypeScript", "Node.js", "Vite", "Docker"},
		},
		{
			ID:         "sliceseeker",
			Name:       "SliceSeeker",
			Kind:       "OSS",
			Summary:    "Segment and chunk video, then search it with multimodal, transcript, and frame embeddings — fused with RRF ranking.",
			Tagline:    "Self-hostable semantic search inside long-form video",
			What:       "Semantic search inside long-form video. Ask in plain text and get near-exact moments back — a line in a talk, a slide, a take you need. Transcript, frame, and multimodal embeddings fused with weighted Reciprocal Rank Fusion; collections keep each search space focused.",
			Why:        "Drop-in semantic video search as an internal service — precise timestamps without owning the indexing stack, with workers you can scale to your workload.",
			Technicals: "Async workers (BullMQ + Valkey) with idempotent jobs, pgvector for embeddings, TUS for long uploads, RustFS object storage, and Vercel's AI Gateway for embedding and transcription.",
			Upcoming: []string{
				"Kubernetes and Helm charts for HA deploys",
				"Multimodal queries — image, video, or speech in, not just text",
				"Worker stress-testing toward metric-based autoscaling",
			},
			GitHub:   "https://github.com/Ashutoshbind15/SliceSeeker",
			Packages: "https://github.com/Ashutoshbind15?tab=packages&repo_name=SliceSeeker",
			More: []ProjectLinkExtra{
				{Label: "docs", URL: "https://sliceseeker.ashutoshbind.com"},
			},
			Tech: []string{"React", "TypeScript", "Node.js", "PostgreSQL", "Gemini", "Docker"},
		},
		{
			ID:         "wafercms",
			Name:       "WaferCMS",
			Kind:       "OSS",
			Summary:    "Self-hostable MIT CMS with diagrams and rich text as first-class fields, plus admin agents for collection work.",
			Tagline:    "Self-hostable CMS for complex content",
			What:       "A self-hostable, MIT-licensed CMS where complex content is first-class — rich text and hand-drawn ScribbleSVG diagrams ship as native field types. AI features are opt-in, including an agent that can manipulate the CMS with the help of tool calls.",
			Why:        "I wanted richer content building blocks as native field types, and content that's easier for agents to reach and reshape. Built for this portfolio (blogs coming soon, powered by this CMS itself), then matured into a standalone product.",
			Technicals: "Schema-driven collections with drafts/publish, TipTap rich text + ScribbleSVG diagram fields, Postgres + S3-compatible assets, and an opt-in OpenRouter agent with in-process CMS tools. Docker images on GHCR — self-hosting is a pull and a compose file away.",
			Upcoming: []string{
				"Remote MCP server so agents can connect directly to your CMS",
				"Agent loops for diagram iterations",
			},
			GitHub:   "https://github.com/Ashutoshbind15/WaferCMS",
			Packages: "https://github.com/Ashutoshbind15?tab=packages&repo_name=WaferCMS",
			More: []ProjectLinkExtra{
				{Label: "docs", URL: "https://wafercms.ashutoshbind.com"},
			},
			Tech: []string{"React", "TypeScript", "Node.js", "Docker"},
		},
		{
			ID:         "termchess",
			Name:       "Term Chess",
			Kind:       "Service",
			Summary:    "SSH TUI chess server — multiplayer time controls or solo games against bots.",
			Tagline:    "Chess over SSH, in your terminal",
			What:       "Chess over SSH — real-time multiplayer with time controls, or solo games against rated bots.",
			Why:        "For a deeper understanding of SSH-based client-server and game architectures — Bubble Tea's update loop follows the Elm Architecture — and for fun!",
			Technicals: "A Go server serves the Bubble Tea TUI over Wish/SSH. Shared managers own multiplayer state, clocks, and bot games; Postgres keeps history. Bots call a self-hosted lc0 engine behind a small JSON move API.",
			GitHub:     "https://github.com/Ashutoshbind15/term-chess",
			Site:       "https://termchess.ashutoshbind.com",
			SSH:        "ssh -p 58303 tchess.ashutoshbind.com",
			Tech:       []string{"Go", "PostgreSQL", "Docker"},
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
			Summary: "Product planning, experimentation, and building products and prototypes — from custom frameworks and browser-based site-builders to chess analysis engines, with deep learning around chess, local LLMs, and Electron apps.",
			Highlights: []string{
				"Shipped multiple open-source projects and developer tools.",
				"Developing SaaS products around chess analysis, local LLMs, UI analysis with remote sandboxes etc.",
			},
			Tech: []string{
				"Product Planning",
				"Experimentation",
				"Custom Frameworks",
				"Site Builders",
				"Chess Analysis",
				"Deep Learning",
				"Local LLMs",
				"Electron",
				"Open Source",
			},
		},
		{
			ID:      "cisco-6m",
			Role:    "Internship",
			Org:     "Cisco Systems",
			Period:  "Jan 2025 – Jun 2025",
			Summary: "Worked on a team that works on system health for Cisco Catalyst Center (formerly DNA Center).",
			Highlights: []string{
				"Networking related features in Java",
				"Improved performance and reliability of Unit Tests",
			},
			Tech: []string{
				"Java",
				"Microservices",
				"Unit Testing",
				"Networking",
				"Testing Frameworks",
			},
		},
		{
			ID:      "cisco-2m",
			Role:    "Internship",
			Org:     "Cisco Systems",
			Period:  "May 2024 – Jul 2024",
			Summary: "Worked with the Cisco ISE team.",
			Highlights: []string{
				"Built APIs and dashboards for system metrics",
				"Used Java and Bash; contributed frontend work",
			},
			Tech: []string{"Java", "Bash", "Cisco ISE", "APIs", "Frontend"},
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
		{Name: "Frontend", Items: []string{"React", "Next.js", "TypeScript", "Tailwind CSS", "Vite", "Electron", "shadcn/ui"}},
		{Name: "Backend", Items: []string{"Node.js", "Go", "Java"}},
		{Name: "Databases", Items: []string{"PostgreSQL", "Drizzle"}},
		{Name: "Cloud & Ops", Items: []string{"Docker", "Linux", "Bash"}},
		{Name: "Testing", Items: []string{"Vitest", "JUnit"}},
		{Name: "Workflow", Items: []string{"Git", "GitHub", "pnpm", "ESLint"}},
		{Name: "AI & LLMs", Items: []string{"Gemini", "Ollama"}},
	}
}

func stackToolCount() int {
	n := 0
	for _, cat := range stackCategories() {
		n += len(cat.Items)
	}
	return n
}
