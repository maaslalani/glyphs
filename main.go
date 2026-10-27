package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/termenv"
)

//go:embed glyphs.json
var glyphsJSON []byte

type Glyph struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

const (
	iconWidth  = 3
	listHeight = 20
)

func main() {
	var selected Glyph

	if err := newForm(&selected).Run(); err != nil {
		return
	}

	termenv.Copy(selected.Icon)
	fmt.Println(selected.Icon)
}

func newForm(selected *Glyph) *huh.Form {
	var glyphs []Glyph

	_ = json.Unmarshal(glyphsJSON, &glyphs)
	glyphOptions := make([]huh.Option[Glyph], len(glyphs))

	for i, g := range glyphs {
		title := " " + g.Icon + strings.Repeat(" ", iconWidth-lipgloss.Width(g.Icon)) + g.Name
		glyphOptions[i] = huh.NewOption(title, g)
	}

	theme := huh.ThemeFunc(func(isDark bool) *huh.Styles {
		theme := huh.ThemeCharm(isDark)
		theme.Focused.Base = theme.Focused.Base.Border(lipgloss.HiddenBorder())
		return theme
	})

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Glyph]().
				Title("Glyphs").
				Description(" ").
				Options(glyphOptions...).
				Value(selected).
				Height(listHeight),
		),
	).WithTheme(theme)

	return form.WithProgramOptions(
		tea.WithOutput(os.Stderr),
		tea.WithFilter(func(_ tea.Model, msg tea.Msg) tea.Msg {
			if size, ok := msg.(tea.WindowSizeMsg); ok {
				resizeForm(form, size.Height)
			}
			return msg
		}),
	)
}

func resizeForm(form *huh.Form, height int) {
	height = min(listHeight, max(3, height-4))
	form.GetFocusedField().WithHeight(height)
	form.WithHeight(height + 3)
}
