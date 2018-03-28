package amqprpc

import (
	_assert "github.com/stretchr/testify/require"
	"testing"
	"github.com/btccom/rpctransport-go/rpc"
)

func TestDefaultEnvVars(t *testing.T) {
	assert := _assert.New(t)
	defaults := DefaultAMQPEnvVars

	assert.Equal("host", varHost)
	assert.Equal("port", varPort)
	assert.Equal("vhost", varVHost)
	assert.Equal("user", varUser)
	assert.Equal("password", varPassword)

	assert.NoError(defaults.Check("host"))
	assert.NoError(defaults.Check("port"))
	assert.NoError(defaults.Check("vhost"))
	assert.NoError(defaults.Check("user"))
	assert.NoError(defaults.Check("password"))

	h, e := defaults.Var("host")
	assert.NoError(e)

	vh, e := defaults.Var("vhost")
	assert.NoError(e)

	p, e := defaults.Var("port")
	assert.NoError(e)

	u, e := defaults.Var("user")
	assert.NoError(e)

	pw, e := defaults.Var("password")
	assert.NoError(e)

	assert.Equal("TRANSPORT_AMQP_HOST", h)
	assert.Equal("TRANSPORT_AMQP_PORT", p)
	assert.Equal("TRANSPORT_AMQP_VHOST", vh)
	assert.Equal("TRANSPORT_AMQP_USER", u)
	assert.Equal("TRANSPORT_AMQP_PASSWORD", pw)
}

func TestConfigChecksEnvVarMap(t *testing.T) {
	assert := _assert.New(t)
	c := &AMQPConfig{}
	err := c.LoadConfigFromEnv(rpc.NewEnvVarMap(nil))
	assert.Error(err)
	assert.EqualError(err, "Missing host from env map")
}
