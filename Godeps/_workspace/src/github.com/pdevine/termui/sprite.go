// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "strings"

type Costume struct {
	Text string
}

type BaseSprite struct {
	Block
	Costumes       []Costume
	CurrentCostume int
	TextFgColor    Attribute
	TextBgColor    Attribute
	WrapLength     int // words wrap limit. Note it may not work properly with multi-width char
	Layer          int
	Alpha          rune
}

type Sprite interface {
	Update(EvtTimer) error
	Buffer() Buffer
}

type SpriteGroup struct {
	Sprites []Sprite
}

type ByLayer []BaseSprite

func (l ByLayer) Len() int           { return len(l) }
func (l ByLayer) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ByLayer) Less(i, j int) bool { return l[i].Layer < l[j].Layer }

func NewSpriteGroup() *SpriteGroup {
	return &SpriteGroup{}
}

func (g *SpriteGroup) Add(s Sprite) {
	g.Sprites = append(g.Sprites, s)
	//sort.Sort(ByLayer(g.Sprites))
}

func (g *SpriteGroup) Render() {
	// XXX - this feels like it's gonna be slow
	b := make([]Bufferer, len(g.Sprites))
	for cnt, s := range g.Sprites {
		b[cnt] = s
	}
	Render(b...)
}

// NewSprite returns a new *Sprite with an single custome.
func NewBaseSprite(s string) *BaseSprite {
	block := NewBlock()
	block.Width, block.Height = getWidthHeight(s)
	block.Border = false
	block.BorderLabel = ""

	c := Costume{s}

	return &BaseSprite{
		Block:       *block,
		Costumes:    []Costume{c},
		TextFgColor: ThemeAttr("par.text.fg"),
		TextBgColor: ThemeAttr("par.text.bg"),
		WrapLength:  0,
		Alpha:       'x',
	}
}

func getWidthHeight(s string) (int, int) {
	var width int
	var height int

	for _, line := range strings.Split(s, "\n") {
		if len(line) > width {
			width = len(line)
		}
		height += 1
	}
	return width, height
}

func (s *BaseSprite) AddCostume(c Costume) error {
	if len(s.Costumes) == 1 && s.Costumes[0].Text == "" {
		s.Block.Width, s.Block.Height = getWidthHeight(c.Text)
		s.Costumes[0] = c
	} else {
		s.Costumes = append(s.Costumes, c)
	}
	return nil
}

func (s *BaseSprite) NextCostume() {
	s.CurrentCostume++
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}
}

func (s *BaseSprite) Update(t EvtTimer) error {
	return nil
}

// Buffer implements Bufferer interface.
func (s BaseSprite) Buffer() Buffer {
	buf := s.Block.Buffer()

	fg, bg := s.TextFgColor, s.TextBgColor
	cs := DefaultTxBuilder.Build(s.Costumes[s.CurrentCostume].Text, fg, bg)

	y, x, n := 0, 0, 0
	for y < s.innerArea.Dy() && n < len(cs) {
		w := cs[n].Width()

		if cs[n].Ch == s.Alpha {
			cs[n].Visible = false
		}

		if cs[n].Ch == '\n' || x+w > s.innerArea.Dx() {
			y++
			x = 0 // set x = 0
			if cs[n].Ch == '\n' {
				n++
			}
			continue
		}

		buf.Set(s.innerArea.Min.X+x, s.innerArea.Min.Y+y, cs[n])

		n++
		x += w
	}

	return buf
}
