/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2021/10/18 11:15 上午
 * @package: wechat
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

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

package main

import (
	"context"
	"fmt"
	"github.com/jageros/hawox/logx"
	"github.com/jroimartin/gocui"
	"os/signal"
	"syscall"
	"time"
)

var (
	viewArr = []string{"msg", "send"}
	active  = 1
	g       *gocui.Gui
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

	out, err := g.View("send")
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
	if v, err := g.SetView("msg", 0, 0, maxX/4*3-1, maxY/5*4-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "message"
		//v.Editable = true
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
	if v, err := g.SetView("v3", maxX/4*3, 0, maxX-1, maxY-1); err != nil {
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

func display(msg string) error {
	v, err := g.View("msg")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(v, msg)
	return err
}

func read() (chan string, error) {
	var ch = make(chan string, 100)

	v, err := g.View("send")
	if err != nil {
		return nil, err
	}
	go func() {
		var str string
		for {
			_, err := fmt.Fscanln(v, &str)
			if err == nil {
				ch <- str
			}
		}
	}()

	return ch, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	logx.Init(logx.InfoLevel, logx.SetFileOut("logs", "wechat"))
	logx.Info("=============")
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		logx.Panic(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		logx.Panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		logx.Panic(err)
	}

	go func() {
		time.Sleep(time.Second*2)
		logx.Info("-----------")
		ch, err := read()
		if err != nil {
			logx.Error(err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				logx.Info(msg)
				err := display(msg)
				if err != nil {
					return
				}
			}
		}
	}()

	err = g.MainLoop()
	if err != nil {
		logx.Error(err)
	}
}
