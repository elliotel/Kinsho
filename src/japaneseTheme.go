package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"os"
)

type japaneseTheme struct{}

var themeVariant = 1

//To make the compiler happy. Remove?
var _ fyne.Theme = (*japaneseTheme)(nil)

func (*japaneseTheme) Font(s fyne.TextStyle) fyne.Resource {
	font := os.Getenv("FYNE_FONT")
	if font == "" {
		return resourceCWINDOWSFontsMEIRYOTTC
	} else {
		return theme.DefaultTheme().Font(s)
	}
}

func (*japaneseTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, (fyne.ThemeVariant)(themeVariant))
}

func (*japaneseTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*japaneseTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func setDarkTheme() {
	themeVariant = 0
}

func setLightTheme() {
	themeVariant = 1
}
