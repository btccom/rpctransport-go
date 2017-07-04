package amqprpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultEnvVars(t *testing.T) {
	defaults := DefaultAMQPEnvVars
	assert.NoError(t, defaults.Check("host"))
	assert.NoError(t, defaults.Check("port"))
	assert.NoError(t, defaults.Check("user"))
	assert.NoError(t, defaults.Check("password"))

	h, e := defaults.Var("host")
	assert.NoError(t, e)

	p, e := defaults.Var("port")
	assert.NoError(t, e)

	u, e := defaults.Var("user")
	assert.NoError(t, e)

	pw, e := defaults.Var("password")
	assert.NoError(t, e)

	assert.Equal(t, "TRANSPORT_AMQP_HOST", h)
	assert.Equal(t, "TRANSPORT_AMQP_PORT", p)
	assert.Equal(t, "TRANSPORT_AMQP_USER", u)
	assert.Equal(t, "TRANSPORT_AMQP_PASSWORD", pw)
}
