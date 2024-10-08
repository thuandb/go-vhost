package vhost

import (
	"crypto/tls"
	"net"
	"testing"
)

func TestSNI(t *testing.T) {
	var testHostname string = "foo.example.com"

	l, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		conf := &tls.Config{ServerName: testHostname}
		conn, err := tls.Dial("tcp", "127.0.0.1:12345", conf)
		if err != nil {
			panic(err)
		}
		conn.Close()
	}()

	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	c, err := TLS(conn)
	if err != nil {
		panic(err)
	}

	if c.Host() != testHostname {
		t.Errorf("Connection Host() is %s, expected %s", c.Host(), testHostname)
	}
}

func TestAlpnProtocols(t *testing.T) {
	rawBytes := []byte{
		0x01,
		0x00, 0x00, 0xc4, // hello msg length
		0x03, 0x03, // protocol version
		0xec, 0x12, 0xdd, 0x17, 0x64, 0xa4, 0x39, 0xfd, 0x7e, 0x8c, 0x85, 0x46, 0xb8, 0x4d, 0x1e, 0xa0, 0x6e, 0xb3, 0xd7, 0xa0, 0x51, 0xf0, 0x3c, 0xb8, 0x17, 0x47, 0x0d, 0x4c, 0x54, 0xc5, 0xdf, 0x72, // Random value
		0x00,       // session id length
		0x00, 0x1c, // cipher suite length
		0xea, 0xea, 0xc0, 0x2b, 0xc0, 0x2f, 0xc0, 0x2c, 0xc0, 0x30, 0xcc, 0xa9, 0xcc, 0xa8, 0xc0, 0x13, 0xc0, 0x14, 0x00, 0x9c, 0x00, 0x9d, 0x00, 0x2f, 0x00, 0x35, 0x00, 0x0a, // cipher suites
		0x01,
		0x00,
		0x00, 0x7f, // length of extensions
		0xda, 0xda, 0x00, 0x00,
		0xff, 0x01, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x16, 0x00, 0x14, 0x00, 0x00, 0x11, 0x77, 0x77, 0x77, 0x2E, 0x6D, 0x79, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x2E, 0x63, 0x6F, 0x6D, // ServerName: www.myexample.com
		0x00, 0x17, 0x00, 0x00,
		0x00, 0x23, 0x00, 0x00,
		0x00, 0x0d, 0x00, 0x14, 0x00, 0x12, 0x04, 0x03, 0x08, 0x04, 0x04, 0x01, 0x05, 0x03, 0x08, 0x05, 0x05, 0x01, 0x08, 0x06, 0x06, 0x01, 0x02, 0x01,
		0x00, 0x05, 0x00, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x12, 0x00, 0x00,
		0x00, 0x10, 0x00, 0x0e, 0x00, 0x0c, 0x02, 0x68, 0x32, 0x08, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x31, 0x2e, 0x31, // AlpnProtocols: [h2 http/1.1]
		0x75, 0x50, 0x00, 0x00,
		0x00, 0x0b, 0x00, 0x02, 0x01, 0x00,
		0x00, 0x0a, 0x00, 0x0a, 0x00, 0x08, 0x1a, 0x1a, 0x00, 0x1d, 0x00, 0x17, 0x00, 0x18,
		0x1a, 0x1a, 0x00, 0x01, 0x00,
	}

	helloMsg := &ClientHelloMsg{}

	if ok := helloMsg.unmarshal(rawBytes); !ok {
		t.Errorf("Failed to parse client hello mesasge!")
	}

	expectedAlpnProtocols := []string{"h2", "http/1.1"}

	if len(helloMsg.AlpnProtocols) != 2 {
		t.Errorf("Failed to parse AlpnProtocols. Expected %d, received %d!", len(expectedAlpnProtocols), len(helloMsg.AlpnProtocols))
	}

	for i := range helloMsg.AlpnProtocols {
		if helloMsg.AlpnProtocols[i] != expectedAlpnProtocols[i] {
			t.Errorf("Alpn protocol mismatched. Expected %s, received %s!", expectedAlpnProtocols[i], helloMsg.AlpnProtocols[i])
		}
	}
}
