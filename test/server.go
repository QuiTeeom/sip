package main

import (
	"net/http"
	"sip/pkg/sip-server"
	"sip/pkg/sip-server/sip-registrar"
)

func main() {
	sipServer := sip_server.SipServer{
		NetOptions: []sip_server.NetOption{
			{
				Type: "udp", Address: "0.0.0.0:45060",
			}, {
				Type: "ws", Address: "/ws-sip",
			},
		},
	}

	r := sip_registrar.Registrar{
		SipServer: &sipServer,
		Domain:    "quitee.com",
	}

	r.Start()

	http.Handle("/", http.FileServer(http.Dir("G:/MyWorkspaces/go2/sip/test/web")))
	http.ListenAndServe(":8080", nil)
}
