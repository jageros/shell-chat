/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    types
 * @Date:    2021/10/20 5:03 下午
 * @package: types
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package types

import (
	"encoding/json"
	"math/rand"
)

type Msg struct {
	MsgID int    `json:"msg_id"`
	Seq   int64  `json:"seq"`
	Msg   string `json:"msg"`
}

type pkg struct {
	N  int
	V1 []byte
	V2 []byte
	V3 []byte
}

func Marshal(msg *Msg) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	pk := &pkg{}
	l := len(data)
	for i := 0; i < l; i += 2 {
		pk.N = 1
		pk.V1 = append(pk.V1, data[i])
		pk.V2 = append(pk.V2, uint8(rand.Intn(256)))
		if i+1 < l {
			pk.N = 2
			pk.V3 = append(pk.V3, data[i+1])
		}
	}
	return json.Marshal(pk)
}

func Unmarshal(data []byte) (*Msg, error) {
	var pk = &pkg{}
	err := json.Unmarshal(data, pk)
	if err != nil {
		return nil, err
	}

	var mds []byte
	if pk.N == 1 {
		l := len(pk.V3)
		for i, v := range pk.V1 {
			mds = append(mds, v)
			if i < l {
				mds = append(mds, pk.V3[i])
			}
		}
	} else {
		l := len(pk.V1)
		for i, v := range pk.V3 {
			if i < l {
				mds = append(mds, pk.V1[i])
			}
			mds = append(mds, v)
		}
	}
	res := &Msg{}
	err = json.Unmarshal(mds, res)
	return res, err
}
