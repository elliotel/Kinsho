package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
)

/*
type image struct {
	url string
}

 */




func main() {
	//image := canvas.NewImageFromFile("jisho_logo.png")
	//image.SetMinSize(fyne.Size{Width: 150, Height: 150})
	//image.FillMode = canvas.ImageFillContain

	f := app.New()

	w := f.NewWindow("")

	f.Settings().SetTheme(&lightTheme{})

	lightResource, err := fyne.LoadResourceFromPath("jisho_logo_light.png")
	darkResource, err := fyne.LoadResourceFromPath("jisho_logo_dark.png")
	logo := widget.NewIcon(lightResource)
	string := "This publication has included material from the EDICT and KANJIDIC dictionary files in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/"
	bottomText := widget.NewLabel(string)
	bottomText.Wrapping = fyne.TextWrapWord
	bottomText.Alignment = fyne.TextAlignCenter
	bottomBox := container.New(
		layout.NewMaxLayout(),
		bottomText,
	)

	if err != nil {
		//Error handling
		log.Println(err)
	}

	darkThemeOn := false
	b1 := widget.NewButton("Toggle Theme", func() {
		if !darkThemeOn {
			f.Settings().SetTheme(&darkTheme{})
			darkThemeOn = true
			logo.SetResource(darkResource)
		} else {
			f.Settings().SetTheme(&lightTheme{})
			darkThemeOn = false
			logo.SetResource(lightResource)
		}
		logo.Refresh()
	})

	b2 := widget.NewButton("Placeholder2", func() { /*Do something*/ })
	
	buttons := container.New(
		layout.NewGridLayoutWithRows(2),
		b1,
		b2,
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("search here")

	results := widget.NewLabel("results")

	searchButton := widget.NewButton("Search", func() {
		results.SetText("Results for " + input.Text)
		canvas.Refresh(results)
	})

	search := container.New(layout.NewBorderLayout(nil,nil,nil,searchButton), searchButton, input)

	findings := container.New(layout.NewHBoxLayout(),results)
	searchAndResult := container.New(layout.NewVBoxLayout(), search, findings)

	w.SetContent(
		container.New(
			layout.NewBorderLayout(
				nil,
				bottomBox,
				nil,
				nil,
				),
			container.New(

			layout.NewGridLayout(1),
			logo,
				container.New(
			layout.NewBorderLayout(
				nil,
				nil,
				buttons,
				nil,
				),
				buttons,
			searchAndResult,
		),
	),
	bottomBox,
		),
	)


	w.Resize(fyne.Size{Height: 320, Width: 480})

	w.ShowAndRun()
}