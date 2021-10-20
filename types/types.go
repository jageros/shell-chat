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

type Msg struct {
	MsgID int    `json:"msg_id"`
	Seq   int64  `json:"seq"`
	Msg   string `json:"msg"`
}
