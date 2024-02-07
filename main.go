package main

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/rivo/uniseg"
)

//go:embed glyphs.json
var glyphsJSON []byte

type Glyph struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

const iconWidth = 3

func main() {
	var glyphs []Glyph
	var selected Glyph

	_ = json.Unmarshal(glyphsJSON, &glyphs)
	glyphOptions := make([]huh.Option[Glyph], len(glyphs))

	for i, g := range glyphs {
		title := " " + g.Icon + strings.Repeat(" ", iconWidth-uniseg.StringWidth(g.Icon)) + g.Name
		glyphOptions[i] = huh.NewOption(title, g)
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base.Border(lipgloss.HiddenBorder())

	_ = huh.NewForm(
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

	termenv.Copy(selected.Icon)
	fmt.Println(selected.Icon)
}
