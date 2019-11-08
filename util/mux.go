package util

import (
	"bytes"
	"fmt"
	"xyr/nexustix/nxcat/nxnet"
)

func insertByte(slice []byte, index int, value byte) []byte {
	s := append(slice, 0 /* use the zero value of the element type */)
	copy(s[index+1:], s[index:])
	s[index] = value
	return s
}

type MuxBasic struct {
}

func NewMuxBasic() *MuxBasic {
	mux := &MuxBasic{}
	return mux
}

func (m *MuxBasic) MultiplexMessage(msg nxnet.Message) []byte {
	var buff []byte

	idString := fmt.Sprint(msg.Client_id)

	buff = append(buff, '(')
	buff = append(buff, []byte(msg.Kind)...)
	buff = append(buff, ' ')
	buff = append(buff, []byte(idString)...)
	buff = append(buff, ' ')
	buff = append(buff, bytes.ReplaceAll(msg.Data, []byte(")"), []byte("))"))...)
	buff = append(buff, ')')

	return buff
}
