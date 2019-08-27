package sip

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	TT := From(
		"22qq", Uri{
			Protocol: "sip",
			Host:     "quitee.com",
		})

	fmt.Println(TT.String())
}

func TestParser(t *testing.T) {
	nameUriHeaderParser("Contact", "<sip:qan5ambg@a2henqvpmmjq.invalid;transport=ws>;+sip.ice;reg-id=1;+sip.instance=\"<urn:uuid:842c78fe-bb1b-4a34-ae6b-1ce69d050748>\";expires=600")
	nameUriHeaderParser("Contact", "asdad <sip:qan5ambg@a2henqvpmmjq.invalid;transport=ws>;+sip.ice;reg-id=1;+sip.instance=\"<urn:uuid:842c78fe-bb1b-4a34-ae6b-1ce69d050748>\";expires=600")
	nameUriHeaderParser("Contact", "asdad sip:qan5ambg@a2henqvpmmjq.invalid;transport=ws;+sip.ice;reg-id=1;+sip.instance=\"<urn:uuid:842c78fe-bb1b-4a34-ae6b-1ce69d050748>\";expires=600")

}
