package amqprpc

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestAmqpDialNeedsValidServer(t *testing.T) {
	c := &AMQPConfig{}
	err := c.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	c.Hostname = "ghash.io"

	a := NewAmqpServer(c, "queue")

	err = a.Dial()
	assert.Error(t, err)
}

func TestAmqpDialCanBeCalledTwice(t *testing.T) {
	c := &AMQPConfig{}
	err := c.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	a := NewAmqpServer(c, "queue")
	err = a.Dial()
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, a.conn)
	assert.NotNil(t, a.ch)

	err = a.Dial()
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, a.conn)
	assert.NotNil(t, a.ch)

	err = a.Close()
	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestAmqpCloseWillErrorIfNotConnected(t *testing.T) {
	c := &AMQPConfig{}
	err := c.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	a := NewAmqpServer(c, "queue")
	err = a.Close()
	assert.Error(t, err)
	assert.Nil(t, a.conn)
	assert.Nil(t, a.ch)
}

func TestAmqpCanSetPrefetch(t *testing.T) {
	c := &AMQPConfig{}
	err := c.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	newPrefetch := 123

	a := NewAmqpServer(c, "queue")
	assert.Equal(t, 1, a.prefetch)
	a.SetPrefetch(newPrefetch)
	assert.Equal(t, newPrefetch, a.prefetch)
}

func TestAmqpDriverConnectDisconnect(t *testing.T) {
	c := &AMQPConfig{}
	err := c.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	a := NewAmqpServer(c, "queue")

	err = a.Dial()
	assert.NoError(t, err)
	assert.NotNil(t, a.conn)
	assert.NotNil(t, a.ch)

	err = a.Close()
	assert.NoError(t, err)
	assert.Nil(t, a.conn)
	assert.Nil(t, a.ch)
}

func TestAmqpConsume(t *testing.T) {
	cfg := &AMQPConfig{}
	err := cfg.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	queueName := "testQueue"

	c := &AmqpClient{}
	c.Init(cfg, queueName)
	assert.NoError(t, err)

	err = c.Dial()
	assert.NoError(t, err)
	defer c.Close()

	server := NewAmqpServer(cfg, queueName)
	err = server.Dial()
	assert.NoError(t, err)
	defer server.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	input := []byte(`input`)
	output := []byte(`output`)

	go func() {
		result, err := c.Request(input)
		wg.Done()
		assert.NoError(t, err)
		assert.Equal(t, output, result)
	}()

	go func() {
		select {
		case msg := <-server.Consume():
			assert.Equal(t, input, msg.Body())
			msg.Respond(output)
			wg.Done()
		}
	}()

	wg.Wait()

}
