package main

import (
	"image"

	ui "github.com/pdevine/termui"
)

type Fish struct {
	ui.BaseSprite
	v image.Point
}

const shark = `xxxxxxxxx,xxxxxx
xxxxxxx__)` + "\\" + `_xxxx
x(` + "\\" + `_.-'    a` + "`" + `-.x
x(/~~` + "````" + `(/~^^` + "`" + `x
xxxxxxxxxxxxxxxx`

func NewFish() *Fish {
	s := ui.NewBaseSprite(shark)
	return &Fish{
		BaseSprite: *s,
		v:          image.Point{1, 0},
	}
}

func (f *Fish) Update(t ui.EvtTimer) error {
	// bounds detection
	if t.Count%40 == 0 {
		win = getwinsize()
		f.X += f.v.X
		if f.X > int(win.ws_col) {
			f.X = -f.Width
		}
	}
	return nil
}
