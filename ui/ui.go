package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fiurthorn/voucher/config"
	"github.com/fiurthorn/voucher/digistore24"
)

type UI struct {
	w fyne.Window

	apiKeyValue     *myEntry
	productsValue   *myEntry
	voucherIdsValue *myEntry

	send *widget.Button

	grid   *widget.Form
	layout *fyne.Container
}

func Create(w fyne.Window) *UI {
	u := &UI{w: w}

	ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlE, u.edit)

	ctrlS := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlS, u.save)

	return u
}

func (u *UI) edit(shortcut fyne.Shortcut) {
	if u.apiKeyValue.Disabled() {
		u.apiKeyValue.Enable()
		u.productsValue.Enable()
		u.w.Canvas().Focus(u.apiKeyValue)
	} else {
		u.apiKeyValue.Disable()
		u.productsValue.Disable()
		u.w.Canvas().Focus(u.voucherIdsValue)
	}
}

func (u *UI) save(shortcut fyne.Shortcut) {
	config.Config.ApiKey = u.apiKeyValue.Text
	config.Config.Products = strings.Split(u.productsValue.Text, ",")
	file, err := config.StoreConfig()
	if err != nil {
		dialog.ShowError(err, u.w)
	} else {
		dialog.ShowInformation("saved", file, u.w)
		u.apiKeyValue.Disable()
		u.productsValue.Disable()
		u.w.Canvas().Focus(u.voucherIdsValue)
	}
}

func (u *UI) Widget() *fyne.Container {
	u.apiKeyValue = NewEntry(u.save, u.edit, u.focus, config.Config.ApiKey, true)

	u.productsValue = NewEntry(u.save, u.edit, u.focus, strings.Join(config.Config.Products, ","), true)

	u.voucherIdsValue = NewEntry(u.save, u.edit, u.focus, "", false)
	u.voucherIdsValue.MultiLine = true
	u.voucherIdsValue.SetMinRowsVisible(10)
	u.voucherIdsValue.PlaceHolder = "enter vochers here"

	u.send = widget.NewButtonWithIcon("create", theme.DocumentCreateIcon(), u.createVouchers())

	u.grid = widget.NewForm(
		widget.NewFormItem("ApiKey", u.apiKeyValue),
		widget.NewFormItem("Products", u.productsValue),
		widget.NewFormItem("VoucherIds", u.voucherIdsValue),
	)
	u.grid.SubmitText = "create"
	u.grid.OnSubmit = u.createVouchers()

	u.layout = container.New(layout.NewMaxLayout(), u.grid)
	u.w.Resize(fyne.Size{Width: 650, Height: u.layout.Size().Height})

	return u.layout
}

func (u *UI) focus(w fyne.Focusable) {
	u.w.Canvas().Focus(w)
}

func (u *UI) Focus() *UI {
	u.focus(u.voucherIdsValue)
	return u
}

func (u *UI) products() string {
	return strings.ReplaceAll(u.productsValue.Text, " ", "")
}

func (u *UI) apiKey() string {
	return strings.ReplaceAll(u.apiKeyValue.Text, " ", "")
}

func (u *UI) vouchers() []string {
	content := strings.ReplaceAll(u.voucherIdsValue.Text, "\r\n", "\n")
	content = strings.ReplaceAll(content, " ", "")
	return strings.Split(content, "\n")
}

func (u *UI) createVouchers() func() {
	return func() {
		u.send.Disable()
		defer u.send.Enable()

		u.voucherIdsValue.Disable()
		defer u.voucherIdsValue.Enable()

		drain := make(chan digistore24.CreateVoucherResult, 5)
		go digistore24.CreateVouchers(u.apiKey(), u.vouchers(), u.products(), drain)

		text := u.voucherIdsValue.Text
		for cr := range drain {
			var replace string = "failure"

			if cr.Err != nil {
				replace = fmt.Sprintf("%s [%d]: %s", cr.Voucher, cr.Status, cr.Err.Error())
			} else if cr.Status == 200 && cr.Result.Result == "error" {
				replace = fmt.Sprintf("%s [%d]: %s %s", cr.Voucher, cr.Result.Code, cr.Result.Message, cr.Result.Result)
			} else if cr.Status == 200 && cr.Result.Result == "success" {
				replace = fmt.Sprintf("%s [%s] %s", cr.Voucher, cr.Result.Data.CouponId, cr.Result.Data.Note)
			}

			text = strings.Replace(text, cr.Voucher, replace, 1)
			u.voucherIdsValue.SetText(text)
		}
	}
}
