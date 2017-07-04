package dummyrpc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDummyInterface(t *testing.T) {
	server := NewDummyServer()
	assert.NoError(t, server.Dial())
	assert.NoError(t, server.Close())
}
