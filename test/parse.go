package main

import (
	"fmt"
	"os"
	"sip/pkg/sip"
)

const DATA_PATH = "G:/MyWorkspaces/go2/sip/test"

func main() {
	f, _ := os.Open(DATA_PATH + "/data/invite.sip")
	defer f.Close()

	sipMessage, err := sip.Marshal(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(sipMessage.String())

	f2, _ := os.Open(DATA_PATH + "/data/200ok.sip")
	sipMessage, err = sip.Marshal(f2)
	if err != nil {
		panic(err)
	}
	fmt.Println(sipMessage.(sip.ResponseMessage).String())
}
