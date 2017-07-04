package amqprpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAmqpClient(t *testing.T) {
	cfg := &AMQPConfig{}
	err := cfg.LoadDefaultConfigFromEnv()
	assert.NoError(t, err)

	queue := "str123"
	client, err := NewAmqpClient(cfg, queue)
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, AmqpClient{}, *client)
	assert.Equal(t, queue, client.workQueue)
}
