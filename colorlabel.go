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
// Now it is also possible to set that too long text is truncated.
//
// Author: Reiner Pröls
// Licence: MIT

package colorlabel

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	_ fyne.WidgetRenderer    = (*ColorLabelRenderer)(nil)
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

	fullText  string
	bgColor   any
	fgColor   any
	textScale float32
	textStyle *fyne.TextStyle
	truncate  bool

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
		fullText:  s,
		textStyle: &fyne.TextStyle{},
	}

	colorLabel.ExtendBaseWidget(colorLabel)

	/*

		fyne.CurrentApp().Settings().AddListener(func(settings fyne.Settings) {
			colorLabel.fgColor = getColor(colorLabel.fgColor)
			colorLabel.bgColor = getColor(colorLabel.bgColor)
			colorLabel.Refresh()
		})
	*/
	return colorLabel
}

// Widget interface
func (l *ColorLabel) CreateRenderer() fyne.WidgetRenderer {
	t := canvas.NewText(l.fullText, getColor(l.fgColor))
	b := canvas.NewRectangle(getColor(l.bgColor))
	return &ColorLabelRenderer{
		w:    l,
		text: t,
		bg:   b,
		objs: []fyne.CanvasObject{b, t},
	}
}

// ColorLabelRenderer implements:
//   - fyne.WidgetRenderer
type ColorLabelRenderer struct {
	w        *ColorLabel
	text     *canvas.Text
	bg       *canvas.Rectangle
	objs     []fyne.CanvasObject
	maxWidth float32
}

// WidgetRenderer interface
func (r *ColorLabelRenderer) Layout(size fyne.Size) {
	pad := theme.Padding()
	s := fyne.NewSize(size.Width-2*pad, size.Height-2*pad)
	s2 := fyne.NewSize(size.Width, size.Height)
	p := fyne.NewPos(pad, pad)
	p2 := fyne.NewPos(0, 0)
	r.maxWidth = size.Width

	r.text.Resize(s)
	r.bg.Resize(s2)
	r.text.Move(p)
	r.bg.Move(p2)
}

// WidgetRenderer interface
func (r *ColorLabelRenderer) MinSize() fyne.Size {
	h := r.text.MinSize().Height + 2*theme.Padding()
	return fyne.NewSize(0, h)
}

// WidgetRenderer interface
func (r *ColorLabelRenderer) Refresh() {
	r.text.TextSize = theme.TextSize() * r.w.textScale
	r.text.TextStyle = *r.w.textStyle
	r.text.Text = r.w.truncateText(r.w.fullText, r.maxWidth, r.text)

	r.text.Color = getColor(r.w.fgColor)
	r.text.Refresh()
	r.bg.FillColor = getColor(r.w.bgColor)
	r.bg.Refresh()
}

// WidgetRenderer interface
func (r *ColorLabelRenderer) Destroy() {
}

func (r *ColorLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objs
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
	l.fullText = s
	l.Refresh()
}

func (l *ColorLabel) truncateText(s string, maxWidth float32, text *canvas.Text) string {
	if !l.truncate {
		return s
	}
	maxWidth -= theme.Padding() * 2
	ellipsis := "…"
	ellW := fyne.MeasureText(ellipsis, text.TextSize, text.TextStyle).Width

	r := []rune(s)
	if fyne.MeasureText(s, text.TextSize, text.TextStyle).Width <= maxWidth {
		return s
	}

	for len(r) > 0 {
		r = r[:len(r)-1]
		if fyne.MeasureText(string(r), text.TextSize, text.TextStyle).Width+ellW <= maxWidth {
			return string(r) + ellipsis
		}
	}
	return ellipsis
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
	l.fgColor = txtColor
	l.Refresh()
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
	l.bgColor = backColor
	l.Refresh()
	return nil
}

// Set new text scale factor
func (l *ColorLabel) SetTextScale(tScale float32) {
	if tScale <= 0 {
		tScale = 1
	}
	l.textScale = tScale
	l.Refresh()
}

// Set a text style
func (l *ColorLabel) SetTextStyle(textStyle *fyne.TextStyle) {
	if textStyle != nil {
		l.textStyle = textStyle
	} else {
		l.textStyle = &fyne.TextStyle{}
	}
	l.Refresh()
}

// Set text and text color
// txtColor is NRGBA or fyne.ThemeColorName
func (l *ColorLabel) SetTextWithColor(txt string, txtColor any) {
	l.fullText = txt
	l.SetTextColor(txtColor)
}

func (l *ColorLabel) SetTruncate(tr bool) {
	l.truncate = tr
	l.Refresh()
}
