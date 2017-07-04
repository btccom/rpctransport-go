package tcprpc

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestTcpConsume(t *testing.T) {
	cfg := &TCPConfig{}
	err := cfg.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	server, err := NewTCPServer(cfg)
	assert.NoError(t, err)

	err = server.Dial()
	assert.NoError(t, err)
	defer server.Close()

	c, err := NewTCPClient(cfg)
	assert.NoError(t, err)

	err = c.Dial()
	assert.NoError(t, err)
	defer c.Close()

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
