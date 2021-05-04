package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

func displayGUI(inputChan chan string, outputChan chan entry, complete chan struct{}) {
	f := app.New()

	w := f.NewWindow("Our「辞書」Dictionary")

	f.Settings().SetTheme(&japaneseTheme{})

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
		log.Fatal(err)
	}

	b1 := widget.NewButton("Toggle Theme", func() {
		if themeVariant != 0 {
			setDarkTheme()
			logo.SetResource(darkResource)
		} else {
			setLightTheme()
			logo.SetResource(lightResource)
		}
		f.Settings().SetTheme(&japaneseTheme{})
		logo.Refresh()
	})

	b2 := widget.NewButton("Update JMdict", func() {
		downloadJMdict()
		decompressAndDeleteGZ(archiveName)
	})

	buttons := container.New(
		layout.NewGridLayoutWithRows(2),
		b1,
		b2,
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("search here")

	results := widget.NewLabel("results")

	searchButton := widget.NewButton("Search", func() {
		go parseDoc(inputChan, outputChan, complete)
		inputChan <- strings.ToLower(input.Text)
		output := ""
		found := false
		finished := false

		for !finished {
			select {
			case response := <-outputChan:
				found = true
				output += response.kanji + "\n" + response.kana + "\n" + response.def + "\n\n"
				results.SetText(output)
				canvas.Refresh(results)
			case <-complete:
				if !found {
					results.SetText("No results found for \"" + input.Text + "\"")
				}
				finished = true
			}
		}
	})

	search := container.New(layout.NewBorderLayout(nil, nil, nil, searchButton), searchButton, input)

	findings := container.NewVScroll(results)
	searchAndResult := container.New(layout.NewBorderLayout(
		search,
		nil,
		nil,
		nil,
		), search, findings)

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

	w.Resize(fyne.Size{Height: 360, Width: 640})

	w.ShowAndRun()
}
