package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/fiurthorn/voucher/config"
	"github.com/fiurthorn/voucher/ui"
)

func init() {
	os.Setenv("FYNE_THEME", "light")
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	log.Printf("%+v", config.Config)

	myApp := app.New()
	myWindow := myApp.NewWindow("Create Voucher | digistore24")

	myWindow.SetContent(ui.Widget())

	// myWindow.Resize(fyne.Size{Width: 800, Height: 600})
	myWindow.Show()
	myApp.Run()
}
