package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"strings"
)

func displayGUI(inputChan chan string, outputChan chan entry, complete chan struct{}) {
	f := app.New()

	w := f.NewWindow("Our「辞書」Dictionary")

	f.Settings().SetTheme(&japaneseTheme{})

	lightResource, err := fyne.LoadResourceFromPath("jisho_logo_light.png")
	darkResource, err := fyne.LoadResourceFromPath("jisho_logo_dark.png")
	logo := widget.NewIcon(lightResource)
	acknowledgement := "This publication has included material from the JMdict dictionary file in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/"
	bottomText := widget.NewLabel(acknowledgement)
	bottomText.Wrapping = fyne.TextWrapWord
	bottomText.Alignment = fyne.TextAlignCenter
	bottomBox := container.New(
		layout.NewMaxLayout(),
		bottomText,
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("search here")

	allResults := make([]fyne.CanvasObject, entryAmount)
	findings := container.NewVBox()
	findingsScroll := container.NewVScroll(findings)

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
		clearContainer(findings)
		findings.Add(container.NewWithoutLayout(widget.NewLabel("Updating the dictionary, please wait")))
		downloadJMdict()
		decompressAndDeleteGZ(archiveName)
		clearContainer(findings)
		findings.Add(container.NewWithoutLayout(widget.NewLabel("Update complete!")))
	})

	buttons := container.New(
		layout.NewGridLayoutWithRows(2),
		b1,
		b2,
	)

	searchButton := widget.NewButton("Search", func() {
		clearContainer(findings)
		go parseDoc(inputChan, outputChan, complete)
		inputChan <- strings.ToLower(input.Text)
		found := false
		finished := false
		i := 0
		for !finished {
			select {
			case response := <-outputChan:
				found = true
				var result string
				for i, r := range response.kanji {
					if i > 0 {
						result += "  ·  "
					}
					result += r
				}
				result += "\n"
				for i, r := range response.kana {
					if i > 0 {
						result += "  |  "
					}
					result += r
				}
				for i, r := range response.def {
					result += "\n" + strconv.Itoa(i+1) + ". " + r
				}
				labResult := widget.NewLabel(result)
				labResult.Wrapping = fyne.TextWrapWord
				allResults[i] = container.NewWithoutLayout(labResult)
				i++
			case <-complete:
				if !found {
					findings.Add(widget.NewLabel("No results found for \"" + input.Text + "\""))
				}
				finished = true
			}
		}

		for j := 0; j < i; j++ {
			findings.Add(allResults[j])
			findings.Refresh()
		}
		findingsScroll.Refresh()
	})

	search := container.New(layout.NewBorderLayout(nil, nil, nil, searchButton), searchButton, input)

	searchAndResult := container.New(layout.NewBorderLayout(
		search,
		nil,
		nil,
		nil,
	), search, findingsScroll)

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

	w.Resize(fyne.Size{Height: 700, Width: 900})

	w.ShowAndRun()
}

func clearContainer(c *fyne.Container) {
	length := len(c.Objects)
	for er := 0; er < length; er++ {
		c.Remove(c.Objects[len(c.Objects)-1])
	}
	canvas.Refresh(c)
}
