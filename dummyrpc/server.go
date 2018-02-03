package dummyrpc

import "github.com/btccom/rpctransport-go/rpc"

func NewDummyServer() *DummyServer {
	return &DummyServer{
		In:       make(chan *pendingRequest),
	}
}

type DummyServer struct {
	In       chan *pendingRequest
}

func (dd *DummyServer) Consume() <-chan rpc.ServerRequest {
	requests := make(chan rpc.ServerRequest)

	go func() {
		for request := range dd.In {
			requests <- &DummyRequest{dd, request}
		}
	}()

	return requests
}

func (dd *DummyServer) Dial() error {
	return nil
}

func (dd *DummyServer) Close() error {
	return nil
}
