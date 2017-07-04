package dummyrpc

import "github.com/btccom/rpctransport-go/rpc"

func NewDummyServer() *DummyServer {
	return &DummyServer{
		Response: make([][]byte, 0),
		In:       make(chan []byte),
	}
}

type DummyServer struct {
	Response [][]byte
	In       chan []byte
}

func (dd *DummyServer) Consume() <-chan rpc.ServerRequest {
	requests := make(chan rpc.ServerRequest)

	go func() {
		for msg := range dd.In {
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
