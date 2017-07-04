package tcprpc

import "net"

type TcpRequest struct {
	conn net.Conn
	body []byte
}

func (r *TcpRequest) Body() []byte {
	return r.body
}

func (r *TcpRequest) Respond(response []byte) error {
	r.conn.Write(response)
	r.conn.Write([]byte{0x0d, 0x0a})
	r.conn.Close()

	return nil
}
