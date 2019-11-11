package nxnet

import (
	"bufio"
	"log"
	"net"

	bp "github.com/nexustix/boilerplate"
)

type Connection struct {
	conn        net.Conn
	alive       bool
	started     bool
	receivechan *chan Message
	sendchan    chan []byte
	tsize       int
	rw          *bufio.ReadWriter
	id          uint
}

func (c *Connection) setup() {
	r := bufio.NewReader(c.conn)
	w := bufio.NewWriter(c.conn)
	c.rw = bufio.NewReadWriter(r, w)
	c.onConnect()

}

func (c *Connection) onConnect() {
	// FIXME send on first message and not on connect ?
	// (prevent spam on TLS servers)
	*c.receivechan <- MakeMessage(MsgKindJoin, c.id, []byte(""))
}

func (c *Connection) Disconnect() {
	if c.alive {
		c.conn.Close()
		*c.receivechan <- MakeMessage(MsgKindLeave, c.id, []byte(""))
	}
	c.alive = false
}

func NewConnection(conn net.Conn, id uint, receivechan *chan Message) *Connection {
	connection := &Connection{
		conn:  conn,
		alive: true,
		//inchan:  make(chan []byte, 8),
		receivechan: receivechan,
		sendchan:    make(chan []byte, 8),
		tsize:       1024,
		id:          id,
		started:     false,
	}
	connection.setup()
	return connection
}

func (c *Connection) handleIncomming() {
	for {
		buff := make([]byte, c.tsize)
		n, err := c.rw.Read(buff)
		if bp.GotError(err) {
			log.Printf("<-> INFO fail socket reading >%s<", err)
			//c.alive = false
			break
		} else {
			*c.receivechan <- MakeMessage(MsgKindData, c.id, buff[0:n])
		}
	}
	c.Disconnect()
}

func (c *Connection) handleOutgoing() {
	for c.alive {
		select {
		case buff := <-c.sendchan:
			_, err := c.rw.Write(buff)
			if bp.GotError(err) {
				log.Printf("<!> INFO fail socket sending >%s<", err)
				//c.alive = false
				break
			} else {
				c.rw.Flush()
			}
		}

	}
	c.Disconnect()
}

func (c *Connection) Start() {
	if !c.started {
		go c.handleIncomming()
		go c.handleOutgoing()
		c.started = true
	}
}

func (c *Connection) IsAlive() bool {
	return c.alive
}

func (c *Connection) SendMessage(msg Message) {
	if c.alive {
		c.sendchan <- msg.Data
	}
}
