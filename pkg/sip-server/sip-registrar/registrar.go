package sip_registrar

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"sip/pkg/sip"
	"sip/pkg/sip-server"
	"strconv"
	"strings"
	"time"
)

type Registrar struct {
	SipServer *sip_server.SipServer
	Dao       *gorm.DB
	Domain    string
}

func (r *Registrar) Start() {
	r.Dao = NewDao()

	r.SipServer.MessageHandler = r.MessageHandler
	r.SipServer.Start()
}

func (r Registrar) MessageHandler(event sip_server.MessageEvent) {
	fmt.Println("收到Sip消息")
	fmt.Println(event.Message.String())

	message := event.Message

	if req, ok := message.(sip.RequestMessage); ok {
		switch string(req.Method) {
		case sip.Method_Register:
			r.onRegister(req, event.ContactAddress)
			break
		}
	}
}

func (r Registrar) onRegister(message sip.RequestMessage, originAddress sip_server.ContactAddress) {
	resp := sip.ResponseMessage{
		Code:         200,
		ReasonPhrase: "OK",
		Headers:      sip.SipHeaders{},
	}

	// 判断 首行 uri 是否正确

	//从 To 中，提取注册的地址
	to := message.Headers.Get(sip.HEADER_To).(sip.NameUriHeader)
	address := to.Uri.Id()
	//Call Id
	callId := message.Headers.Get(sip.HEADER_Call_Id).GetValue()

	//整体的expires
	expire := 600
	expireHeader := message.Headers.Get(sip.HEADER_Expires)
	if expireHeader != nil {
		expire = expireHeader.(sip.IntHeader).Value
	}

	//Cseq
	v := message.Headers.Get(sip.HEADER_CSeq).GetValue()
	cseq, _ := strconv.Atoi(strings.Split(v, " ")[0])

	//contact
	contacts := message.Headers.GetM(sip.HEADER_Contact)
	for _, contact := range contacts {
		header := contact.(sip.NameUriHeader)

		ex := header.Parameters.GetInt("expires")
		if ex != -1 {
			expire = ex
		}

		sipInstance := header.Parameters.Get("+sip.instance")
		if sipInstance != "" {
			sipInstance = strings.Trim(sipInstance, " ")
			sipInstance = strings.Trim(sipInstance, "\"")
			sipInstance = strings.Trim(sipInstance, "<")
			sipInstance = strings.Trim(sipInstance, ">")
		}

		contactInfo := ContactInfo{
			Address:     address,
			CallId:      callId,
			ExpireDate:  time.Now().Add(time.Second * time.Duration(expire)),
			Cseq:        cseq,
			Contact:     header.Uri.String(),
			GlobalRoute: sipInstance,
		}

		r.Dao.Create(&contactInfo)
	}

	resp.Headers.AddHeader(
		message.Headers.Get(sip.HEADER_VIA),
		message.Headers.Get(sip.HEADER_From),
		sip.To(to.DisplayName, to.Uri, sip.Parameters{"tag": uuid.NewV1().String()}),
		message.Headers.Get(sip.HEADER_Call_Id),
		message.Headers.Get(sip.HEADER_CSeq),
	)

	r.SipServer.Send(resp, originAddress)
}
