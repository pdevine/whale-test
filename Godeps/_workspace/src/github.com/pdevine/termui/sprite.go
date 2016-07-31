// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"strings"
	"syscall"
	"unsafe"
)

type Costume struct {
	Text string
}

type Vector struct {
	X float32
	Y float32
}

func NewVector(x, y int) *Vector {
	return &Vector{
		X: float32(x),
		Y: float32(y),
	}
}

type BaseSprite struct {
	Block
	Costumes       []Costume
	Position       Vector
	Velocity       Vector
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

type WinSize struct {
	Rows    uint16
	Columns uint16
	Xpixel  uint16
	Ypixel  uint16
}

func GetWinSize() WinSize {
	ws := WinSize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	return ws
}

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
	s.UpdatePosition()
	return nil
}

func (s *BaseSprite) SetPosition(x, y int) {
	s.Position.X = float32(x)
	s.Position.Y = float32(y)
	s.X = int(s.Position.X)
	s.Y = int(s.Position.Y)
}

func (s *BaseSprite) UpdatePosition() {
	// XXX - add proportional positions
	s.Position.X += s.Velocity.X
	s.Position.Y += s.Velocity.Y
	s.X = int(s.Position.X)
	s.Y = int(s.Position.Y)
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

func (s *BaseSprite) Bounce() {
	win := GetWinSize()
	if s.X+s.Width > int(win.Columns) {
		s.Position.X = float32(int(win.Columns) - s.Width)
		s.Velocity.X = -s.Velocity.X
	}
	if s.X <= 0 {
		s.Position.X = 0
		s.Velocity.X = -s.Velocity.X
	}
	if s.Y+s.Height > int(win.Rows) {
		s.Position.Y = float32(int(win.Rows) - s.Height)
		s.Velocity.Y = -s.Velocity.Y
	}
	if s.Y <= 0 {
		s.Position.Y = 0
		s.Velocity.Y = -s.Velocity.Y
	}
	s.X, s.Y = int(s.Position.X), int(s.Position.Y)
}

func (s *BaseSprite) Wrap() {
	win := GetWinSize()
	if s.X >= int(win.Columns) {
		s.Position.X = float32(0 - s.Width)
	}
	if s.X+s.Width < 0 {
		s.Position.X = float32(win.Columns)
	}
	if s.Y >= int(win.Rows) {
		s.Position.Y = float32(0 - s.Height)
	}
	if s.Y+s.Height < 0 {
		s.Position.Y = float32(win.Rows)
	}
	s.X, s.Y = int(s.Position.X), int(s.Position.Y)
}
