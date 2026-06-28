package crypto

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	MaxFrameSize = 32 * 1024
	HeaderSize   = 2
)

type SecureConn struct {
	net.Conn
	aead       cipher.AEAD
	readNonce  []byte
	writeNonce []byte
	readBuf    []byte
	leftover   []byte
	writeMu    sync.Mutex
	readMu     sync.Mutex
}

func NewSecureConn(baseConn net.Conn, masterKey []byte) (*SecureConn, error) {
	aead, err := chacha20poly1305.New(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create chacha20poly1305 instance: %w", err)
	}

	readNonce := make([]byte, aead.NonceSize())
	writeNonce := make([]byte, aead.NonceSize())

	return &SecureConn{
		Conn:       baseConn,
		aead:       aead,
		readNonce:  readNonce,
		writeNonce: writeNonce,
		readBuf:    make([]byte, HeaderSize+MaxFrameSize+aead.Overhead()),
	}, nil
}

func SetClientRoles(sc *SecureConn) {
	sc.writeNonce[0] = 1
}

func SetServerRoles(sc *SecureConn) {
	sc.readNonce[0] = 1
}

func (c *SecureConn) Write(b []byte) (n int, err error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	left := len(b)
	src := b

	for left > 0 {
		chunkSize := left
		if chunkSize > MaxFrameSize {
			chunkSize = MaxFrameSize
		}

		encryptedFrame, err := c.encryptFrame(src[:chunkSize])
		if err != nil {
			return n, err
		}

		_, err = c.Conn.Write(encryptedFrame)
		if err != nil {
			return n, err
		}

		n += chunkSize
		left -= chunkSize
		src = src[chunkSize:]
	}

	return n, nil
}

func (c *SecureConn) Read(b []byte) (n int, err error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()

	if len(c.leftover) > 0 {
		n = copy(b, c.leftover)
		c.leftover = c.leftover[n:]
		return n, nil
	}

	decryptedPayload, err := c.readAndDecryptFrame()
	if err != nil {
		return 0, err
	}

	n = copy(b, decryptedPayload)
	if n < len(decryptedPayload) {
		c.leftover = append(c.leftover, decryptedPayload[n:]...)
	}

	return n, nil
}

func (c *SecureConn) encryptFrame(payload []byte) ([]byte, error) {
	payloadLen := len(payload)
	frame := make([]byte, HeaderSize+payloadLen+c.aead.Overhead())

	binary.BigEndian.PutUint16(frame[0:HeaderSize], uint16(payloadLen+c.aead.Overhead()))

	c.aead.Seal(frame[HeaderSize:HeaderSize], c.writeNonce, payload, nil)
	incrementNonce(c.writeNonce)

	return frame, nil
}

func (c *SecureConn) readAndDecryptFrame() ([]byte, error) {
	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(c.Conn, header); err != nil {
		return nil, err
	}

	frameLen := binary.BigEndian.Uint16(header)
	if int(frameLen) > MaxFrameSize+c.aead.Overhead() {
		return nil, fmt.Errorf("incoming frame size exceeds maximum limit")
	}

	encryptedBuf := make([]byte, frameLen)
	if _, err := io.ReadFull(c.Conn, encryptedBuf); err != nil {
		return nil, err
	}

	decrypted, err := c.aead.Open(nil, c.readNonce, encryptedBuf, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt frame: %w", err)
	}
	incrementNonce(c.readNonce)

	return decrypted, nil
}

func incrementNonce(nonce []byte) {
	for i := 1; i < len(nonce); i++ {
		nonce[i] = nonce[i] + 1
		if nonce[i] != 0 {
			break
		}
	}
}
