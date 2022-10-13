package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fiurthorn/voucher/config"
	"github.com/fiurthorn/voucher/digistore24"
)

func Widget() *fyne.Container {
	apiKeyLabel := widget.NewLabel("ApiKey")
	apiKeyValue := widget.NewLabel(config.Config.ApiKey)

	productsLabel := widget.NewLabel("Products")
	productsValue := widget.NewLabel(strings.Join(config.Config.Products, " | "))

	voucherIdsLabel := widget.NewLabel("VoucherIds")
	voucherIdsValue := widget.NewEntry()
	voucherIdsValue.SetMinRowsVisible(20)
	voucherIdsValue.MultiLine = true
	voucherIdsValue.PlaceHolder = "enter vochers here"

	var send *widget.Button
	send = widget.NewButtonWithIcon("create", theme.DocumentCreateIcon(), func() {
		send.Disable()
		voucherIdsValue.Disable()
		drain := make(chan digistore24.CallResult, 5)
		go digistore24.Call(voucherIdsValue.Text, drain)

		text := voucherIdsValue.Text
		for cr := range drain {
			var replace string = "failure"

			if cr.Err != nil {
				replace = fmt.Sprintf("%d: %s", cr.Status, cr.Err.Error())
			} else if cr.Status == 200 && cr.Result.Result == "error" {
				replace = fmt.Sprintf("%d: %s", cr.Result.Code, cr.Result.Message)
			} else if cr.Status == 200 && cr.Result.Result == "success" {
				replace = fmt.Sprintf("%s [%s] %s", cr.Result.Data.Code, cr.Result.Data.Code, cr.Result.Data.Note)
			}
			text = strings.Replace(text, cr.Voucher, replace, 1)
			voucherIdsValue.SetText(text)
		}

		voucherIdsValue.Enable()
		send.Enable()
	})

	grid := container.New(layout.NewFormLayout(),
		apiKeyLabel, apiKeyValue,
		productsLabel, productsValue,
		voucherIdsLabel, voucherIdsValue,
		widget.NewLabel(""), send,
	)

	gridMax := container.New(layout.NewMaxLayout(), grid)

	return gridMax

}
