package sip

import "fmt"

type SipMessage interface {
	String() string
	GetHeaders() SipHeaders
	GetType() string
}

type RequestMessage struct {
	Method  string
	Uri     Uri
	Version string
	Headers SipHeaders
	Body    []byte
}

func (r RequestMessage) GetType() string {
	return "request"
}

func (r RequestMessage) GetHeaders() SipHeaders {
	return r.Headers
}

func (r RequestMessage) String() string {
	return fmt.Sprintf("%s %s %s\r\n%s\r\n%s", r.Method, r.Uri.String(), r.Version, r.Headers.String(), string(r.Body))
}

type ResponseMessage struct {
	Code         int
	ReasonPhrase string
	Headers      SipHeaders
	Body         []byte
}

func (r ResponseMessage) GetHeaders() SipHeaders {
	return r.Headers
}

func (r ResponseMessage) String() string {
	return fmt.Sprintf("SIP/2.0 %d %s\r\n%s\r\n%s", r.Code, r.ReasonPhrase, r.Headers.String(), string(r.Body))
}

func (r ResponseMessage) GetType() string {
	return "response"
}

const (
	Method_Register string = "REGISTER"
)
