package tcprpc

import "net"

type ServerTcpRequest struct {
	conn net.Conn
	body []byte
}

func (r *ServerTcpRequest) Body() []byte {
	return r.body
}

func (r *ServerTcpRequest) Respond(response []byte) error {
	r.conn.Write(response)
	r.conn.Write([]byte{0x0d, 0x0a})
	r.conn.Close()

	return nil
}

