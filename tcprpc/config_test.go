package tcprpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultEnvVars(t *testing.T) {
	defaults := DefaultTCPEnvVars

	assert.Equal(t, "host", varHost)
	assert.Equal(t, "port", varPort)
	assert.Equal(t, "tls", varTLS)

	assert.NoError(t, defaults.Check("host"))
	assert.NoError(t, defaults.Check("port"))
	assert.NoError(t, defaults.Check("tls"))

	h, e := defaults.Var("host")
	assert.NoError(t, e)

	p, e := defaults.Var("port")
	assert.NoError(t, e)

	tls, e := defaults.Var("tls")
	assert.NoError(t, e)

	assert.Equal(t, "TRANSPORT_TCP_HOST", h)
	assert.Equal(t, "TRANSPORT_TCP_PORT", p)
	assert.Equal(t, "TRANSPORT_TCP_TLS", tls)
}
