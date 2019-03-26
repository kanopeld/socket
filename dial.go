package socket

import (
	"bufio"
	"net"
	"sync"
)

type dial struct {
	*clientHandler
	*defaultEmitter
	conn net.Conn
	id string
	on bool
	disc bool
}

func NewDial(addr string) (Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	d := &dial{
		conn:conn,
		defaultEmitter:&defaultEmitter{c:conn},
	}
	d.clientHandler = newClientHandler(d, &baseHandler{events:make(map[string]*caller, 0), name:BASE_HANDLER_DEFAULT_NAME, evMu:sync.Mutex{}})
	go d.loop()
	return d, nil
}

func (d *dial) ID() string {
	return d.id
}

func (d *dial) Connection() net.Conn {
	return d.conn
}

func (d *dial) sendConnect() {
	p := NewPacket(_PACKET_TYPE_CONNECT)
	_, _ = d.conn.Write(p.MarshalBinary())
}

func (d *dial) Disconnect() {
	d.on = false
	if d.disc {
		return
	}
	d.disc = true
	_ = d.send(&Package{PT:_PACKET_TYPE_DISCONNECT})
	_ = d.call(DISCONNECTION_NAME, nil)
	_ = d.conn.Close()
}

func (d *dial) Broadcast(event string, msg []byte) error {
	return nil
}

func (d *dial) loop() {
	d.on = true
	defer d.Disconnect()
	d.sendConnect()

	for d.on {
		msg, err := bufio.NewReader(d.conn).ReadBytes('\n')
		if err != nil {
			return
		}

		p, err := DecodePackage(msg)
		if err != nil {
			return
		}

		switch p.PT {
		case _PACKET_TYPE_CONNECT:
			d.id = string(p.Payload)
			if err := d.call(CONNECTION_NAME, nil); err != nil {
				return
			}
		case _PACKET_TYPE_DISCONNECT:
			return
		case _PACKET_TYPE_EVENT:
			msg ,err := DecodeMessage(p.Payload)
			if err != nil {
				return
			}

			if err := d.call(msg.EventName, msg.Data); err != nil {
				return
			}
		}
	}
}
