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
//
// Implements a label with colored text and background as fyne widget.
// You can use color names defined by Fyne theme or direct NRGBA values.
// The labels can also be clicked (primary and secondary) and double clicked if needed.
// You can also set a Text style for bold and italic and Monospace font.
//
// Author: Reiner Pröls
// Licence: MIT

package colorlabel

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	_ fyne.Widget            = (*ColorLabel)(nil)
	_ fyne.Tappable          = (*ColorLabel)(nil)
	_ fyne.DoubleTappable    = (*ColorLabel)(nil)
	_ fyne.SecondaryTappable = (*ColorLabel)(nil)
	_ desktop.Mouseable      = (*ColorLabel)(nil)
)

// Color label text with color (text and background), also clickable if needed
// Implements
//   - fyne.Widget
//   - fyne.Tappable
//   - fyne.DoubleTappable
//	 - fyne.SecondaryTappable
//   - desktop.Mouseable

type ColorLabel struct {
	widget.BaseWidget
	text *canvas.Text
	bg   *canvas.Rectangle

	bgColor   any
	fgColor   any
	textScale float32
	textStyle *fyne.TextStyle

	OnTapped          func()
	OnTappedSecondary func()
	OnDoubleTapped    func()
	lastKeyModifier   fyne.KeyModifier
}

func getColor(c any) color.Color {
	switch v := any(c).(type) {
	case string:
		return theme.Color(fyne.ThemeColorName(v))
	case fyne.ThemeColorName:
		return theme.Color(v)
	case color.NRGBA:
		return v
	case color.Alpha16:
		return v
	case color.Gray16:
		return v
	}
	return color.Transparent
}

// Creates a new ColorLabel
// txtColor is NRGBA or fyne.ThemeColorName
// backColor is NRGBA or fyne.ThemeColorName
func NewColorLabel(s string, txtColor, backColor any, tScale float32) *ColorLabel {
	switch c := any(backColor).(type) {
	case fyne.ThemeColorName, string:
		if c == "" {
			backColor = color.Transparent
		}
	case color.NRGBA:
		backColor = c
	case color.Alpha16:
		backColor = c
	case color.Gray16:
		backColor = c
	default:
		return nil
	}

	switch c := any(txtColor).(type) {
	case fyne.ThemeColorName, string:
		if c == "" {
			txtColor = theme.ColorNameForeground
		}
	case color.NRGBA:
		txtColor = c
	case color.Alpha16:
		txtColor = c
	case color.Gray16:
		txtColor = c
	default:
		return nil
	}

	if tScale <= 0 {
		tScale = 1
	}

	colorLabel := &ColorLabel{
		bgColor:   backColor,
		fgColor:   txtColor,
		textScale: tScale,
		text:      canvas.NewText(s, getColor(txtColor)),
		bg:        canvas.NewRectangle(getColor(backColor)),
	}

	colorLabel.text.TextSize = theme.TextSize() * tScale
	if colorLabel.textStyle != nil {
		colorLabel.text.TextStyle = *colorLabel.textStyle
	}
	colorLabel.ExtendBaseWidget(colorLabel)
	fyne.CurrentApp().Settings().AddListener(func(settings fyne.Settings) {
		colorLabel.text.Color = getColor(colorLabel.fgColor)
		colorLabel.text.TextSize = theme.TextSize() * tScale
		colorLabel.bg.FillColor = getColor(colorLabel.bgColor)
		colorLabel.text.Refresh()
		colorLabel.bg.Refresh()
	})
	return colorLabel
}

// Widget interface
func (l *ColorLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		l.bg,
		container.NewPadded(l.text),
	))
}

// Tappable interface
func (l *ColorLabel) Tapped(ev *fyne.PointEvent) {
	if l.OnTapped != nil {
		l.OnTapped()
	}
}

// SecondaryTappable interface
func (l *ColorLabel) TappedSecondary(*fyne.PointEvent) {
	if l.OnTappedSecondary != nil {
		l.OnTappedSecondary()
	}
}

// DoubleTappable interface
func (l *ColorLabel) DoubleTapped(ev *fyne.PointEvent) {
	if l.OnDoubleTapped != nil {
		l.OnDoubleTapped()
	}
}

// Mouseable interface
func (l *ColorLabel) MouseDown(ev *desktop.MouseEvent) {
}

// Mouseable interface
func (l *ColorLabel) MouseUp(ev *desktop.MouseEvent) {
	l.lastKeyModifier = ev.Modifier
}

// User functions
// Get the last keyboard modifier
func (l *ColorLabel) GetLastKeyModifier() fyne.KeyModifier {
	return l.lastKeyModifier
}

// Set new text
func (l *ColorLabel) SetText(s string) {
	l.text.Text = s
	l.text.Refresh()
}

// Set new text color
// txtColor is NRGBA or fyne.ThemeColorName
func (l *ColorLabel) SetTextColor(txtColor any) error {
	switch c := txtColor.(type) {
	case fyne.ThemeColorName, string:
		if c == "" {
			txtColor = theme.ColorNameForeground
		}
	case color.NRGBA:
		txtColor = c
	case color.Alpha16:
		txtColor = c
	case color.Gray16:
		txtColor = c
	default:
		return errors.New("fyne.ThemeColorName or color.NRGBA required")
	}
	l.text.Color = getColor(txtColor)
	l.text.Refresh()
	return nil
}

// Set new background color
// backColor is NRGBA or fyne.ThemeColorName
func (l *ColorLabel) SetBackgroundColor(backColor any) error {
	switch c := backColor.(type) {
	case fyne.ThemeColorName, string:
		if c == "" {
			backColor = color.Transparent
		}
	case color.NRGBA:
		backColor = c
	case color.Alpha16:
		backColor = c
	case color.Gray16:
		backColor = c
	default:
		return errors.New("fyne.ThemeColorName or color.NRGBA required")
	}
	l.bg.FillColor = getColor(backColor)
	l.bg.Refresh()
	return nil
}

// Set new text scale factor
func (l *ColorLabel) SetTextScale(tScale float32) {
	if tScale <= 0 {
		tScale = 1
	}
	l.text.TextSize = theme.TextSize() * tScale
	l.textScale = tScale
	l.text.Refresh()
}

// Set a text style
func (l *ColorLabel) SetTextStyle(textStyle *fyne.TextStyle) {
	l.textStyle = textStyle
	l.text.TextStyle = *textStyle
	l.text.Refresh()
}

// Set text and text color
// txtColor is NRGBA or fyne.ThemeColorName
func (l *ColorLabel) SetTextWithColor(txt string, txtColor any) {
	l.text.Text = txt
	l.SetTextColor(txtColor)
}
