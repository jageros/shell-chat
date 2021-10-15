/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    main
 * @Date:    2021/10/15 3:24 下午
 * @package: server
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/flags"
	"github.com/jageros/hawox/httpx"
	"github.com/jageros/hawox/ws"
	"gopkg.in/olahol/melody.v1"
)

const appName = "wechat"

func main() {
	ctx, wait := flags.Parse(appName)
	defer wait()
	httpx.InitializeHttpServer(ctx, func(engine *gin.Engine) {
		r := engine.Group("ws")
		ws.Init(ctx, r, flags.Source())
	})
	ws.OnMessage(func(session *melody.Session, bytes []byte) {
		ws.Broadcast()
	})
}
