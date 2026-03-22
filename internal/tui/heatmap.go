package tui

import (
	"slices"
	"sort"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/is386/git-contributions-visualizer/internal/git"
)

type Model struct {
	heatmaps    []string
	years       []int
	currentPage int
}

func NewModel(contributionsMap git.ContributionsMap) Model {
	m := Model{currentPage: 0}
	for year := range contributionsMap {
		m.years = append(m.years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(m.years)))

	for _, year := range m.years {
		contribs := contributionsMap[year]
		var heatmap strings.Builder
		maxVal := slices.Max(contribs[:])
		rowNum := 0

		for i, c := range contribs {
			if rowNum < i/52 {
				rowNum = i / 52
				heatmap.WriteString("\n")
			}
			if rowNum > 6 {
				break
			}

			var color string
			if c == 0 {
				color = "240"
			} else if c < maxVal/3 {
				color = "22"
			} else if c >= maxVal/3 && c < 2*(maxVal/3) {
				color = "34"
			} else {
				color = "34"
			}

			heatmap.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("■ "))
		}

		m.heatmaps = append(m.heatmaps, heatmap.String())
	}
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

// TODO: nav bar with years on left
// TODO: 210 contributions in 2026
// TODO: centered on screen with box around
// TODO: add help bar
func (m Model) View() tea.View {
	heatmap := m.heatmaps[m.currentPage]
	year := m.years[m.currentPage]
	v := tea.NewView(lipgloss.NewStyle().Render(strconv.Itoa(year) + "\n" + heatmap))
	v.AltScreen = true
	return v
}
