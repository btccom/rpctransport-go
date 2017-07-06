package tcprpc

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

func NewTCPClient(cfg *TCPConfig) (*TCPClient, error) {
	return &TCPClient{
		cfg:          cfg,
		sem:          make(chan struct{}, 1),
		dueResponses: make(chan *pendingRequest),
	}, nil
}

type pendingRequest struct {
	res  []byte
	err  error
	done chan *pendingRequest
}

type TCPClient struct {
	sync.RWMutex
	cfg          *TCPConfig
	conn         net.Conn
	dueResponses chan *pendingRequest
	sem          chan struct{}
}

func (c *TCPClient) Dial() error {
	conn, err := net.Dial("tcp", c.cfg.Dsn())
	if err != nil {
		return err
	}

	c.conn = conn

	go func() {
		for {
			pending := <-c.dueResponses
			response, err := c.readResponse()
			if err != nil {
				pending.err = err
			} else {
				pending.res = response[:len(response)-2]
			}

			pending.done <- pending
			close(pending.done)
		}
	}()

	return nil
}

func (c *TCPClient) Close() error {
	if c.conn == nil {
		return fmt.Errorf("TCP Connection not open")
	}

	c.conn.Close()
	c.conn = nil

	return nil
}

func (c *TCPClient) readResponse() ([]byte, error) {
	data := make([]byte, 0)
	n := 0
	for {
		buf := make([]byte, 1024)

		// Read the incoming connection into the buffer.
		reqLen, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}

		data = append(data, buf[:reqLen]...)
		n += reqLen

		if bytes.Equal(data[n-2:], msgEOF) {
			break
		}
	}

	return data, nil
}

func (c *TCPClient) RequestAsync(body []byte) (<- chan []byte, <- chan error) {

	c.Lock()
	c.conn.Write(body)
	c.conn.Write([]byte{0x0d, 0x0a})
	request := &pendingRequest{
		done: make(chan *pendingRequest),
	}
	resultChan := make(chan []byte)
	errorChan := make(chan []byte)
	c.dueResponses <- request
	go func() {
		complete := <-request.done
		if complete.err != nil {
			errorChan <- complete.err
			close(errorChan)
			close(resultChan)
		} else {
			resultChan <- complete.res
			close(errorChan)
			close(resultChan)
		}
	}()

	c.Unlock()

	return resultChan, errorChan

}

func (c *TCPClient) Request(body []byte) ([]byte, error) {
	resultChan, errorChan := c.RequestAsync(body)
	select {
	case result := <- resultChan:
		return result, nil
	case err := <- errorChan:
		return nil, err
	}
}
