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

package main

import (
	"encoding/json"
	"flag"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/wsc"
	"log"
	"net/http"
	"os"
	"wechat/types"
	"wechat/view"
)

func main() {
	uid := flag.String("phone", "", "输入手机号码参数")
	flag.Parse()
	if *uid == "" {
		log.Fatal("请携带手机号码参数启动，--phone=10086")
		os.Exit(-1)
	}
	ctx, cancel := contextx.Default()
	defer cancel()

	h := http.Header{}
	h.Add("uid", *uid)

	m := wsc.New(ctx)
	sess, err := m.ConnectWithHeader("ws://127.0.0.1:8888/ws/wechat/1", h, map[string]interface{}{"uid": uid})
	if err != nil {
		panic(err)
	}

	view.OnSendMsg(func(msg string) {
		err = sess.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	})

	m.HandleMessageBinary(func(session *wsc.Session, bytes []byte) {
		//uid, _ := session.Get("uid")
		//roomId, _ := session.Get("roomId")
		msg := &types.Msg{}
		err := json.Unmarshal(bytes, msg)
		if err != nil {
			log.Panicf("msg.Unmarshal err: %v", err)
			return
		}
		switch msg.MsgID {
		case 1:
			view.OnMessage(msg.Msg)

		case 2:
			view.UpdateOnline(msg.Msg)

		default:
			log.Printf("MsgId=%d Msg=%s", msg.MsgID, msg.Msg)
		}
	})

	err = view.Run()
	if err != nil {
		log.Fatal(err)
	}
}
