package main

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/rivo/uniseg"
)

//go:embed glyphs.json
var glyphsJSON []byte

type Glyph struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

const iconWidth = 3

var iconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).MarginLeft(1)
var nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))

func main() {
	var glyphs []Glyph
	var selected Glyph

	json.Unmarshal(glyphsJSON, &glyphs)
	var glyphOptions = make([]huh.Option[Glyph], len(glyphs))

	for i, g := range glyphs {
		title := iconStyle.Render(g.Icon) + strings.Repeat(" ", iconWidth-uniseg.StringWidth(g.Icon)) + nameStyle.Render(g.Name)
		glyphOptions[i] = huh.NewOption(title, g)
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base.Border(lipgloss.HiddenBorder())

	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Glyph]().
				Title("Glyphs").
				Description(" ").
				Options(glyphOptions...).
				Value(&selected).
				Height(20),
			huh.NewNote(),
		),
	).WithTheme(theme).Run()

	clipboard.WriteAll(selected.Icon)
	fmt.Println(selected.Icon)
}
