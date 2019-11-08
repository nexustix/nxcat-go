package main

import (
	"xyr/nexustix/nxcat/nxnet"
	"xyr/nexustix/nxcat/util"
)

func main() {
	localReceiveBuff := make(chan nxnet.Message, 8)
	localSendBuff := make(chan nxnet.Message, 8)

	//server := nxnet.NewServerTCP("0.0.0.0", "8080")
	server := nxnet.NewServerTCPxTSL("0.0.0.0", "8080")
	go server.Listen()

	writeSTDIO := util.WrapWriterSTDIO(&localSendBuff)
	readSTDIO := util.WrapReaderSTDIO(&localReceiveBuff)

	go func() {
		for {
			readSTDIO()
		}
	}()

	for {
		select {
		case msg := <-server.Messages:
			localSendBuff <- msg
			writeSTDIO()
		case msg := <-localReceiveBuff:
			server.SendMessage(msg)
		}
	}
}
