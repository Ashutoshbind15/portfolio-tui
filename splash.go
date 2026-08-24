package main

import (
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/harmonica"
)

const splashFPS = 24

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

	grid := make([][]splashCell, h)
	for y := 0; y < h; y++ {
		grid[y] = make([]splashCell, w)
		for x := 0; x < w; x++ {
			grid[y][x] = waveCell(x, y, w, h, s.t)
		}
	}

	name, bind := splashBanner(w, h)
	p := profile()
	settled := 1 - clamp01(math.Abs(s.namePos)/8)
	shift := int(math.Round(s.namePos))

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
	}, true)

	by := oy + len(name) + gap
	stampBanner(grid, bind, ox+(blockW-bannerWidth(bind))/2, by, func(int, int) string {
		return colorAccent
	}, true)

	my := by + len(bind) + 1
	stampLine(grid, p.Role, ox+(blockW-lipgloss.Width(p.Role))/2, my, colorMuted, false)

	pulse := 0.45 + 0.55*settled*(0.5+0.5*math.Sin(s.t*2.4))
	hint = "  " + hint + "  "
	stampLine(grid, hint, max(0, (w-lipgloss.Width(hint))/2), h-2, mixHex(colorFaint, colorAccent, pulse), false)

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

func stampBanner(grid [][]splashCell, lines []string, ox, oy int, fg func(col, row int) string, bold bool) {
	for i, line := range lines {
		for j, r := range []rune(line) {
			if r == ' ' {
				stampCell(grid, ox+j, oy+i, splashCell{r: ' ', fg: colorFaint})
				continue
			}
			stampCell(grid, ox+j, oy+i, splashCell{r: r, fg: fg(j, i), bold: bold})
		}
	}
}

func stampLine(grid [][]splashCell, text string, ox, oy int, fg string, bold bool) {
	for j, r := range []rune(text) {
		stampCell(grid, ox+j, oy, splashCell{r: r, fg: fg, bold: bold})
	}
}

func stampCell(grid [][]splashCell, x, y int, c splashCell) {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return
	}
	grid[y][x] = c
}

func waveCell(x, y, w, h int, t float64) splashCell {
	if h <= 0 || w <= 0 {
		return splashCell{r: ' '}
	}
	fx := float64(x)
	fy := float64(y)
	nh := float64(h)

	swell := math.Sin(fx*0.13 + t*1.65)
	chop := math.Sin(fx*0.07 - t*0.78 + 1.15)
	ripple := math.Sin(fx*0.28 + t*2.35)
	surface := nh*0.72 + swell*2.2 + chop*1.35 + ripple*0.55
	rel := fy - surface

	switch {
	case rel < -3.2:
		if sparkle(x, y, t) {
			return splashCell{r: '·', fg: colorFaint}
		}
		return splashCell{r: ' '}
	case rel < -1.1:
		if swell+chop > 1.05 {
			return splashCell{r: '~', fg: colorFaint}
		}
		return splashCell{r: ' '}
	case rel < 0:
		return splashCell{r: '~', fg: mixHex(colorMuted, colorAccent, 0.45)}
	case rel < 1.2:
		return splashCell{r: '≈', fg: colorAccent, bold: true}
	case rel < 2.5:
		return splashCell{r: '~', fg: colorAccent}
	case rel < 4.8:
		if math.Sin(fx*0.2+t*1.4+fy*0.55) > -0.2 {
			return splashCell{r: '~', fg: colorMuted}
		}
		return splashCell{r: '·', fg: colorFaint}
	default:
		if math.Sin(fx*0.16+t*0.95+fy*0.4) > 0.15 {
			return splashCell{r: '~', fg: colorFaint}
		}
		if sparkle(x, y+23, t*0.45) {
			return splashCell{r: '.', fg: colorFaint}
		}
		return splashCell{r: ' '}
	}
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
