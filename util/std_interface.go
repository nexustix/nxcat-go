package util

import (
	"bufio"
	"log"
	"os"

	"github.com/nexustix/nxcat-go/nxnet"

	bp "github.com/nexustix/boilerplate"
)

func WrapWriterSTDIO(pipe *chan nxnet.Message) func() {
	h := bufio.NewWriter(os.Stdout)
	mux := NewMuxBasic()
	return func() {
		select {
		case msg := <-*pipe:
			mMsg := mux.MultiplexMessage(msg)
			_, err := h.Write(mMsg)
			//log.Printf("<wrt< %s \n", mMsg)
			h.Flush()
			if bp.GotError(err) {
				log.Printf("<!> CRITICAL fail writing to STDOUT >%s<", err)
			}
		}
	}
}

func WrapReaderSTDIO(pipe *chan nxnet.Message) func() {
	h := bufio.NewReader(os.Stdin)
	demux := NewDemuxBasic()
	return func() {
		buff := make([]byte, 1024)
		n, err := h.Read(buff)
		//log.Printf(">red> %s \n", buff[0:n])
		if bp.GotError(err) {
			log.Printf("<!> CRITICAL fail STDIN reading >%s<", err)
		} else {
			demux.DemultiplexBytes(buff[0:n])
			demux.FindMessages()
		}
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
