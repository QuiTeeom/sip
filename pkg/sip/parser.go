package sip

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var reg_start_line = regexp.MustCompile(`(SIP/2.0|.*)\s+(.*)\s+(SIP/2.0|.*)`)

func Marshal(reader io.Reader) (SipMessage, error) {
	bufferReader := bufio.NewReader(reader)
	startLine, _, _ := bufferReader.ReadLine()
	startLiveValue := reg_start_line.FindStringSubmatch(string(startLine))
	if len(startLiveValue) != 4 {
		return nil, nil
	}

	headerStr := ""
	for true {
		lineBuf, _, _ := bufferReader.ReadLine()
		if lineBuf == nil || len(lineBuf) == 0 {
			break
		}
		line := string(lineBuf)
		headerStr = headerStr + line + "\r\n"
	}
	headers, _ := phraseSipHeader(headerStr)

	contentLengthHeader := headers.Get(HEADER_Content_Length_U)
	contentLength := 0
	var body []byte

	if contentLengthHeader != nil {
		contentLength = headers.Get(HEADER_Content_Length_U).(IntHeader).Value
	}

	if contentLength > 0 {
		content := make([]byte, contentLength)
		_, err := bufferReader.Read(content)
		if err != nil {
			fmt.Println(err)
		}
		body = content
	} else {
		body = make([]byte, 0)
	}

	var msg SipMessage
	if startLiveValue[1] == "SIP/2.0" {
		code, _ := strconv.Atoi(startLiveValue[2])
		msg = ResponseMessage{
			Code:         code,
			ReasonPhrase: startLiveValue[3],
			Headers:      headers,
			Body:         body,
		}
	} else {
		msg = RequestMessage{
			Method:  startLiveValue[1],
			Uri:     uriParser(startLiveValue[2]),
			Version: startLiveValue[3],
			Headers: headers,
			Body:    body,
		}
	}
	return msg, nil
}

//解析sip头
func phraseSipHeader(message string) (SipHeaders, error) {
	lines := strings.Split(message, "\r\n")
	sipHeaders := SipHeaders{
		headers: make(map[string][]SipHeader),
	}
	aHeaderLine := ""
	for _, line := range lines {
		//多行的 header-value
		if strings.Index(line, " ") == 0 || strings.Index(line, "\t") == 0 {
			line = strings.TrimLeft(line, " ")
			line = strings.TrimLeft(line, "\t")
			aHeaderLine = aHeaderLine + " " + line
			continue
		}
		if aHeaderLine != "" {
			headers := parseHeadLine(aHeaderLine)
			sipHeaders.AddHeaders(headers)
		}
		aHeaderLine = line
	}
	return sipHeaders, nil
}

var headerNeedToDivide = map[string]bool{
	HEADER_VIA_U: true,
}

func parseHeadLine(line string) []SipHeader {
	headerPair := strings.SplitN(line, ":", 2)
	headerName := strings.Trim(headerPair[0], " ")
	headerNameUpper := strings.ToUpper(headerName)
	headerValue := strings.Trim(headerPair[1], " ")
	if headerNeedToDivide[headerNameUpper] {
		headerValues := strings.Split(headerValue, ",")
		res := make([]SipHeader, len(headerValues))
		for i, singleValue := range headerValues {
			res[i] = parseOneHeadValue(headerName, singleValue)
		}
		return res
	} else {
		return []SipHeader{parseOneHeadValue(headerName, headerValue)}
	}
}

func parseOneHeadValue(name, value string) SipHeader {
	parser := parserMap[strings.ToUpper(name)]
	if parser == nil {
		return commonHeaderParser(name, value)
	} else {
		return parser(name, value)
	}
}

var parserMap = map[string]func(name, value string) SipHeader{
	HEADER_From_U:           nameUriHeaderParser,
	HEADER_To_U:             nameUriHeaderParser,
	HEADER_VIA_U:            viaHeaderParser,
	HEADER_Call_Id_U:        commonHeaderParser,
	HEADER_CSeq_U:           commonHeaderParser,
	HEADER_Contact_U:        nameUriHeaderParser,
	HEADER_Content_Type_U:   commonHeaderParser,
	HEADER_Content_Length_U: intHeaderParser,
	HEADER_Expires_U:        intHeaderParser,
}

var reg_name_uri1 = regexp.MustCompile(`\s*(.*\s|.*?)\s*(?:<)(.*)(?:>)\s*(;.*|$)`)
var reg_name_uri2 = regexp.MustCompile(`\s*(.*\s|.*?)\s*(.*?)\s*(;.*|$)`)

func nameUriHeaderParser(name, value string) SipHeader {
	valueGroup := reg_name_uri1.FindStringSubmatch(value)
	if valueGroup == nil {
		valueGroup = reg_name_uri2.FindStringSubmatch(value)
	}
	return NameUriHeader{
		Name:        name,
		DisplayName: strings.Trim(valueGroup[1], " "),
		Uri:         uriParser(valueGroup[2]),
		Parameters:  stringToParameters(valueGroup[3]),
	}
}

// SIP/2.0/UDP r9ik7446ccqm.invalid;branch=z9hG4bK5075365
var reg_via = regexp.MustCompile(`(.*)(?:/)(.*?)\s+(.*?)(;.*|$)`)

func viaHeaderParser(name, value string) SipHeader {
	valueGroup := reg_via.FindStringSubmatch(value)
	return ViaHeader{
		Version:    valueGroup[1],
		Protocol:   valueGroup[2],
		Address:    valueGroup[3],
		Parameters: stringToParameters(valueGroup[4]),
	}
}

//<sip:watson@worcester.bell-telephone.com> ;q=0.7; expires=3600")
var reg_uri = regexp.MustCompile(`(.*?):([^@]*)@?(.*?|.*?)\s*(;.*|$)`)

func uriParser(uriStr string) Uri {
	uri := reg_uri.FindStringSubmatch(uriStr)
	if uri[3] == "" {
		uri[2], uri[3] = uri[3], uri[2]
	}
	return Uri{
		Protocol:   uri[1],
		User:       uri[2],
		Host:       uri[3],
		Parameters: stringToParameters(uri[4]),
	}
}

var reg_int_header = reg_common_header

func intHeaderParser(name, value string) SipHeader {
	valueGroup := reg_int_header.FindStringSubmatch(value)
	v, _ := strconv.Atoi(valueGroup[1])

	return IntHeader{
		Name:       name,
		Value:      v,
		Parameters: stringToParameters(valueGroup[2]),
	}
}

var reg_common_header = regexp.MustCompile(`\s*(.*?)\s*(;.*|$)`)

func commonHeaderParser(name, value string) SipHeader {
	valueGroup := reg_common_header.FindStringSubmatch(value)
	return CommonHeader{
		Name:       name,
		Value:      valueGroup[1],
		Parameters: stringToParameters(valueGroup[2]),
	}
}
