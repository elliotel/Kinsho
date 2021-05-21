package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

const (
	internetErrorMsg = "Unable to download JMdict, please check your internet connection"
)

func displayGUI(f fyne.App, inputChan chan string, outputChan chan entry, complete chan struct{}) {
	w := f.NewWindow("「近書」Kinsho")

	f.Settings().SetTheme(&japaneseTheme{})

	lightImage := &canvas.Image{
		File:     "img/jisho_logo_light.png", // file path to image
		FillMode: canvas.ImageFillContain,    // constrains aspect ratio
	}
	darkImage := &canvas.Image{
		File:     "img/jisho_logo_dark.png", // file path to image
		FillMode: canvas.ImageFillContain,   // constrains aspect ratio
	}

	lightImage.SetMinSize(fyne.Size{Width: 700, Height: 150})
	darkImage.SetMinSize(fyne.Size{Width: 700, Height: 150})

	logo := container.NewMax(lightImage)
	acknowledgement := "This publication has included material from the JMdict_e dictionary file in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/"
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

	b1 := widget.NewButton("Toggle Theme", func() {
		if themeVariant != 0 {
			setDarkTheme()
			logo.Remove(lightImage)
			logo.Add(darkImage)
		} else {
			setLightTheme()
			logo.Remove(darkImage)
			logo.Add(lightImage)
		}
		f.Settings().SetTheme(&japaneseTheme{})
		logo.Refresh()
	})

	b2 := widget.NewButton("Update JMdict", func() {
		if !connected() {
			clearContainer(findings)
			findings.Add(container.NewWithoutLayout(widget.NewLabel(internetErrorMsg)))
		} else {
			clearContainer(findings)
			findings.Add(container.NewWithoutLayout(widget.NewLabel("Updating the dictionary, please wait")))
			downloadJMdict()
			decompressAndDeleteGZ(archivePath)
			splitXML()
			clearContainer(findings)
			findings.Add(container.NewWithoutLayout(widget.NewLabel("Update complete!")))
		}
	})

	buttons := container.New(
		layout.NewGridLayoutWithRows(2),
		b1,
		b2,
	)

	searching := false
	searchButton := widget.NewButton("Search", func() {
		go func() {
			if searching {
				return
			}
			searching = true
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
			searching = false
		}()
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
				logo,
				bottomBox,
				nil,
				nil,
			),
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
			bottomBox,
		),
	)

	w.Resize(fyne.Size{Height: 700, Width: 900})

	w.ShowAndRun()
}

func displayConnectionError(f fyne.App) fyne.Window {
	w := f.NewWindow("Connection Error")
	w.SetContent(widget.NewLabel(internetErrorMsg))
	return w
}

func clearContainer(c *fyne.Container) {
	length := len(c.Objects)
	for er := 0; er < length; er++ {
		c.Remove(c.Objects[len(c.Objects)-1])
	}
	canvas.Refresh(c)
}
