package main

import ui "github.com/pdevine/termui"

type Whale struct {
	ui.BaseSprite
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
	s.Velocity = ui.Vector{1.0, 1.0}
	return &Whale{
		BaseSprite: *s,
	}
}

func (w *Whale) Update(t ui.EvtTimer) error {
	w.Bounce()

	if t.Count%1000 == 0 {
		w.NextCostume()
	}

	w.UpdatePosition()
	return nil
}
