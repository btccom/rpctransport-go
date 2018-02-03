package dummyrpc

type DummyRequest struct {
	dummyServer *DummyServer
	pending *pendingRequest
}

func (dr *DummyRequest) Respond(respond []byte) error {
	dr.pending.resultChan <- respond
	return nil
}

func (dr *DummyRequest) Body() []byte {
	return dr.pending.request
}
