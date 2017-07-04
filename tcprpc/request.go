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
	r.conn.Write(msgEOF)
	r.conn.Close()

	return nil
}
