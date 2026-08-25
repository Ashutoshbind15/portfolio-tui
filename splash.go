package main

import (
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/harmonica"
)

const splashFPS = 60

type splashFrameMsg time.Time

func tickSplash() tea.Cmd {
	return tea.Tick(time.Second/splashFPS, func(t time.Time) tea.Msg {
		return splashFrameMsg(t)
	})
}

type splashCell struct {
	r    rune
	fg   string
	bold bool
}

type splashModel struct {
	width, height int
	t             float64
	namePos       float64
	nameVel       float64
	spring        harmonica.Spring
}

func newSplashModel() splashModel {
	s := splashModel{
		spring: harmonica.NewSpring(harmonica.FPS(splashFPS), 7.0, 0.55),
	}
	s.Reset(0, 0)
	return s
}

func (s *splashModel) Reset(width, height int) {
	s.width = width
	s.height = height
	s.t = 0
	s.nameVel = 0
	if height > 0 {
		s.namePos = math.Min(10, float64(height)*0.28)
	} else {
		s.namePos = 8
	}
}

func (s splashModel) Init() tea.Cmd { return tickSplash() }

func (s *splashModel) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s splashModel) Update(msg tea.Msg) (splashModel, tea.Cmd) {
	switch msg.(type) {
	case splashFrameMsg:
		s.t += 1.0 / splashFPS
		s.namePos, s.nameVel = s.spring.Update(s.namePos, s.nameVel, 0)
		return s, tickSplash()
	}
	return s, nil
}

func (s splashModel) View() string {
	w, h := s.width, s.height
	if w <= 0 || h <= 0 {
		return ""
	}

	grid := renderWaves(w, h, s.t)

	name, bind := splashBanner(w, h)
	p := profile()
	settled := 1 - clamp01(math.Abs(s.namePos)/8)
	shift := int(math.Round(s.namePos))
	reveal := clamp01((settled - 0.1) / 0.75)

	hint := "enter  ·  open the portfolio"
	gap := 1
	if h < 18 {
		gap = 0
	}
	blockH := len(name) + gap + len(bind) + 2
	blockW := max(bannerWidth(name), bannerWidth(bind), lipgloss.Width(p.Role))
	ox := max(0, (w-blockW)/2)
	oy := max(1, h/8) + shift
	if bottom := oy + blockH; bottom > h-3 {
		oy = max(0, h-3-blockH)
	}

	stampBanner(grid, name, ox+(blockW-bannerWidth(name))/2, oy, func(col, _ int) string {
		return mixHex(colorAccent, colorBright, float64(col)/float64(max(1, bannerWidth(name)-1)))
	}, true, reveal)

	by := oy + len(name) + gap
	stampBanner(grid, bind, ox+(blockW-bannerWidth(bind))/2, by, func(int, int) string {
		return colorAccent
	}, true, clamp01(reveal*1.25-0.25))

	my := by + len(bind) + 1
	stampLine(grid, p.Role, ox+(blockW-lipgloss.Width(p.Role))/2, my, colorMuted, false, clamp01(reveal*1.5-0.5))

	pulse := 0.45 + 0.55*settled*(0.5+0.5*math.Sin(s.t*2.4))
	hint = "  " + hint + "  "
	stampLine(grid, hint, max(0, (w-lipgloss.Width(hint))/2), h-2, mixHex(colorFaint, colorAccent, pulse), false, 1)

	return renderSplash(grid)
}

func splashBanner(width, height int) (name, bind []string) {
	switch {
	case width >= 90 && height >= 18:
		return bannerSlantASHUTOSH, bannerSlantBIND
	case width >= 58 && height >= 16:
		return bannerDoomASHUTOSH, bannerDoomBIND
	case width >= 50 && height >= 12:
		return bannerSmallASHUTOSH, bannerSmallBIND
	default:
		return []string{strings.ToUpper(profile().Name)}, nil
	}
}

func bannerWidth(lines []string) int {
	w := 0
	for _, l := range lines {
		if n := lipgloss.Width(l); n > w {
			w = n
		}
	}
	return w
}

func stampBanner(grid [][]splashCell, lines []string, ox, oy int, fg func(col, row int) string, bold bool, reveal float64) {
	bw := max(1, bannerWidth(lines)-1)
	for i, line := range lines {
		for j, r := range []rune(line) {
			alpha := clamp01((reveal - float64(j)/float64(bw)*0.85) * 5)
			if alpha <= 0 {
				continue
			}
			if r == ' ' {
				stampCell(grid, ox+j, oy+i, splashCell{r: ' ', fg: colorFaint})
				continue
			}
			stampCell(grid, ox+j, oy+i, splashCell{r: r, fg: mixHex(colorBg, fg(j, i), alpha), bold: bold})
		}
	}
}

func stampLine(grid [][]splashCell, text string, ox, oy int, fg string, bold bool, reveal float64) {
	tw := max(1, lipgloss.Width(text)-1)
	for j, r := range []rune(text) {
		alpha := clamp01((reveal - float64(j)/float64(tw)*0.8) * 5)
		if alpha <= 0 {
			continue
		}
		stampCell(grid, ox+j, oy, splashCell{r: r, fg: mixHex(colorBg, fg, alpha), bold: bold})
	}
}

func stampCell(grid [][]splashCell, x, y int, c splashCell) {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return
	}
	grid[y][x] = c
}

// Each terminal cell is a 2×4 braille pixel, matching the site's dotted hero.
var brailleDot = [4][2]byte{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

func renderWaves(w, h int, t float64) [][]splashCell {
	grid := make([][]splashCell, h)
	for y := 0; y < h; y++ {
		grid[y] = make([]splashCell, w)
		for x := 0; x < w; x++ {
			grid[y][x] = splashCell{r: ' '}
		}
	}
	if w <= 0 || h <= 0 {
		return grid
	}

	pw, ph := w*2, h*4
	bits := make([]byte, w*h)
	crest := make([]bool, w*h)
	surf := make([]float64, pw)
	for px := 0; px < pw; px++ {
		surf[px] = waveSurface(px, ph, t)
	}

	plot := func(px, py int, isCrest bool) {
		if px < 0 || py < 0 || px >= pw || py >= ph {
			return
		}
		cx, cy := px/2, py/4
		idx := cy*w + cx
		bits[idx] |= brailleDot[py%4][px%2]
		if isCrest {
			crest[idx] = true
		}
	}

	const ribbon = 18
	for px := 0; px < pw; px++ {
		yi := int(math.Round(surf[px]))
		if px+1 < pw {
			y2 := int(math.Round(surf[px+1]))
			step := 1
			if y2 < yi {
				step = -1
			}
			for py := yi; py != y2; py += step {
				plot(px, py, true)
			}
		}
		bot := min(ph, yi+ribbon)
		for py := yi - 1; py < bot; py++ {
			depth := py - yi
			if depth <= 1 {
				plot(px, py, true)
				continue
			}
			thresh := float64(depth) / float64(ribbon)
			n := math.Sin(float64(px)*12.9898 + float64(py)*78.233 + t*0.2)
			n = n * 43758.5453
			n -= math.Floor(n)
			if n > thresh*0.55 {
				plot(px, py, false)
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			b := bits[idx]
			if b == 0 {
				if sparkle(x, y, t) {
					grid[y][x] = splashCell{r: '·', fg: colorFaint}
				}
				continue
			}
			sx := min(x*2, pw-1)
			rel := (float64(y) + 0.5) - surf[sx]/4
			cell := splashCell{r: rune(0x2800 + int(b))}
			switch {
			case crest[idx] || rel < 0.55:
				cell.fg = colorAccent
				cell.bold = true
			case rel < 2.0:
				cell.fg = mixHex(colorAccent, colorMuted, 0.3)
			case rel < 4.2:
				cell.fg = colorMuted
			default:
				cell.fg = colorFaint
			}
			grid[y][x] = cell
		}
	}
	return grid
}

func waveSurface(px, ph int, t float64) float64 {
	fx := float64(px)
	swell := math.Sin(fx*0.062 + t*1.55)
	chop := math.Sin(fx*0.034 - t*0.82 + 1.2)
	ripple := math.Sin(fx*0.148 + t*2.35)
	return float64(ph)*0.80 + swell*7.5 + chop*4.2 + ripple*1.8
}

func sparkle(x, y int, t float64) bool {
	n := math.Sin(float64(x)*12.9898 + float64(y)*78.233 + t*0.18)
	n = n * 43758.5453
	n -= math.Floor(n)
	return n > 0.985
}

func renderSplash(grid [][]splashCell) string {
	if len(grid) == 0 {
		return ""
	}
	rows := make([]string, len(grid))
	for y, line := range grid {
		var b strings.Builder
		var run []rune
		var fg string
		var bold bool
		flush := func() {
			if len(run) == 0 {
				return
			}
			st := lipgloss.NewStyle().
				Foreground(lipgloss.Color(fg)).
				Background(lipgloss.Color(colorBg))
			if bold {
				st = st.Bold(true)
			}
			b.WriteString(st.Render(string(run)))
			run = run[:0]
		}
		for _, c := range line {
			r := c.r
			if r == 0 {
				r = ' '
			}
			cellFG := c.fg
			if cellFG == "" {
				cellFG = colorFaint
			}
			if len(run) > 0 && (cellFG != fg || c.bold != bold) {
				flush()
			}
			if len(run) == 0 {
				fg = cellFG
				bold = c.bold
			}
			run = append(run, r)
		}
		flush()
		rows[y] = b.String()
	}
	return strings.Join(rows, "\n")
}

func mixHex(a, b string, t float64) string {
	t = clamp01(t)
	ca, cb := parseHexColor(a), parseHexColor(b)
	if ca == nil || cb == nil {
		if t < 0.5 {
			return a
		}
		return b
	}
	ar, ag, ab, _ := ca.RGBA()
	br, bg, bb, _ := cb.RGBA()
	mix := func(x, y uint32) uint8 {
		xf := float64(x) / 257
		yf := float64(y) / 257
		return uint8(xf + (yf-xf)*t)
	}
	return sprintfRGB(mix(ar, br), mix(ag, bg), mix(ab, bb))
}

func sprintfRGB(r, g, b uint8) string {
	const hex = "0123456789ABCDEF"
	out := [7]byte{'#', 0, 0, 0, 0, 0, 0}
	out[1] = hex[r>>4]
	out[2] = hex[r&0x0f]
	out[3] = hex[g>>4]
	out[4] = hex[g&0x0f]
	out[5] = hex[b>>4]
	out[6] = hex[b&0x0f]
	return string(out[:])
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

var bannerDoomASHUTOSH = []string{
	"  ___   _____  _   _  _   _  _____  _____  _____  _   _",
	" / _ \\ /  ___|| | | || | | ||_   _||  _  |/  ___|| | | |",
	"/ /_\\ \\\\ `--. | |_| || | | |  | |  | | | |\\ `--. | |_| |",
	"|  _  | `--. \\|  _  || | | |  | |  | | | | `--. \\|  _  |",
	"| | | |/\\__/ /| | | || |_| |  | |  \\ \\_/ //\\__/ /| | | |",
	"\\_| |_/\\____/ \\_| |_/ \\___/   \\_/   \\___/ \\____/ \\_| |_/",
}

var bannerDoomBIND = []string{
	"______  _____  _   _ ______",
	"| ___ \\|_   _|| \\ | ||  _  \\",
	"| |_/ /  | |  |  \\| || | | |",
	"| ___ \\  | |  | . ` || | | |",
	"| |_/ / _| |_ | |\\  || |/ /",
	"\\____/  \\___/ \\_| \\_/|___/",
}

var bannerSmallASHUTOSH = []string{
	"    _     ___   _  _   _   _   _____    ___    ___   _  _",
	"   /_\\   / __| | || | | | | | |_   _|  / _ \\  / __| | || |",
	"  / _ \\  \\__ \\ | __ | | |_| |   | |   | (_) | \\__ \\ | __ |",
	" /_/ \\_\\ |___/ |_||_|  \\___/    |_|    \\___/  |___/ |_||_|",
}

var bannerSmallBIND = []string{
	"  ___   ___   _  _   ___",
	" | _ ) |_ _| | \\| | |   \\",
	" | _ \\  | |  | .` | | |) |",
	" |___/ |___| |_|\\_| |___/",
}

var bannerSlantASHUTOSH = []string{
	"    ___    _____    __  __   __  __  ______   ____    _____    __  __",
	"   /   |  / ___/   / / / /  / / / / /_  __/  / __ \\  / ___/   / / / /",
	"  / /| |  \\__ \\   / /_/ /  / / / /   / /    / / / /  \\__ \\   / /_/ /",
	" / ___ | ___/ /  / __  /  / /_/ /   / /    / /_/ /  ___/ /  / __  /",
	"/_/  |_|/____/  /_/ /_/   \\____/   /_/     \\____/  /____/  /_/ /_/",
}

var bannerSlantBIND = []string{
	"    ____     ____    _   __    ____",
	"   / __ )   /  _/   / | / /   / __ \\",
	"  / __  |   / /    /  |/ /   / / / /",
	" / /_/ /  _/ /    / /|  /   / /_/ /",
	"/_____/  /___/   /_/ |_/   /_____/",
}
