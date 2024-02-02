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

const iconWidth = 4

func main() {
	var glyphs []Glyph
	var selected Glyph

	json.Unmarshal(glyphsJSON, &glyphs)
	var glyphOptions = make([]huh.Option[Glyph], len(glyphs))

	for i, g := range glyphs {
		icon := g.Icon
		iconWidth := uniseg.StringWidth(icon)
		title := g.Icon + strings.Repeat(" ", 4-iconWidth) + g.Name
		glyphOptions[i] = huh.NewOption(title, g)
	}

	theme := huh.ThemeCharm()
	theme.Focused.Base.Border(lipgloss.Border{})

	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Glyph]().
				Title("Glyphs").
				Options(glyphOptions...).
				Value(&selected).
				Height(20),
		),
	).WithTheme(theme).Run()

	clipboard.WriteAll(selected.Icon)
	fmt.Println(selected.Icon)
}
