// Copyright (c) 2025 Reiner Pröls
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT

package main

import (
	"image/color"

	"github.com/bytemystery-com/colorlabel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.NewWithID("com.bytemystery.colorlabel")
	w := a.NewWindow("ColorLabel")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(1000, 800))
	w.CenterOnScreen()

	label1 := colorlabel.NewColorLabel("Hallo", "", color.Transparent, 1.0)
	label2 := colorlabel.NewColorLabel("Text in red", color.NRGBA{R: 255, G: 0, B: 0, A: 255}, "", 1.0)
	label3 := colorlabel.NewColorLabel("Click me", theme.ColorNameForeground, theme.ColorNameSelection, 1.0)
	label3.OnTapped = func() {
		appearance := settings.NewSettings().LoadAppearanceScreen(w)
		dialog.ShowCustom("Fyne theme settings", "Ok", appearance, w)
	}
	label3.OnTappedSecondary = func() {
		n := fyne.NewNotification("Secondary click", "Secondary click on colorlabel")
		a.SendNotification(n)
	}
	label4 := colorlabel.NewColorLabel("Status text", theme.ColorNameForegroundOnError, theme.ColorNameError, 1.0)
	label5 := colorlabel.NewColorLabel("Blue text on gray", color.NRGBA{R: 0, G: 0, B: 255, A: 255}, color.NRGBA{R: 192, G: 192, B: 192, A: 255}, 1.0)
	var label6 *colorlabel.ColorLabel
	label6 = colorlabel.NewColorLabel("Click for changing color or size", color.NRGBA{R: 0, G: 0, B: 255, A: 255}, color.NRGBA{R: 192, G: 192, B: 192, A: 255}, 1.0)
	label6.OnTapped = func() {
		label6.SetTextWithColor("Now in red", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		label6.SetTextScale(1.0)
	}
	label6.OnTappedSecondary = func() {
		label6.SetTextScale(2.0)
	}
	label6.OnDoubleTapped = func() {
		label6.SetTextWithColor("Double tapped", color.NRGBA{R: 255, G: 255, B: 0, A: 255})
		label6.SetTextScale(1.0)
	}
	var label7 *colorlabel.ColorLabel
	label7 = colorlabel.NewColorLabel("With textstyle", color.NRGBA{R: 0, G: 0, B: 255, A: 255}, color.NRGBA{R: 192, G: 192, B: 192, A: 255}, 1.0)
	label7.SetTextStyle(&fyne.TextStyle{Bold: true, Italic: true})

	var label8 *colorlabel.ColorLabel
	label8 = colorlabel.NewColorLabel("Monospace - il1234W", color.NRGBA{R: 0, G: 0, B: 255, A: 255}, color.NRGBA{R: 192, G: 192, B: 192, A: 255}, 1.0)
	label8.SetTextStyle(&fyne.TextStyle{Monospace: true})

	var label9 *colorlabel.ColorLabel
	label9 = colorlabel.NewColorLabel("A very long text - Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.", "", "", 1.0)
	label9.SetTruncate(true)

	vbox := container.NewGridWrap(fyne.NewSize(w.Canvas().Size().Width, 50), label1, label2, label3, label4, label5, label6, label7, label8, label9)
	w.SetContent(vbox)

	w.ShowAndRun()
}
