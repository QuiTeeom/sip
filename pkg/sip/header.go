package sip

import (
	"fmt"
	"strconv"
	"strings"
)

type SipHeader interface {
	String() string
	GetName() string
	GetValue() string
	GetParameters() Parameters
}

type SipHeaders struct {
	headers    map[string][]SipHeader
	linkedName []string
}

func (header SipHeaders) String() string {
	res := ""
	for _, n := range header.linkedName {
		for _, v := range header.headers[n] {
			res += v.String() + "\r\n"
		}
	}
	return res
}

func (headers *SipHeaders) AddHeader(header ...SipHeader) SipHeaders {
	if headers.headers == nil {
		headers.headers = make(map[string][]SipHeader)
	}
	for _, h := range header {
		n := strings.ToUpper(h.GetName())
		hs := headers.headers[n]
		if hs == nil {
			headers.headers[n] = make([]SipHeader, 0)
			hs = headers.headers[n]
			headers.linkedName = append(headers.linkedName, n)
		}
		headers.headers[n] = append(hs, h)
	}

	return *headers
}

func (headers *SipHeaders) AddHeaders(header []SipHeader) SipHeaders {
	for _, h := range header {
		headers.AddHeader(h)
	}
	return *headers
}

func (header SipHeaders) Get(name string) SipHeader {
	hs := header.GetM(name)
	if hs == nil || len(hs) == 0 {
		return nil
	}
	return hs[0]
}

func (header SipHeaders) GetM(name string) []SipHeader {
	return header.headers[strings.ToUpper(name)]
}

//parameter
type Parameters map[string]string

func (p Parameters) String() string {
	return parametersToString(p)
}

func (p Parameters) Get(key string) string {
	return p[key]
}

func (p Parameters) GetInt(key string) int {
	str := p[key]
	if str == "" {
		return -1
	} else {
		res, _ := strconv.Atoi(str)
		return res
	}
}

func newParameters(p ...[]map[string]string) Parameters {
	if len(p) == 1 && len(p[0]) == 1 {
		return Parameters(p[0][0])
	} else {
		return Parameters(make(map[string]string))
	}
}

//A common header  head-name:head-Value(;head-param1:head-param-value1 ...)
type CommonHeader struct {
	Name       string
	Value      string
	Parameters Parameters
}

func (c CommonHeader) String() string {
	return fmt.Sprintf("%s: %s%s", c.Name, c.Value, c.Parameters.String())
}

func (c CommonHeader) GetName() string {
	return c.Name
}

func (c CommonHeader) GetValue() string {
	return c.Value
}

func (c CommonHeader) GetParameters() Parameters {
	return c.Parameters
}

var Header = func(Name string, Value string, Parameters ...map[string]string) CommonHeader {
	return CommonHeader{
		Name: Name, Value: Value, Parameters: newParameters(Parameters),
	}
}

func newHeaderMaker(name string) func(Value string, Parameters ...map[string]string) CommonHeader {
	return func(Value string, Parameters ...map[string]string) CommonHeader {
		return CommonHeader{
			Name: name, Value: Value, Parameters: newParameters(Parameters),
		}
	}
}

//URI  sip:(user@)host(;param1:value1;param1:value1 ...)
type Uri struct {
	Protocol   string     //required
	User       string     //optional
	Host       string     //required
	Parameters Parameters //optional
}

func (uri Uri) String() string {
	res := fmt.Sprintf("%s:%s%s", uri.Protocol, uri.Id(), parametersToString(uri.Parameters))
	return res
}

func (uri Uri) Id() string {
	if uri.User == "" {
		return uri.Host
	} else {
		return uri.User + "@" + uri.Host
	}
}

//Header Like : (displayName) (<)uri(>) (;param1:value1;param1:value1 ...)
type NameUriHeader struct {
	Name        string
	DisplayName string
	Uri         Uri
	Parameters  Parameters
}

func (header NameUriHeader) String() string {
	if header.DisplayName == "" {
		return fmt.Sprintf("%s: <%s>%s", header.Name, header.Uri.String(), header.Parameters.String())
	}
	return fmt.Sprintf("%s: %s <%s>%s", header.Name, header.DisplayName, header.Uri.String(), header.Parameters.String())
}

func (header NameUriHeader) GetName() string {
	return header.Name
}

func (header NameUriHeader) GetValue() string {
	return fmt.Sprintf("%s <%s>", header.DisplayName, header.Uri.String())
}

func (header NameUriHeader) GetParameters() Parameters {
	return header.Parameters
}

func newNameUriHeaderMaker(name string) func(displayName string, uri Uri, parameters ...map[string]string) NameUriHeader {
	return func(displayName string, uri Uri, parameters ...map[string]string) NameUriHeader {
		return NameUriHeader{
			DisplayName: displayName,
			Uri:         uri,
			Parameters:  newParameters(parameters),
			Name:        name,
		}
	}
}

//int header
type IntHeader struct {
	Name       string
	Value      int
	Parameters Parameters
}

func (c IntHeader) GetValue() string {
	return strconv.Itoa(c.Value)
}

func (c IntHeader) String() string {
	return fmt.Sprintf("%s: %d%s", c.Name, c.Value, c.Parameters.String())
}

func (c IntHeader) GetName() string {
	return c.Name
}

func (c IntHeader) GetIntValue() int {
	return c.Value
}

func (c IntHeader) GetParameters() Parameters {
	return c.Parameters
}

func newIntHeaderMaker(name string) func(name string, value int, parameters ...map[string]string) IntHeader {
	return func(name string, value int, parameters ...map[string]string) IntHeader {
		return IntHeader{
			Name:       name,
			Value:      value,
			Parameters: newParameters(parameters),
		}
	}
}

//Via
type ViaHeader struct {
	Version    string
	Protocol   string
	Address    string
	Parameters Parameters
}

func (header ViaHeader) String() string {
	return fmt.Sprintf("%s: %s%s", header.GetName(), header.GetValue(), header.GetParameters().String())
}

func (header ViaHeader) GetName() string {
	return HEADER_VIA
}

func (header ViaHeader) GetValue() string {
	return fmt.Sprintf("%s/%s %s", header.Version, header.Protocol, header.Address)
}

func (header ViaHeader) GetParameters() Parameters {
	return header.Parameters
}

const HEADER_VIA = "Via"

var HEADER_VIA_U = strings.ToUpper(HEADER_VIA)
var Via = func(version, protocol, address string, parameters ...map[string]string) ViaHeader {
	return ViaHeader{
		Version:    version,
		Protocol:   protocol,
		Address:    address,
		Parameters: newParameters(parameters),
	}
}

//From
const HEADER_From = "From"

var HEADER_From_U = strings.ToUpper(HEADER_From)
var From = newNameUriHeaderMaker(HEADER_From)

//To
const HEADER_To = "To"

var HEADER_To_U = strings.ToUpper(HEADER_To)
var To = newNameUriHeaderMaker(HEADER_To)

//Contact
const HEADER_Contact = "Contact"

var HEADER_Contact_U = strings.ToUpper(HEADER_Contact)
var Contact = newNameUriHeaderMaker(HEADER_Contact)

//Call-Id
const HEADER_Call_Id = "Call-Id"

var HEADER_Call_Id_U = strings.ToUpper(HEADER_Call_Id)
var CallId = newHeaderMaker(HEADER_Call_Id)

//CSeq
const HEADER_CSeq = "CSeq"

var HEADER_CSeq_U = strings.ToUpper(HEADER_CSeq)
var CSeq = newHeaderMaker(HEADER_CSeq_U)

//Content-Type
const HEADER_Content_Type = "Content-Type"

var HEADER_Content_Type_U = strings.ToUpper(HEADER_Content_Type)
var ContentType = newHeaderMaker(HEADER_Content_Type)

//Content-Length
const HEADER_Content_Length = "Content-Length"

var HEADER_Content_Length_U = strings.ToUpper(HEADER_Content_Length)
var ContentLength = newIntHeaderMaker(HEADER_Content_Length)

//Content-Length
const HEADER_Expires = "Expires"

var HEADER_Expires_U = strings.ToUpper(HEADER_Expires)
var Expires = newIntHeaderMaker(HEADER_Expires)
