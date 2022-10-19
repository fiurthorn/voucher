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

	app := app.New()
	window := app.NewWindow("Create Voucher | digistore24")
	u := ui.Create(window)

	window.SetContent(u.Widget())
	window.Show()
	u.Focus()
	app.Run()
}
