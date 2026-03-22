package tui

import (
	"slices"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/is386/gcv/internal/git"
)

var daysOfWeek = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

type Model struct {
	heatmaps    []string
	years       []int
	totals      []int
	currentPage int
}

func NewModel(contributionsMap git.ContributionsMap) Model {
	m := Model{currentPage: 0}
	for year := range contributionsMap {
		m.years = append(m.years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(m.years)))
	m.buildHeatMaps(contributionsMap)
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "j", "h", "down", "left":
			m.currentPage++
			if m.currentPage > len(m.years)-1 {
				m.currentPage = 0
			}
			return m, nil
		case "k", "l", "up", "right":
			m.currentPage--
			if m.currentPage < 0 {
				m.currentPage = len(m.years) - 1
			}
			return m, nil
		default:
			return m, nil
		}
	case tea.WindowSizeMsg:
		return m, nil
	}
	return m, cmd
}

func (m Model) View() tea.View {
	heatmap := m.heatmaps[m.currentPage]
	year := m.years[m.currentPage]
	total := m.totals[m.currentPage]
	title := lipgloss.NewStyle().Bold(true).Render(lipgloss.Sprintf("%d contributions in %d", total, year))
	help := lipgloss.NewStyle().Foreground(lipgloss.Color(("240"))).Render("← prev • → next • q quit")
	v := tea.NewView(lipgloss.NewStyle().Render(lipgloss.Sprintf("%s\n\n%s\n\n%s", title, heatmap, help)))
	v.AltScreen = true
	return v
}

func (m *Model) buildHeatMaps(contributionsMap git.ContributionsMap) {
	for _, year := range m.years {
		total := 0
		heatmapRows := make([]string, 7)
		contribs := contributionsMap[year]
		maxVal := slices.Max(contribs[:])
		firstDayOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Weekday()
		isLeapYear := (year%4 == 0 && year%100 != 0) || year%400 == 0

		for i := range heatmapRows {
			heatmapRows[i] += lipgloss.Sprintf("%s ", daysOfWeek[i])
			if i < int(firstDayOfYear) {
				heatmapRows[i] += "  "
			}
		}

		rowNum := firstDayOfYear

		for i := 1; i < len(contribs); i++ {
			if i == 366 && !isLeapYear {
				break
			}

			c := contribs[i]
			var color string
			if c == 0 {
				color = "236"
			} else if c < maxVal/3 {
				color = "22"
			} else if c >= maxVal/3 && c < 2*(maxVal/3) {
				color = "34"
			} else {
				color = "34"
			}

			heatmapRows[rowNum] += lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("■ ")
			total += c

			rowNum++
			if rowNum > 6 {
				rowNum = 0
			}
		}
		m.heatmaps = append(m.heatmaps, strings.Join(heatmapRows, "\n"))
		m.totals = append(m.totals, total)
	}
}
