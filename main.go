package main

import (
	"math/rand"
	"syscall"
	"time"
	"unsafe"

	ui "github.com/pdevine/termui"
)

type winsize struct {
	ws_row, ws_col       uint16
	ws_xpixel, ws_ypixel uint16
}

func getwinsize() winsize {
	ws := winsize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	return ws
}

var win winsize

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.DefaultEvtStream.Merge("timer", ui.NewTimerCh(time.Millisecond))

	sg := ui.NewSpriteGroup()

	f := NewFish()
	f.TextFgColor = ui.ColorWhite
	f.X = 7
	f.Y = 7

	sg.Add(f)

	f = NewFish()
	f.X = 20
	f.Y = 14
	f.v.X = 2

	sg.Add(f)

	for cnt := 0; cnt < 5; cnt++ {
		p := NewWhale()
		p.TextFgColor = ui.ColorWhite
		p.Y = rand.Intn(50) + 1
		p.X = rand.Intn(150) + 1
		sg.Add(p)
	}

	// event handler...
	// handle key q pressing
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		p := NewWhale()
		p.TextFgColor = ui.ColorWhite
		p.Y = rand.Intn(50) + 1
		p.X = rand.Intn(150) + 1
		sg.Add(p)
	})

	ui.Handle("/sys/kbd/z", func(ui.Event) {
		sg.Sprites = sg.Sprites[:len(sg.Sprites)-1]
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
		// handle Ctrl + x combination
	})

	ui.Handle("/sys/kbd", func(ui.Event) {
		// handle all other key pressing
	})

	// handle a 1s timer
	ui.Handle("/timer/1ms", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		// t is a EvtTimer
		if t.Count%20 == 0 {
			for _, s := range sg.Sprites {
				s.Update(t)
			}
		}

		if t.Count%10 == 0 {
			ui.Clear()
			sg.Render()
		}
	})

	ui.Loop() // block until StopLoop is called
}
