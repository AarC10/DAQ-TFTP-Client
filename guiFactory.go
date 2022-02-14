/**
Factory methods for creating GUI elements.

@Author Aaron Chan
*/

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

/**
Factory method for creating a text field
*/
func makeEntryField(text string, validationType string) *widget.Entry {
	newField := widget.NewEntry()
	newField.SetPlaceHolder(text)

	// Set validator function
	if validationType == "ip" {
		newField.Validator = ipAddrValidator
	} else {
		newField.Validator = portValidator
	}

	return newField
}

/**
Factory method for creating texts
*/
func makeNewText(text string) *canvas.Text {
	newText := canvas.NewText(text, color.White)
	newText.Alignment = fyne.TextAlignCenter

	return newText
}

/**
Factory method for creating selection box
*/
func makeNewSelection(text string, options []string) *widget.Select {
	newSelection := widget.NewSelect(options, func(value string) {
		// log.Println(text, "set to", value)
	})

	return newSelection
}
