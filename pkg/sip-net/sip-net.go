package sip_net

import "sip/pkg/sip"

const DEFAULT_UDP_BUFFER_SIZE = 64 * 1024

//Conn
type Conn interface {
	Receive() (sip.SipMessage, string)
	Send(message sip.SipMessage, address string)
	Id() string
	Type() string
}

//start listen and return conn
func Listen(netType, address string) Conn {
	switch netType {
	case "udp":
		return startListenUdp(address)
	case "ws":
		return startListenWebrtc(address)
	}
	return nil
}

func startListenUdp(address string) Conn {
	conn := UdpConn{
		address: address,
		buf:     make([]byte, DEFAULT_UDP_BUFFER_SIZE),
	}
	conn.Start()
	return conn
}

func startListenWebrtc(address string) Conn {
	conn := WebsocketConn{
		address: address,
	}
	conn.Start()
	return conn
}
