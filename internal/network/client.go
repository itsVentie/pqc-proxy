package network

import (
	"fmt"
	"io"
	"net"

	"pqc-proxy/internal/crypto"
)

type Client struct {
	listenAddr string
	serverAddr string
	listener   net.Listener
}

func NewClient(listenAddr, serverAddr string) *Client {
	return &Client{
		listenAddr: listenAddr,
		serverAddr: serverAddr,
	}
}

func (c *Client) Start() error {
	var err error
	c.listener, err = net.Listen("tcp", c.listenAddr)
	if err != nil {
		return fmt.Errorf("client failed to listen: %w", err)
	}

	for {
		localConn, err := c.listener.Accept()
		if err != nil {
			return nil
		}
		go c.handleConnection(localConn)
	}
}

func (c *Client) Stop() {
	if c.listener != nil {
		c.listener.Close()
	}
}

func (c *Client) handleConnection(localConn net.Conn) {
	defer localConn.Close()

	serverConn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return
	}
	defer serverConn.Close()

	ecdhPriv, mlkemPriv, clientBlob, err := crypto.GenerateClientInception()
	if err != nil {
		return
	}

	if _, err := serverConn.Write(clientBlob); err != nil {
		return
	}

	serverResponseBlob := make([]byte, 32+1088)
	if _, err := io.ReadFull(serverConn, serverResponseBlob); err != nil {
		return
	}

	masterKey, err := crypto.ClientHandleResponse(ecdhPriv, mlkemPriv, serverResponseBlob)
	if err != nil {
		return
	}

	secureServerConn, err := crypto.NewSecureConn(serverConn, masterKey)
	if err != nil {
		return
	}
	crypto.SetClientRoles(secureServerConn)

	errChan := make(chan error, 2)
	go func() { errChan <- proxyPipe(localConn, secureServerConn) }()
	go func() { errChan <- proxyPipe(secureServerConn, localConn) }()
	<-errChan
}
