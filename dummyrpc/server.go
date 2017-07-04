package dummyrpc

import "github.com/btccom/rpctransport-go/rpc"

func NewDummyServer() *DummyServer {
	return &DummyServer{
		response: make([][]byte, 0),
		in:       make(chan []byte),
	}
}

type DummyServer struct {
	response [][]byte
	in       chan []byte
}

func (dd *DummyServer) Consume() <-chan rpc.ServerRequest {
	requests := make(chan rpc.ServerRequest)

	go func() {
		for msg := range dd.in {
			requests <- &DummyRequest{dd, msg}
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
