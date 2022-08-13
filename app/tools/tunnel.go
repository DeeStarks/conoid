package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/DeeStarks/conoid/utils"
)

type (
	ITunnel interface {
		AllocateHost() (ITunnelHost, error)
	}

	tunnel struct {
		name         string
		tunnelServer string
		openConns    chan<- net.Conn
	}

	ITunnelHost interface {
		OpenTunnel(string, chan<- string)
		SubDomain() string
		FullURL() string
		MaxConnectionCount() int
		PortNumber() int
	}

	allocatedHost struct {
		Id           string `json:"id"`
		Port         int    `json:"port"`
		MaxConnCount int    `json:"max_conn_count"`
		Url          string `json:"url"`
		openConns    chan<- net.Conn
	}
)

// Create a new tunnel
func NewTunnel(name string, connCh chan<- net.Conn) ITunnel {
	return &tunnel{
		name:         name,
		tunnelServer: "http://localtunnel.me/",
		openConns:    connCh,
	}
}

// Creates a remote host for the service to be tunnelled
func (t *tunnel) AllocateHost() (ITunnelHost, error) {
	var host allocatedHost
	host.openConns = t.openConns

	// Names are required to be at least 4 in length
	subdomain := t.name
	if len(t.name) < 4 {
		suffix := []byte("1234")
		subdomain += string(suffix[:4-len(t.name)])
	}

	res, err := http.Get(t.tunnelServer + subdomain)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &host); err != nil {
		return nil, err
	}
	return &host, nil
}

func (h *allocatedHost) OpenTunnel(conoidServer string, connectedAddressCh chan<- string) {
	// Connect to server
	localConn, err := net.Dial("tcp", conoidServer)
	if err != nil {
		utils.Log("Error occured while tunneling:", err)
		return
	}
	// Add to open connections
	h.openConns <- localConn

	// Connect to remote host
	// Parse url
	pUrl, err := url.Parse(h.Url)
	if err != nil {
		utils.Log(err)
		return
	}
	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", pUrl.Host, h.Port))
	if err != nil {
		utils.Log("Error occured while", err)
		return
	}
	// Add to open connections
	h.openConns <- remoteConn

	// Send the address connected to localtunnel, to allow
	// conoid server know the local connection it is to connect
	connectedAddressCh <- localConn.LocalAddr().String()

	// This will check if a connection is closed
	// to stop goroutine from accepting more requests
	var isClosedConn = func(err error) bool {
		return strings.Contains(err.Error(), ": use of closed network connection")
	}

	// Establish a point-to-point connection between remote server and the local server
	go func() {
		for {
			_, err = io.Copy(localConn, remoteConn)
			if err != nil {
				if isClosedConn(err) {
					return
				}
				utils.Log(err)
			}
		}
	}()

	go func() {
		for {
			_, err = io.Copy(remoteConn, localConn)
			if err != nil {
				if isClosedConn(err) {
					return
				}
				utils.Log(err)
			}
		}
	}()
}

func (h *allocatedHost) SubDomain() string {
	return h.Id
}

func (h *allocatedHost) FullURL() string {
	return h.Url
}

func (h *allocatedHost) MaxConnectionCount() int {
	return h.MaxConnCount
}

func (h *allocatedHost) PortNumber() int {
	return h.Port
}
