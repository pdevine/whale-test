package main

import ui "github.com/pdevine/termui"

type Fish struct {
	ui.BaseSprite
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
	}
}

func (f *Fish) Update(t ui.EvtTimer) error {
	// bounds detection
	if t.Count%40 == 0 {
		f.Wrap()
		f.UpdatePosition()
	}
	return nil
}
