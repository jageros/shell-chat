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
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jroimartin/gocui"
	"log"
)

var (
	viewArr = []string{"v1", "v2"}
	active  = 1
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

	out, err := g.View("v2")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "\nGoing from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 1 {
		fmt.Fprintln(out, "")
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX/4*3-1, maxY/5*4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "message"
		//v.Editable = true
		v.Wrap = true
	}

	if v, err := g.SetView("v2", 0, maxY/5*4-1, maxX/4*3-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "send"
		v.Editable = true
		v.Overwrite = true
		v.Wrap = true
		v.Autoscroll = true
		if _, err = setCurrentViewOnTop(g, "v2"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("v3", maxX/4*3-1, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "online"
		v.Wrap = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func Init(ctx contextx.Context) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	ctx.Go(func(ctx contextx.Context) error {
		return g.MainLoop()
	})
}
