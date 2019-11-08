package nxnet

import (
	"crypto/tls"
	"fmt"
	"log"
)

type ServerTCPxTSL struct {
	hostname    string
	port        string
	idCounter   uint
	connections map[uint]*Connection
	Messages    chan Message
}

func NewServerTCPxTSL(hostname string, port string) *ServerTCPxTSL {
	server := &ServerTCPxTSL{
		hostname:    hostname,
		port:        port,
		idCounter:   1,
		connections: make(map[uint]*Connection),
		Messages:    make(chan Message, 8),
	}
	return server
}

func (s *ServerTCPxTSL) Listen() {
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("<!> ERR loading cerificates: %s\n", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err := tls.Listen("tcp", fmt.Sprintf("%s:%s", s.hostname, s.port), config)
	if err != nil {
		log.Fatal("<!> ERR starting tcp server")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
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

func (s *ServerTCPxTSL) SendMessage(msg Message) {
	if msg.Client_id == 0 {
		switch msg.Kind {
		case MsgKindData:
			for i, v := range s.connections {
				if i != 0 {
					v.SendMessage(msg)
				}
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
