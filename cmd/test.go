package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/flopp/go-findfont"
	"log"
	"os"
)

func init() {
	fontPath, err := findfont.Find("simkai.ttf")
	if err != nil {
		fontPath, err = findfont.Find("STHeiti Medium.ttc")
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FYNE_FONT", fontPath)
}

func main() {
	a := app.New()
	w := a.NewWindow("Wechat")
	w.Resize(fyne.Size{Width: 900, Height: 600})

	//account := widget.NewTextGrid()
	//w.SetContent(container.NewVBox(
	//	acount,
	//	widget.NewButton("Hi!", func() {
	//		acount.SetText("我勒个去 :)")
	//	}),
	//))

	w.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}
