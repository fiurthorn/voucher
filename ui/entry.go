package ui

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Handler func(shortcut fyne.Shortcut)
type FocusHandler func(fyne.Focusable)

func NewEntry(save, edit Handler, focus FocusHandler, text string, disable bool) *myEntry {
	e := &myEntry{
		save:  save,
		edit:  edit,
		focus: focus,
	}
	e.ExtendBaseWidget(e)
	e.Wrapping = fyne.TextTruncate
	e.TextStyle = fyne.TextStyle{Monospace: true}
	if disable {
		e.Validator = func(s string) error {
			if strings.ContainsRune(s, ' ') {
				return fmt.Errorf("space not allowed")
			}
			return nil
		}
		e.action = widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
			if e.Disabled() {
				e.Enable()
				e.action.SetIcon(theme.VisibilityOffIcon())
				e.focus(e)
			} else {
				e.Disable()
				e.action.SetIcon(theme.VisibilityIcon())
			}
		})
		e.ActionItem = e.action
		e.Refresh()
		e.SetText(text)
		e.Disable()
	} else {
		e.Validator = func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("empty field")
			}
			return nil
		}
	}
	return e
}

type myEntry struct {
	widget.Entry

	focus FocusHandler
	save  Handler
	edit  Handler

	action *widget.Button
}

func (m *myEntry) TypedShortcut(s fyne.Shortcut) {
	if c, ok := s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		return
	} else {
		if c.Modifier == fyne.KeyModifierControl {
			if c.KeyName == fyne.KeyE {
				m.edit(s)
				return
			} else if c.KeyName == fyne.KeyS {
				m.save(s)
				return
			}
		}
	}

	log.Println("Shortcut typed:", s)
}
