package rpc

import (
	"github.com/davyxu/cellnet"
	"time"
)

// ud: peer/session,   reqMsg:请求用的消息, userCallback: 返回消息类型回调 func( ackMsg *ackMsgType)
func Call(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, userCallback func(raw interface{})) {

	ses, err := getPeerSession(sesOrPeer)

	if err != nil {
		cellnet.QueuedCall(ses, func() {
			userCallback(err)
		})

		return
	}

	// 发送RPC请求
	req := createRequest(userCallback)

	req.Send(ses, reqMsg)

	// 等待RPC回复
	time.AfterFunc(timeout, func() {

		if requestExists(req.id) {
			cellnet.QueuedCall(ses, func() {
				userCallback(ErrTimeout)
			})
		}
	})
}
