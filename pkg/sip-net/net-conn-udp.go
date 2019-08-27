package sip_net

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net"
	"sip/pkg/sip"
)

type UdpConn struct {
	address string
	udpConn *net.UDPConn
	buf     []byte
	id      string
}

func (conn UdpConn) Type() string {
	return "udp"
}

func (conn UdpConn) Id() string {
	return conn.id
}

func (conn UdpConn) Receive() (sip.SipMessage, string) {
	udp := conn.udpConn
	buf := conn.buf
	i, _, _, addr, err := udp.ReadMsgUDP(buf, nil)
	if err != nil {
		log.Println(err)
	}
	msg, _ := sip.Marshal(bytes.NewReader(buf[:i]))
	return msg, addr.String()
}

func (conn UdpConn) Send(message sip.SipMessage, address string) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	str := message.String()
	fmt.Println("发送数据: ->", address)
	fmt.Println(str)
	conn.udpConn.WriteToUDP([]byte(str), addr)
}

func (conn *UdpConn) Start() {
	conn.id = uuid.NewV1().String()
	addr, _ := net.ResolveUDPAddr("udp", conn.address)
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
	} else {
		conn.udpConn = udpConn
	}
}
