package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"os"
)

type darkTheme struct{}

//To make the compiler happy. Remove?
var _ fyne.Theme = (*darkTheme)(nil)

func (*darkTheme) Font(s fyne.TextStyle) fyne.Resource {
	font := os.Getenv("FYNE_FONT")
	if font == "" {
		return resourceCWINDOWSFontsMEIRYOTTC
	} else {
		return theme.DefaultTheme().Font(s)
	}
}

func (*darkTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return theme.DarkTheme().Color(n, 0)
}

func (*darkTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(n)
}

func (*darkTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(n)
}