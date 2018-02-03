package dummyrpc

import (
	_assert "github.com/stretchr/testify/require"
	"testing"
	"sync"
	"strconv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDummyInterface(t *testing.T) {
	server := NewDummyServer()
	_assert.NoError(t, server.Dial())
	_assert.NoError(t, server.Close())
}

func TestDummyClientCanPassBackErrors(t *testing.T) {
	server := NewDummyServer()
	_assert.NoError(t, server.Dial())
	_assert.NoError(t, server.Close())

	var wg sync.WaitGroup
	wg.Add(1)
	go func(server *DummyServer) {
		msg := <- server.Consume()
		dummyReq, ok := msg.(*DummyRequest)
		_assert.True(t, ok)
		_assert.IsType(t, DummyRequest{}, *dummyReq)
		dummyReq.pending.errorChan <- errors.New("oops-a-testing-error-occurred")
		wg.Done()
	}(server)

	client := NewDummyClient(server)
	_assert.NoError(t, client.Dial())

	_, err := client.Request([]byte{})
	assert.Error(t, err)

	server.Close()
	wg.Wait()
	client.Close()
}

func TestDummyServer(t *testing.T) {
	server := NewDummyServer()
	_assert.NoError(t, server.Dial())
	_assert.NoError(t, server.Close())

	var wg sync.WaitGroup
	go func(server *DummyServer) {
		for msg := range server.Consume() {
			bodyStr := string(msg.Body())
			intReq, err := strconv.Atoi(bodyStr)
			_assert.NoError(t, err)

			intRes := intReq * 2
			strIntRes := strconv.Itoa(intRes)
			msg.Respond([]byte(strIntRes))
		}
		wg.Done()
	}(server)

	client := NewDummyClient(server)
	_assert.NoError(t, client.Dial())

	fixtures := []struct{
		request string
		expecting int
	}{
		{
			request: "1",
			expecting: 2,
		},
		{
			request: "2",
			expecting: 4,
		},
		{
			request: "4",
			expecting: 8,
		},
		{
			request: "8",
			expecting: 16,
		},
	}

	var waitFixtures sync.WaitGroup
	for i := 0; i < len(fixtures); i++ {
		waitFixtures.Add(1)
		go func(fixture struct{
			request string
			expecting int}) {

			response, err := client.Request([]byte(fixture.request))
			_assert.NoError(t, err)
			_assert.NotNil(t, response)

			decoded := string(response)
			decodedInt, err := strconv.Atoi(decoded)
			_assert.NoError(t, err)

			_assert.Equal(t, fixture.expecting, decodedInt)
			waitFixtures.Done()
		}(fixtures[i])
	}

	waitFixtures.Wait()
	server.Close()
	wg.Wait()
	client.Close()
}
