package nxnet

const (
	MsgKindData  string = "dat"
	MsgKindJoin  string = "new"
	MsgKindLeave string = "gon"
)

type Message struct {
	Kind      string
	Client_id uint
	Data      []byte
}

func MakeMessage(kind string, client_id uint, data []byte) Message {
	msg := Message{
		Kind:      kind,
		Client_id: client_id,
		Data:      data,
	}
	return msg
}

func NewMessage(kind string, client_id uint, data []byte) *Message {
	msg := MakeMessage(kind, client_id, data)
	//msg := &Message{
	//	kind:      kind,
	//	client_id: client_id,
	//	data:      data,
	//}
	return &msg
}
