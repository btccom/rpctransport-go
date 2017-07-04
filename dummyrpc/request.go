package dummyrpc

type DummyRequest struct {
	dummyServer *DummyServer
	body        []byte
}

func (dr *DummyRequest) Respond(respond []byte) error {
	dr.dummyServer.response = append(dr.dummyServer.response, respond)
	return nil
}

func (dr *DummyRequest) Body() []byte {
	return dr.body
}
