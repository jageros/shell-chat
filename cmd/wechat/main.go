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
	"log"
	"wechat/view"
)

func main() {
	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	//defer cancel()

	view.OnMessage("dekdmkwenkwndklwenklndk\n")
	view.UpdateOnline("杰（13160676597）\n哲（10086）\n文（10010）\n")
	view.UpdateOnline("杰（13160676597）\n哲（10086）\n文（10010）\n")


	err := view.Run()
	if err != nil {
		log.Fatal(err)
	}
}
