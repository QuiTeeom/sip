package sip_server

import (
	"sip/pkg/sip"
	"sip/pkg/sip-net"
)

type SipServer struct {
	NetOptions     []NetOption
	MessageHandler func(messageEvent MessageEvent)

	//private
	msgChan chan MessageEvent

	conns []sip_net.Conn
}

func (server *SipServer) Start() {
	prepareSipServer(server)
	for i, net := range server.NetOptions {
		conn := sip_net.Listen(net.Type, net.Address)
		server.conns[i] = conn
		go server.listen(conn)
	}

	go server.onMessage()
}

func (server SipServer) listen(conn sip_net.Conn) {
	for {
		msg, address := conn.Receive()
		if msg != nil {
			server.msgChan <- MessageEvent{
				Message: msg,
				ContactAddress: ContactAddress{
					TransportType: conn.Type(),
					TransportId:   conn.Id(),
					Address:       address,
				},
			}
		}
	}
}

func (server SipServer) onMessage() {
	for {
		select {
		case event := <-server.msgChan:
			if server.MessageHandler != nil {
				server.MessageHandler(event)
			}
		}
	}
}

func (server SipServer) Send(message sip.SipMessage, contactAddress ContactAddress) {
	for _, c := range server.conns {
		if c.Type() == contactAddress.TransportType && c.Id() == contactAddress.TransportId {
			c.Send(message, contactAddress.Address)
			break
		}
	}

}

func prepareSipServer(server *SipServer) {
	if server.NetOptions == nil {
		server.NetOptions = make([]NetOption, 0)
	}
	if len(server.NetOptions) == 0 {
		server.NetOptions = append(server.NetOptions, NetOption{
			Type: "udp", Address: ":45060",
		})
	}

	if server.msgChan == nil {
		server.msgChan = make(chan MessageEvent)
	}

	if server.MessageHandler == nil {
		server.MessageHandler = nil
	}

	server.conns = make([]sip_net.Conn, len(server.NetOptions))
}

type NetOption struct {
	Type    string
	Address string
}

type MessageEvent struct {
	Message        sip.SipMessage
	ContactAddress ContactAddress
}

type ContactAddress struct {
	TransportType string
	TransportId   string
	Address       string
}
