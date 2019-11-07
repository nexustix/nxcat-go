package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"xyr/nexustix/nxcat/nxnet"
	"xyr/nexustix/nxcat/util"

	bp "github.com/nexustix/boilerplate"
)

/*
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
*/

func wrapWriterSTDIO(pipe *chan nxnet.Message) func() {
	h := bufio.NewWriter(os.Stdout)
	mux := util.NewMuxBasic()
	return func() {
	loop:
		for {
			select {
			case msg := <-*pipe:
				//log.Printf("DING!\n")
				mMsg := mux.MultiplexMessage(msg)
				_, err := h.Write(mMsg)
				h.Flush()
				if bp.GotError(err) {
					log.Printf("<!> CRITICAL fail writing to STDOUT >%s<", err)
				}
			default:
				//return
				break loop
			}
		}
	}
}

func readMessageSTDIN(pipe *chan nxnet.Message, reader io.Reader) {
	demux := util.NewDemuxBasic()
	for {
		buff := make([]byte, 1024)
		n, err := reader.Read(buff)
		if bp.GotError(err) {
			log.Printf("<!> CRITICAL fail STDIN reading >%s<", err)
		} else {
			demux.DemultiplexBytes(buff[0:n])
			demux.FindMessages()
			//*pipe <- MakeMessage(MsgKindData, c.id, buff[0:n])
		cake:
			for {
				select {
				case msg := <-demux.Messages:
					*pipe <- msg
				default:
					break cake
				}
			}
		}
	}
}

func main() {
	localReceiveBuff := make(chan nxnet.Message, 8)
	localSendBuff := make(chan nxnet.Message, 8)
	//stdioWriter := bufio.NewWriter(os.Stdout)

	//server := nxnet.NewServerTCP("0.0.0.0", "8080")
	server := nxnet.NewServerTCPxTSL("0.0.0.0", "8080")
	//mux := util.NewMuxBasic()
	//demux := util.NewDemuxBasic()

	go server.Listen()
	//go demux.FindMessagesForever()
	//go writeMessageSTDOUT(&demux.Messages, stdioWriter)
	writeSTDIO := wrapWriterSTDIO(&localSendBuff)
	go readMessageSTDIN(&localReceiveBuff, os.Stdin)

	//for cheese := range localReceiveBuff {
	//	log.Printf("CAKE! >%v<\n", cheese)
	//}

	for {
		//log.Printf("tick\n")
		select {
		case msg := <-server.Messages:
			//fmt.Printf(">%v<", msg.Data)
			/*
				log.Printf("msg%v<\n", msg)
				mMsg := mux.MultiplexMessage(msg)
				log.Printf("muxed>%s<\n", mMsg)

				demux.DemultiplexBytes(mMsg)
				demux.FindMessages()

				writeSTDIO()
			*/
			localSendBuff <- msg
			writeSTDIO()
		//default:
		//	time.Sleep(100 * time.Millisecond)
		case msg := <-localReceiveBuff:
			server.SendMessage(msg)
		}
		//log.Printf("tock\n")

		/*
			msg := <-server.Messages
			//fmt.Printf(">%v<", msg.Data)
			log.Printf("msg%v<\n", msg)
			mMsg := mux.MultiplexMessage(msg)
			log.Printf("muxed>%s<\n", mMsg)

			demux.DemultiplexBytes(mMsg)
			demux.FindMessages()

			writeSTDIO()
		*/
	}
}

/*
(dat 1 dfasdf)
*/
