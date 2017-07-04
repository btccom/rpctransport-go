package tcprpc

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"github.com/btccom/rpctransport-go/rpc"
)

type ServerDriverTcp struct {
	cfg      *TCPConfig
	listener net.Listener
}

func (td *ServerDriverTcp) Init(cfg *TCPConfig) error {
	td.cfg = cfg
	return nil
}

func (td *ServerDriverTcp) Close() error {
	if td.listener == nil {
		return fmt.Errorf("TCP Server not connected")
	}

	td.listener.Close()
	td.listener = nil

	return nil
}

func (td *ServerDriverTcp) Dial() error {
	if td.listener != nil {
		return nil
	}

	l, err := net.Listen(td.cfg.Type, td.cfg.Host+":"+strconv.Itoa(td.cfg.Port))
	if err != nil {
		return err
	}

	td.listener = l
	return nil
}

func (td *ServerDriverTcp) Consume() <-chan rpc.ServerRequest {
	requests := make(chan rpc.ServerRequest)

	ending := []byte{0x0d, 0x0a}

	go func(listener net.Listener) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go func() {
				data := make([]byte, 0)
				n := 0
				for {
					// Make a buffer to hold incoming data.
					buf := make([]byte, 1024)

					// Read the incoming connection into the buffer.
					reqLen, err := conn.Read(buf)
					if err != nil {
						conn.Close()
						return
					}

					data = append(data, buf[:reqLen]...)
					n += reqLen

					if bytes.Equal(data[n-2:], ending) {
						break
					}
				}

				requests <- &ServerTcpRequest{conn, data[:n-2]}
			}()
		}
	}(td.listener)

	return requests
}
