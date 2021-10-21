/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    service
 * @Date:    2021/7/8 5:30 下午
 * @package: service
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jageros/hawox/contextx"
	"github.com/jageros/hawox/errcode"
	"github.com/jageros/hawox/httpx"
	"github.com/jageros/hawox/logx"
	"gopkg.in/olahol/melody.v1"
	"sync"
	"time"
	"wechat/types"
)

var (
	seq   int64
	mux   sync.Mutex
	names = map[string]string{
		"13160676597": "杰",
		"13612225480": "文",
		"13750043941": "哲",
	}
)

var ss *service

type service struct {
	ctx         contextx.Context
	m           *melody.Melody
	callTimeout time.Duration
	online      map[string]struct{}
	mx          sync.Mutex
}

func Init(ctx contextx.Context, r *gin.RouterGroup, relativePath string) {
	ss = &service{
		ctx:         ctx,
		m:           melody.New(),
		callTimeout: time.Second * 5,
		online:      map[string]struct{}{},
	}
	ss.m.HandleMessage(ss.handleMessage)
	ss.m.HandleConnect(ss.onConnect)
	ss.m.HandleDisconnect(ss.onDisconnect)
	r.GET(relativePath, ss.handler)
}

func (s *service) handler(c *gin.Context) {
	uid := c.GetHeader("uid")
	if _, ok := names[uid]; !ok {
		logx.Infof("=====uid=%s", uid)
		httpx.ErrInterrupt(c, errcode.InvalidParam.WithMsg(uid))
		return
	}
	err := s.m.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"uid": uid})
	if err != nil {
		httpx.ErrInterrupt(c, errcode.WithErrcode(-1, err))
	}
}

func (s *service) onConnect(session *melody.Session) {
	uid, exist := session.Get("uid")
	logx.Infof("on connect uid=%s", uid)
	if !exist {
		return
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	s.online[uid.(string)] = struct{}{}
	s.updateOnline()
}

func (s *service) onDisconnect(session *melody.Session) {
	uid, exist := session.Get("uid")
	logx.Infof("on disconnect uid=%s", uid)
	if !exist {
		return
	}
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.online, uid.(string))
	s.updateOnline()
}

func (s *service) updateOnline() {
	var msg string
	for id := range s.online {
		name := names[id]
		msg = fmt.Sprintf("%s%s(%s)\n", msg, name, id)
	}
	var resp = &types.Msg{
		MsgID: 2,
		Msg:   msg,
	}
	bty, err := json.Marshal(resp)
	if err != nil {
		return
	}
	err = s.m.BroadcastBinary(bty)
	if err != nil {
		logx.Error(err)
	}
	logx.Infof("update %s", msg)
}

func (s *service) handleMessage(session *melody.Session, bytes []byte) {
	start := time.Now()
	uid, exist := session.Get("uid")
	logx.Infof("on msg uid=%s", uid)
	if !exist {
		return
	}

	name := names[uid.(string)]

	mux.Lock()
	defer mux.Unlock()
	seq += 1

	msg := fmt.Sprintf("[%d]%s(%s): %s\n", seq, name, time.Now().Format("15:04:05"), string(bytes))

	var resp = &types.Msg{
		MsgID: 1,
		Msg:   msg,
	}
	bty, err := json.Marshal(resp)
	if err != nil {
		return
	}
	err = s.m.BroadcastBinary(bty)
	if err != nil {
		logx.Error(err)
	}
	take := time.Now().Sub(start)
	if take > time.Millisecond*100 {
		logx.Warnf("send msg take: %s", take.String())
	}
}
