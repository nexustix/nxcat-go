package nxnet

import (
	"fmt"
	"log"
	"net"
)

type ServerTCP struct {
	hostname    string
	port        string
	idCounter   uint
	connections map[uint]*Connection
	Messages    chan Message
}

func NewServerTCP(hostname string, port string) *ServerTCP {
	server := &ServerTCP{
		hostname:    hostname,
		port:        port,
		idCounter:   1,
		connections: make(map[uint]*Connection),
		Messages:    make(chan Message, 8),
	}
	return server
}

func (s *ServerTCP) GetMessageChannel() chan Message {
	return s.Messages
}

func (s *ServerTCP) Listen() {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.hostname, s.port))
	if err != nil {
		log.Fatal("<!> ERR starting tcp server")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		if err != nil {
			log.Printf("<-> WARN fail accepting: %s\n", err)
		} else {
			connection := NewConnection(conn, s.idCounter, &s.Messages)
			s.connections[s.idCounter] = connection
			s.connections[s.idCounter].Start()
			s.idCounter = s.idCounter + 1
		}
	}
}

func (s *ServerTCP) SendMessage(msg Message) {
	if msg.Client_id == 0 {
		switch msg.Kind {
		case MsgKindData:
			for _, v := range s.connections {
				v.SendMessage(msg)
			}
		}

	} else {
		// TODO merge if statements ?
		if val, ok := s.connections[msg.Client_id]; ok {
			if val.IsAlive() {
				switch msg.Kind {
				case MsgKindData:
					val.SendMessage(msg)
				}
			} else {
				log.Printf("<!> WARN trying to send to dead connection: %v \n", msg.Client_id)
			}
		} else {
			log.Printf("<!> WARN trying to send to invalid connection: %v \n", msg.Client_id)
		}
	}
}

// values for future tests
//Aloha
//the quick brown fox jumps over the lazy dog
//hello)
//hello ()
//hello ())
//this is ((testing))) ()))
//(testing)) ())
