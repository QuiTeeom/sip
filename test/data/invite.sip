INVITE sip:bob@example.com SIP/2.0
Via: SIP/2.0/WS r9ik7446ccqm.invalid;branch=z9hG4bK5075365
Max-Forwards: 69
To: <sip:bob@example.com>
From: <sip:10000@quitee.com>;tag=kq74oia6ra
Call-ID: mrvmadori1cu6gm3rusa
CSeq: 8757 INVITE
Contact: <sip:as3sju0d@r9ik7446ccqm.invalid;transport=ws;ob>
Content-Type: application/sdp
Session-Expires: 90
Allow: INVITE,ACK,CANCEL,BYE,UPDATE,MESSAGE,OPTIONS,REFER,INFO
Supported: timer,ice,replaces,outbound
User-Agent: JsSIP 3.3.6
Content-Length: 267

v=0
o=mozilla...THIS_IS_SDPARTA-68.0.2 7271742767429479545 0 IN IP4 0.0.0.0
s=-
t=0 0
a=sendrecv
a=fingerprint:sha-256 16:AD:30:F4:4D:62:C8:22:AE:3B:10:E2:47:06:A2:19:4B:BB:73:2E:87:7D:AB:97:D4:D2:A6:84:18:F0:6D:8C
a=ice-options:trickle
a=msid-semantic:WMS *