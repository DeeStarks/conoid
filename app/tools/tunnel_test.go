package tools_test

import (
	"net"
	"testing"
	"time"

	"github.com/DeeStarks/conoid/app/tools"
)

func TestTunnel(t *testing.T) {
	tests := []struct {
		name, svr string
	}{
		{name: "abc", svr: ":30001"},
		{name: "cd", svr: ":30002"},
		{name: "agha", svr: ":30003"},
		{name: "abasdjjc", svr: ":30004"},
		{name: "a", svr: ":30005"},
	}
	openConns := make(chan net.Conn, len(tests)*2)

	// Create main server
	mSrv, err := net.Listen("tcp", ":30000")
	if err != nil {
		t.Error(err)
	}

	for _, tc := range tests {
		tunnel := tools.NewTunnel(tc.name, openConns)
		h, err := tunnel.AllocateHost()
		if err != nil {
			t.Error("Error allocating host:", err)
		}

		// Open tunnel
		_, err = net.Listen("tcp", tc.svr)
		if err != nil {
			t.Error(err)
		}
		remoteAddrCh := make(chan string, 1)
		h.OpenTunnel(mSrv.Addr().String(), remoteAddrCh)
	}

L:
	for {
		select {
		case conn := <-openConns:
			err := conn.Close()
			if err != nil {
				t.Error("Error closing connection", err)
			}
		case <-time.After(time.Second * 5):
			break L
		}
	}
}
