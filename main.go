package main

import (
	"bufio"
	"log"
	"os"
	"xyr/nexustix/nxcat/nxnet"
	"xyr/nexustix/nxcat/util"

	bp "github.com/nexustix/boilerplate"
)

func readMessageSTDIN(pipe *chan<- nxnet.Message) {

}

func writeMessageSTDOUT(pipe *chan nxnet.Message) {
	h := bufio.NewWriter(os.Stdout)
	mux := util.NewMuxBasic()
	//for msg := range *pipe {
	//	log.Printf("DING!\n")
	//	mMsg := mux.MultiplexMessage(msg)
	//	_, err := h.Write(mMsg)
	//	if bp.GotError(err) {
	//		log.Printf("<!> CRITICAL fail writing to STDOUT >%s<", err)
	//	}
	//}

	for {

		select {
		case msg := <-*pipe:
			log.Printf("DING!\n")
			mMsg := mux.MultiplexMessage(msg)
			_, err := h.Write(mMsg)
			h.Flush()
			if bp.GotError(err) {
				log.Printf("<!> CRITICAL fail writing to STDOUT >%s<", err)
			}
		default:
			return
		}
	}

}

func wrapWriterSTDIO(pipe *chan nxnet.Message) func() {
	h := bufio.NewWriter(os.Stdout)
	mux := util.NewMuxBasic()
	return func() {
		for {

			select {
			case msg := <-*pipe:
				log.Printf("DING!\n")
				mMsg := mux.MultiplexMessage(msg)
				_, err := h.Write(mMsg)
				h.Flush()
				if bp.GotError(err) {
					log.Printf("<!> CRITICAL fail writing to STDOUT >%s<", err)
				}
			default:
				return
			}
		}
	}
}

func main() {
	//localReceiveBuff := make(chan nxnet.Message, 8)
	//localSendBuff := make(chan nxnet.Message, 8)
	//stdioWriter := bufio.NewWriter(os.Stdout)

	server := nxnet.NewServerTCP("0.0.0.0", "8080")
	mux := util.NewMuxBasic()
	demux := util.NewDemuxBasic()

	go server.Listen()
	//go demux.FindMessagesForever()
	//go writeMessageSTDOUT(&demux.Messages, stdioWriter)
	writeSTDIO := wrapWriterSTDIO(&demux.Messages)

	for {
		msg := <-server.Messages
		//fmt.Printf(">%v<", msg.Data)
		log.Printf("msg%v<\n", msg)
		mMsg := mux.MultiplexMessage(msg)
		log.Printf("muxed>%s<\n", mMsg)

		demux.DemultiplexBytes(mMsg)
		demux.FindMessages()

		//select {
		//case xmsg := <-demux.Messages:
		//	log.Printf("xmsg%v<\n", xmsg)
		//default:
		//}
		//writeMessageSTDOUT(&demux.Messages)
		writeSTDIO()
	}
}
