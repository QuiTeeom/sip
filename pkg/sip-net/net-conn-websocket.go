package sip_net

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"sip/pkg/sip"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"sip"},
}

type WebsocketConn struct {
	address string
	msgChan chan messageAddressBundle
	connMap map[string]*websocket.Conn
	id      string
}

func (conn WebsocketConn) Id() string {
	return conn.id
}
func (WebsocketConn) Type() string {
	return "ws"
}
func (conn WebsocketConn) Receive() (sip.SipMessage, string) {
	messageAddressBundle := <-conn.msgChan
	return messageAddressBundle.message, messageAddressBundle.address
}

func (conn WebsocketConn) Send(message sip.SipMessage, address string) {
	str := message.String()
	fmt.Println("发送数据: ->", address)
	fmt.Println(str)
	c := conn.connMap[address]
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(str))
	}
}

func (conn *WebsocketConn) Start() {
	conn.msgChan = make(chan messageAddressBundle)
	conn.connMap = make(map[string]*websocket.Conn)
	conn.id = uuid.NewV1().String()
	http.HandleFunc(conn.address, conn.sipServer)
}

func (conn *WebsocketConn) sipServer(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	id := uuid.NewV4().String()
	conn.connMap[id] = c
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		sip, _ := sip.Marshal(bytes.NewReader(message))
		conn.msgChan <- messageAddressBundle{
			message: sip,
			address: id,
		}
	}
}

type messageAddressBundle struct {
	message sip.SipMessage
	address string
}
