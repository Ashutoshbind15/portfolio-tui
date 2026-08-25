package main

// Static portfolio content, kept in sync with portfolio-new.
// Home and the project list use the site's card copy (tagline + summary).
// Opening a project uses the individual page text (what / why / where /
// technicals / upcoming) and the same Where-chip labels, without the gallery.

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
		LinkedIn: "in.linkedin.com/in/ashutosh-bind-56806b22b",
		SSH:      "go run github.com/Ashutoshbind15/portfolio-tui@main",
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

// ProjectLink is a Where-section chip from the site (label + href/value).
type ProjectLink struct {
	Label string
	Value string
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
	Links      []ProjectLink
}

func projects() []Project {
	return []Project{
		{
			ID:         "varypane",
			Name:       "VaryPane",
			Kind:       "SaaS",
			Summary:    "Cloud-hosted parallel UI agents with starter templates, prompt libraries, and streaming previews.",
			Tagline:    "Parallel UI-generation agents in the cloud",
			What:       "Parallel UI-generation agents in the cloud. Pick a starter template, keep a reusable prompt library, attach skills, and compare models side by side — results stream in real time. Runs pi-coding-agent inside every sandbox.",
			Why:        "Visual comparison — and picking the best bits out of generated UIs — is usually very helpful, especially for developers doing backend work, people new to design, or anyone who just doesn't have the time for a skill this demanding. The cloud version scales and improves that workflow for hosting and for end users.",
			Technicals: "BYOK for OpenRouter and Upstash Box. Parallel agent sandboxes per session, a dedicated preview gateway, labelled box ownership, and self-hosted secret storage via Infisical.",
			Upcoming: []string{
				"Screenshot-to-agent input — iterate from a wireframe or your current design",
				"GitHub repo connections and agent-loop iterations on one-shot designs",
				"Scripted environment setup",
			},
			Links: []ProjectLink{
				{Label: "site", Value: "https://varypane.com"},
			},
		},
		{
			ID:         "scribblesvg",
			Name:       "ScribbleSVG",
			Kind:       "OSS",
			Summary:    "TypeScript core and React utilities for hand-drawn SVG diagrams on an infinite canvas.",
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
			Links: []ProjectLink{
				{Label: "GitHub", Value: "https://github.com/Ashutoshbind15/ScribbleSVG"},
				{Label: "core", Value: "https://www.npmjs.com/package/@scribblesvg/core"},
				{Label: "react-utils", Value: "https://www.npmjs.com/package/@scribblesvg/react-utils"},
				{Label: "Playground", Value: "https://scribblesvg.ashutoshbind.com"},
				{Label: "Docs", Value: "https://scribblesvg-docs.ashutoshbind.com"},
			},
		},
		{
			ID:         "sliceseeker",
			Name:       "SliceSeeker",
			Kind:       "OSS",
			Summary:    "Ask in plain text, get back exact moments from long-form video via fused transcript and frame embeddings.",
			Tagline:    "Self-hostable semantic search inside long-form video",
			What:       "Semantic search inside long-form video. Ask in plain text and get near-exact moments back — a line in a talk, a slide, a take you need. Transcript, frame, and multimodal embeddings fused with weighted Reciprocal Rank Fusion; collections keep each search space focused.",
			Why:        "Drop-in semantic video search as an internal service — precise timestamps without owning the indexing stack, with workers you can scale to your workload.",
			Technicals: "Async workers (BullMQ + Valkey) with idempotent jobs, pgvector for embeddings, TUS for long uploads, RustFS object storage, and Vercel's AI Gateway for embedding and transcription.",
			Upcoming: []string{
				"Kubernetes and Helm charts for HA deploys",
				"Multimodal queries — image, video, or speech in, not just text",
				"Worker stress-testing toward metric-based autoscaling",
			},
			Links: []ProjectLink{
				{Label: "GitHub", Value: "https://github.com/Ashutoshbind15/SliceSeeker"},
				{Label: "GHCR images", Value: "https://github.com/Ashutoshbind15?tab=packages&repo_name=SliceSeeker"},
				{Label: "Docs", Value: "https://sliceseeker.ashutoshbind.com"},
			},
		},
		{
			ID:         "termchess",
			Name:       "TermChess",
			Kind:       "Service",
			Summary:    "Real-time multiplayer chess with time controls, or solo games against rated bots: all over SSH.",
			Tagline:    "Chess over SSH, in your terminal",
			What:       "Chess over SSH — real-time multiplayer with time controls, or solo games against rated bots.",
			Why:        "For a deeper understanding of SSH-based client-server and game architectures — Bubble Tea's update loop follows the Elm Architecture — and for fun!",
			Technicals: "A Go server serves the Bubble Tea TUI over Wish/SSH. Shared managers own multiplayer state, clocks, and bot games; Postgres keeps history. Bots call a self-hosted lc0 engine behind a small JSON move API.",
			Links: []ProjectLink{
				{Label: "ssh", Value: "ssh -p 58303 tchess.ashutoshbind.com"},
				{Label: "Marketing site", Value: "https://termchess.ashutoshbind.com"},
				{Label: "GitHub", Value: "https://github.com/Ashutoshbind15/term-chess"},
			},
		},
		{
			ID:         "wafercms",
			Name:       "WaferCMS",
			Kind:       "OSS",
			Summary:    "Schema-driven CMS with rich text, diagram fields, and opt-in AI agents: self-hostable.",
			Tagline:    "Self-hostable CMS for complex content",
			What:       "A self-hostable, MIT-licensed CMS where complex content is first-class — rich text and hand-drawn ScribbleSVG diagrams ship as native field types. AI features are opt-in, including an agent that can manipulate the CMS with the help of tool calls.",
			Why:        "I wanted richer content building blocks as native field types, and content that's easier for agents to reach and reshape. Built for this portfolio (blogs coming soon, powered by this CMS itself), then matured into a standalone product.",
			Technicals: "Schema-driven collections with drafts/publish, TipTap rich text + ScribbleSVG diagram fields, Postgres + S3-compatible assets, and an opt-in OpenRouter agent with in-process CMS tools. Docker images on GHCR — self-hosting is a pull and a compose file away.",
			Upcoming: []string{
				"Remote MCP server so agents can connect directly to your CMS",
				"Agent loops for diagram iterations",
			},
			Links: []ProjectLink{
				{Label: "GitHub", Value: "https://github.com/Ashutoshbind15/WaferCMS"},
				{Label: "GHCR images", Value: "https://github.com/Ashutoshbind15?tab=packages&repo_name=WaferCMS"},
				{Label: "Docs", Value: "https://wafercms.ashutoshbind.com"},
			},
		},
		{
			ID:         "aiuicomparator",
			Name:       "AI UI Comparator",
			Kind:       "OSS",
			Summary:    "Compare how AI models generate UI locally: isolated sandboxes, reusable prompts, side-by-side results.",
			Tagline:    "Local side-by-side AI UI generation",
			What:       "A local tool for side-by-side comparison of how AI models generate UI. A custom pi-agent runs a tight look → read → write/edit loop with four tools, and each worker is isolated by Anthropic's sandbox-runtime.",
			Why:        "Comparing UI-capable models meant juggling prompts, starters, and sandboxes by hand. This is the local answer — reusable prompts, isolated workers, side-by-side results, all on your machine.",
			Technicals: "Installable on Linux with a single script. Express API + React UI; Postgres for batches, prompts, and run state. Each session gets an isolated workdir from a starter template with network isolation except OpenRouter for model calls; agent events stream over SSE; completed runs launch managed Vite previews with a port registry that can reap orphans.",
			Upcoming: []string{
				"Replace the Postgres container with SQLite to simplify operations",
			},
			Links: []ProjectLink{
				{Label: "GitHub", Value: "https://github.com/Ashutoshbind15/ai-ui-comparator"},
			},
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
				"SaaS: VaryPane, and more in the works — chess analysis and improvement app, Electron app for local models, component library for content software.",
				"Open source: ScribbleSVG, SliceSeeker, TermChess, WaferCMS, and AI UI Comparator.",
			},
			Tech: []string{
				"TypeScript",
				"React",
				"Next.js",
				"Go",
				"PostgreSQL",
				"Docker",
				"DevOps",
				"AI Agents",
				"Local LLMs",
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
