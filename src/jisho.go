
/*
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
)

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("Box Layout")

	text1 := canvas.NewText("Hello", color.Black)
	text2 := canvas.NewText("There", color.Black)
	text3 := canvas.NewText("(right)", color.Black)
	content := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)
	test := canvas.NewRectangle(color.RGBA{255, 0, 255, 255})
	test.SetMinSize(fyne.NewSize(1024, 55))
	image := canvas.NewImageFromFile("logo.png")
	image.SetMinSize(fyne.Size{Width: 150, Height: 150})
	image.FillMode = canvas.ImageFillContain
	logo := container.New(layout.NewBorderLayout(nil,nil,nil,nil), image)
	logoContainer := container.New(layout.NewVBoxLayout(), content, layout.NewSpacer(), logo, layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer())

	//text4 := canvas.NewText("centered", color.Black)
	//centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())
	myWindow.SetContent(logoContainer)
	myWindow.Resize(fyne.NewSize(50,50))
	myWindow.ShowAndRun()



}
*/

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
)

type image struct {
	url string
}

func main() {
	//image := canvas.NewImageFromFile("jisho_logo.png")
	//image.SetMinSize(fyne.Size{Width: 150, Height: 150})
	//image.FillMode = canvas.ImageFillContain

	f := app.New()
	w := f.NewWindow("")

	b1 := widget.NewButton("Placeholder1", func() { /*Do something*/ })

	b2 := widget.NewButton("Placeholder2", func() { /*Do something else*/ })

	resource, err := fyne.LoadResourceFromPath("jisho_logo.png")

	if err != nil {
		//Error handling
		log.Println(err)
	}

	i1 := widget.NewIcon(resource)
	i1.ExtendBaseWidget(i1)
		string := "This publication has included material from the EDICT and KANJIDIC dictionary files in accordance with the licence provisions of the Electronic Dictionaries Research Group. See http://www.edrdg.org/1"
		bottomText := widget.NewLabel(string)
	bottomText.Wrapping = fyne.TextWrapWord
	bottomText.Alignment = fyne.TextAlignCenter
	bottomBox := container.New(
		layout.NewMaxLayout(),
		bottomText,
	)

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
			i1,
				container.New(
			layout.NewGridLayout(5),

			container.New(
				layout.NewGridLayoutWithRows(2),
				b1,
				b2,
				),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
			layout.NewSpacer(),
		),
	),
	bottomBox,
		),
	)


	w.Resize(fyne.Size{Height: 320, Width: 480})

	w.ShowAndRun()
}