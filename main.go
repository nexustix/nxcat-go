package main

import (
	"flag"

	"github.com/nexustix/nxcat-go/nxnet"
	"github.com/nexustix/nxcat-go/util"
)

type Server interface {
	Listen()
	SendMessage(msg nxnet.Message)
	GetMessageChannel() chan nxnet.Message
}

func main() {

	hostnamePtr := flag.String("hostname", "0.0.0.0", "adress to accept connections from")
	portPtr := flag.String("port", "8080", "nxcat port")
	sslPtr := flag.Bool("ssl", false, "use ssl")

	localReceiveBuff := make(chan nxnet.Message, 8)
	localSendBuff := make(chan nxnet.Message, 8)

	flag.Parse()

	var server Server

	// fixme "inherit" server structs so
	if *sslPtr == true {
		//server = nxnet.NewServerTCPxTSL("0.0.0.0", "8080")
		server = nxnet.NewServerTCPxTSL(*hostnamePtr, *portPtr)
	} else {
		//server = nxnet.NewServerTCP("0.0.0.0", "8080")
		server = nxnet.NewServerTCP(*hostnamePtr, *portPtr)
	}
	serverMsg := server.GetMessageChannel()

	//server := nxnet.NewServerTCP("0.0.0.0", "8080")

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
		case msg := <-serverMsg:
			localSendBuff <- msg
			writeSTDIO()
		case msg := <-localReceiveBuff:
			server.SendMessage(msg)
		}
	}
}
