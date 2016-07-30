package main

import (
	"image"

	ui "github.com/pdevine/termui"
)

type Whale struct {
	ui.BaseSprite
	v     image.Point
	blink bool
}

const whale_c0 = `xxxxxxxxxxxxxxx##xxxxxxxx.xxxxx
xxxxxxxxx##x##x##xxxxxxx==xxxxx
xxxxxx##x##x##x##xxxxxx===xxxxx
xx/""""""""""""""""\___/x===xxx
x{                      /xx===x
xx\______ o          __/xxxxxxx
xxxx\    \        __/xxxxxxxxxx
xxxxx\____\______/xxxxxxxxxxxxx`

const whale_c1 = `xxxxxxxxxxxxxxx##xxxxxxxx.xxxxx
xxxxxxxxx##x##x##xxxxxxxx==xxxx
xxxxxx##x##x##x##xxxxxxx===xxxx
xx/""""""""""""""""\___/===xxxx
x{                      /x===xx
xx\______ -          __/xxxxxxx
xxxx\    \        __/xxxxxxxxxx
xxxxx\____\______/xxxxxxxxxxxxx`

func NewWhale() *Whale {
	s := ui.NewBaseSprite(whale_c0)
	//s.AddCostume(ui.Costume{whale_c0})
	s.AddCostume(ui.Costume{whale_c1})
	return &Whale{
		BaseSprite: *s,
		v:          image.Point{1, 1},
	}
}

func (w *Whale) Update(t ui.EvtTimer) error {
	// bounds detection
	win = getwinsize()
	if w.X+w.Width > int(win.ws_col) || w.X <= 0 {
		w.v.X = -w.v.X
	}
	if w.Y+w.Height > int(win.ws_row) || w.Y <= 0 {
		w.v.Y = -w.v.Y
	}
	w.X += w.v.X
	w.Y += w.v.Y

	if t.Count%1000 == 0 {
		w.NextCostume()
	}
	return nil
}
