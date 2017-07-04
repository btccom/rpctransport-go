package dummyrpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummyInterface(t *testing.T) {
	server := NewDummyServer()
	assert.NoError(t, server.Dial())
	assert.NoError(t, server.Close())
}
