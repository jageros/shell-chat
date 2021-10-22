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
	"fmt"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/wsc"
	"log"
	"net/http"
	"unicode"
	"wechat/types"
	"wechat/view"
)

func main() {
	ctx, cancel := contextx.Default()
	defer cancel()
	var err error

	var sess *wsc.Session
	m := wsc.New(ctx)

	view.OnMessage("请发送本人手机号码加入聊天室！\n")

	view.OnSendMsg(func(msg string) {
		if sess == nil {
			if len(msg) != 11 {
				view.OnMessage("手机号码格式错误！\n")
				return
			}

			for _, r := range msg {
				if !unicode.Is(unicode.Number, r) {
					view.OnMessage("手机号码格式错误！\n")
					return
				}
			}

			uid := msg
			h := http.Header{}
			h.Add("uid", uid)
			sess, err = m.ConnectWithHeader("ws://wechat.hawtech.cn/ws/wechat/1", h, map[string]interface{}{"uid": uid})
			if err != nil {
				view.OnMessage(fmt.Sprintf("创建websocket链接错误，手机号码:%s 错误信息:%v\n", uid, err))
				return
			}
			view.OnMessage("登录成功！\n")
			return
		}
		data := &types.Msg{
			MsgID: 1,
			Msg:   msg,
		}
		bytes, err := types.Marshal(data)
		if err != nil {
			view.OnMessage(fmt.Sprintf("编码错误: %v", err))
			return
		}
		err = sess.WriteBinary(bytes)
		if err != nil {

		}
	})

	m.HandleMessageBinary(func(session *wsc.Session, bytes []byte) {
		msg, err := types.Unmarshal(bytes)
		if err != nil {
			view.OnMessage(fmt.Sprintf("解码错误：%v", err))
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
