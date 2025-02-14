/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    view
 * @Date:    2021/10/15 6:27 下午
 * @package: view
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package view

import (
	"github.com/rocket049/gocui"
	"log"
)

var (
	viewArr = []string{"msg", "send"}
	active  = 1
	gg      *gocui.Gui
	onSend  func(msg string)
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", 0, 0, maxX/4*3-1, maxY/5*4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "message"
		v.Autoscroll = true
		v.Wrap = true
	}

	if v, err := g.SetView("send", 0, maxY/5*4, maxX/4*3-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "send"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		if _, err = setCurrentViewOnTop(g, "send"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("online", maxX/4*3, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "online"
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func OnSendMsg(fn func(msg string)) {
	onSend = fn
}

func OnMessage(msg string) {
	gg.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("msg")
		if err != nil {
			return err
		}
		_, err = v.Write([]byte(msg))
		return err
	})
}

func UpdateOnline(msg string) {
	gg.Update(func(gui *gocui.Gui) error {
		v, err := gui.View("online")
		if err != nil {
			return err
		}
		v.Clear()
		_, err = v.Write([]byte(msg))
		return err
	})
}

func sendMsg(g *gocui.Gui, v *gocui.View) error {
	byts := v.ReadEditor()
	if len(byts) <= 0 {
		v.Clear()
		return v.SetCursor(0, 0)
	}
	str := string(byts)

	if onSend != nil {
		onSend(str)
	}

	v.Clear()
	return v.SetCursor(0, 0)
}

func arrowUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func arrowDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func backspace(g *gocui.Gui, v *gocui.View) error {
	v.EditDelete(true)
	return nil
}

func clear(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	return v.SetCursor(0, 0)
}

func init() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("send", gocui.KeyDelete, gocui.ModNone, backspace); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("send", gocui.KeyBackspace, gocui.ModNone, backspace); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("send", gocui.KeyCtrlD, gocui.ModNone, clear); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("send", gocui.KeyBackspace2, gocui.ModNone, backspace); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("send", gocui.KeyEnter, gocui.ModNone, sendMsg); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("msg", gocui.KeyArrowUp, gocui.ModNone, arrowUp); err != nil {
		log.Panic(err)
	}

	if err := g.SetKeybinding("msg", gocui.KeyArrowDown, gocui.ModNone, arrowDown); err != nil {
		log.Panic(err)
	}
	gg = g
}

func Run() error {
	defer gg.Close()
	return gg.MainLoop()
}
