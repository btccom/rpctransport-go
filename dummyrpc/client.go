package dummyrpc

import (
	"sync"
	"github.com/btccom/rpctransport-go/rpc"
)

func NewDummyClient(server *DummyServer) rpc.Client {
	return &DummyClient{
		server: server,
	}
}

type pendingRequest struct {
	request []byte
	resultChan chan []byte
	errorChan  chan error
}

func newPendingRequest(request []byte) *pendingRequest {
	return &pendingRequest{
		request: request,
		resultChan: make(chan []byte),
		errorChan:  make(chan error),
	}
}

type DummyClient struct {
	server *DummyServer
	sendLock sync.RWMutex
}

func (ad *DummyClient) Close() error {
	return nil
}

func (ad *DummyClient) Dial() error {
	return nil
}

func (ad *DummyClient) RequestAsync(request []byte) (<-chan []byte, <-chan error) {
	pending := newPendingRequest(request)
	ad.server.In <- pending
	return pending.resultChan, pending.errorChan
}

func (ad *DummyClient) Request(request []byte) ([]byte, error) {
	resultChan, errorChan := ad.RequestAsync(request)

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

