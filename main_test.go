package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestGlyphWidths(t *testing.T) {
	for _, tt := range []struct {
		icon  string
		width int
	}{
		{"𝕏", 1},
		{"↗︎", 1},
		{"＊", 2},
		{"☰", 2},
		{"ッ", 2},
		{"🜂", 1},
		{"✌︎", 1},
		{"☝︎", 1},
	} {
		t.Run(tt.icon, func(t *testing.T) {
			if width := lipgloss.Width(tt.icon); width != tt.width {
				t.Fatalf("width of %q = %d, want %d", tt.icon, width, tt.width)
			}
		})
	}
}

func TestFormScroll(t *testing.T) {
	var glyphs []Glyph
	if err := json.Unmarshal(glyphsJSON, &glyphs); err != nil {
		t.Fatal(err)
	}

	for _, size := range []tea.WindowSizeMsg{
		{Width: 90, Height: 45},
		{Width: 40, Height: 24},
		{Width: 40, Height: 12},
	} {
		t.Run(fmt.Sprintf("%dx%d", size.Width, size.Height), func(t *testing.T) {
			var selected Glyph
			form := newForm(&selected)
			form.Init()
			resizeForm(form, size.Height)
			form.Update(size)
			height := lipgloss.Height(form.View())
			nameColumn := -1

			check := func(glyph Glyph) {
				t.Helper()
				view := ansi.Strip(form.View())
				if h := lipgloss.Height(view); h != height || h > size.Height {
					t.Fatalf("%s: view height = %d, want %d within %d", glyph.Name, h, height, size.Height)
				}
				if !strings.Contains(view, "Glyphs") || !strings.Contains(view, "filter") {
					t.Fatalf("%s: title or help missing:\n%s", glyph.Name, view)
				}
				found := false
				for _, line := range strings.Split(view, "\n") {
					if width := lipgloss.Width(line); width > size.Width {
						t.Fatalf("%s: row width = %d, exceeds %d: %q", glyph.Name, width, size.Width, line)
					}
					if !strings.Contains(line, "> ") {
						continue
					}
					found = true
					if !strings.Contains(line, glyph.Icon) || !strings.HasSuffix(strings.TrimSpace(line), glyph.Name) {
						t.Fatalf("selected row = %q, want %s", line, glyph.Name)
					}
					column := lipgloss.Width(line[:strings.LastIndex(line, glyph.Name)])
					if nameColumn < 0 {
						nameColumn = column
					}
					if column != nameColumn {
						t.Fatalf("%s: name column = %d, want %d", glyph.Name, column, nameColumn)
					}
				}
				if !found {
					t.Fatalf("%s: selected row missing:\n%s", glyph.Name, view)
				}
			}

			check(glyphs[0])
			for _, glyph := range glyphs[1:] {
				form.Update(tea.KeyPressMsg{Code: tea.KeyDown})
				check(glyph)
			}
			for i := len(glyphs) - 2; i >= 0; i-- {
				form.Update(tea.KeyPressMsg{Code: tea.KeyUp})
				check(glyphs[i])
			}
		})
	}
}

func TestFormResize(t *testing.T) {
	var selected Glyph
	form := newForm(&selected)
	form.Init()
	resizeForm(form, 45)
	form.Update(tea.WindowSizeMsg{Width: 90, Height: 45})
	for range 52 {
		form.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}

	for _, size := range []tea.WindowSizeMsg{
		{Width: 40, Height: 12},
		{Width: 90, Height: 45},
		{Width: 40, Height: 8},
		{Width: 90, Height: 45},
	} {
		resizeForm(form, size.Height)
		form.Update(size)
		view := ansi.Strip(form.View())
		if !strings.Contains(view, ">  🜂  Triangle") {
			t.Fatalf("%dx%d: selection lost after resize:\n%s", size.Width, size.Height, view)
		}
		if height := lipgloss.Height(view); height != min(listHeight+4, size.Height) {
			t.Fatalf("%dx%d: unexpected view height %d", size.Width, size.Height, height)
		}
		if !strings.Contains(view, "Glyphs") || !strings.Contains(view, "filter") {
			t.Fatalf("%dx%d: title or help missing:\n%s", size.Width, size.Height, view)
		}
	}
}

func TestFormFilterAndSelect(t *testing.T) {
	var selected Glyph
	form := newForm(&selected)
	form.Init()
	resizeForm(form, 12)
	form.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	for _, r := range "/horizontal" {
		form.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}

	view := ansi.Strip(form.View())
	if !strings.Contains(view, ">  ☰ Horizontal Lines") {
		t.Fatalf("filtered glyph missing or misaligned:\n%s", view)
	}
	if strings.Contains(view, "Twitter") {
		t.Fatalf("unmatched glyph remains visible:\n%s", view)
	}

	form.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if selected != (Glyph{Icon: "☰", Name: "Horizontal Lines"}) {
		t.Fatalf("selected = %+v, want Horizontal Lines", selected)
	}
}
