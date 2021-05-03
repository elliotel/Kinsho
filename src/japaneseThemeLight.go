package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type lightTheme struct{}

//To make the compiler happy. Remove?
var _ fyne.Theme = (*lightTheme)(nil)

func (*lightTheme) Font(s fyne.TextStyle) fyne.Resource {
	return resourceCWindowsFontsUDDIGIKYOKASHONRTTC
}

func (*lightTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return theme.LightTheme().Color(n, 1)
}

func (*lightTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(n)
}

func (*lightTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(n)
}