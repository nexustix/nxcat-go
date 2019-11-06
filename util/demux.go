package util

import (
	"bytes"
	"log"
	"strconv"
	"xyr/nexustix/nxcat/nxnet"

	bp "github.com/nexustix/boilerplate"
)

func byteEqAt(data []byte, index int, value byte) bool {
	if index >= 0 && index < len(data) {
		return data[index] == value
	}
	return false
}

func findByte(data []byte, search byte, offset int, escapable bool) int {
	i := offset

	for i < len(data) {
		if byteEqAt(data, i, search) {
			if byteEqAt(data, i+1, search) {
				i = i + 1
			} else {
				return i
			}
		}
		i = i + 1
	}
	return -1
}

type DemuxBasic struct {
	buffer   []byte
	Messages chan nxnet.Message
}

func NewDemuxBasic() *DemuxBasic {
	demux := &DemuxBasic{
		Messages: make(chan nxnet.Message, 8),
	}
	return demux
}

func (d *DemuxBasic) findMessage() bool {
	headerStart := findByte(d.buffer, '(', 0, false)
	headerEnd := -1

	if headerStart > -1 {
		headerEnd = findByte(d.buffer, ')', headerStart, true)
	}

	if headerEnd > -1 {
		header := d.buffer[headerStart+1 : headerEnd]
		segs := bytes.SplitN(header, []byte(" "), 3)
		//log.Printf("%v\n", segs)
		kind := string(segs[0])
		//id := string(segs[1])
		//FIXME not handling error
		id, err := strconv.ParseUint(string(segs[1]), 10, 32)
		if bp.GotError(err) {
			log.Printf("<-> CRITICAL failed to decode client ID >%s<", err)
		}
		//data := segs[2]
		data := bytes.ReplaceAll(segs[2], []byte("))"), []byte(")"))

		msg := nxnet.MakeMessage(kind, uint(id), data)
		//XXX potential deadlock (when channel only queried after call to FindMessage)
		//will be made private and intended to run in a go routine
		//d.Messages <- msg
		d.Messages <- msg
		//XXX comodification potential
		d.buffer = d.buffer[headerEnd+1:]
		return true
	}
	return false
}

func (d *DemuxBasic) FindMessages() {
	for d.findMessage() {
	}
}

func (d *DemuxBasic) FindMessagesForever() {
	for {
		d.FindMessages()
	}
}

func (d *DemuxBasic) DemultiplexBytes(data []byte) {
	d.buffer = append(d.buffer, data...)
}
