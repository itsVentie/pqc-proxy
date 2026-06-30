package network

import (
	"errors"
	"fmt"
	"io"
	"net"

	"pqc-proxy/internal/crypto"
)

type Server struct {
	listenAddr string
	targetAddr string
	listener   net.Listener
	pqcKeys    *crypto.PQCKeyPair
}

func NewServer(listenAddr, targetAddr string, keys *crypto.PQCKeyPair) *Server {
	return &Server{
		listenAddr: listenAddr,
		targetAddr: targetAddr,
		pqcKeys:    keys,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("server failed to listen: %w", err)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			continue
		}

		InjectChaos(conn)
		go s.handleConnection(conn)
	}
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	clientInceptionBlob := make([]byte, 32+1184)
	if _, err := io.ReadFull(clientConn, clientInceptionBlob); err != nil {
		return
	}

	masterKey, responseBlob, err := crypto.ServerHandleInception(clientInceptionBlob)
	if err != nil {
		return
	}

	if _, err := clientConn.Write(responseBlob); err != nil {
		return
	}

	secureClientConn, err := crypto.NewSecureConn(clientConn, masterKey)
	if err != nil {
		return
	}
	crypto.SetServerRoles(secureClientConn)

	targetConn, err := net.Dial("tcp", s.targetAddr)
	if err != nil {
		return
	}
	defer targetConn.Close()

	errChan := make(chan error, 2)
	go func() { errChan <- proxyPipe(secureClientConn, targetConn) }()
	go func() { errChan <- proxyPipe(targetConn, secureClientConn) }()
	<-errChan
}

func proxyPipe(dst io.Writer, src io.Reader) error {
	bufPtr := GetBuffer()
	defer PutBuffer(bufPtr)
	buf := *bufPtr
	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
