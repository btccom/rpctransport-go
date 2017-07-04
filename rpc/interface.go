package rpc

import (
	"fmt"
	"sync"
)

type Client interface {
	Request(req []byte) ([]byte, error)
	Dial() error
	Close() error
}

type Server interface {
	Consume() <-chan ServerRequest
	Dial() error
	Close() error
}

type ServerRequest interface {
	Respond(response []byte) error
	Body() []byte
}

type EnvVarMap struct {
	sync.RWMutex
	M map[string]string
}

func NewEnvVarMap(vars map[string]string) *EnvVarMap {
	m := make(map[string]string, len(vars))
	for k, v := range vars {
		m[k] = v
	}
	return &EnvVarMap{
		M: m,
	}
}
func (m *EnvVarMap) Var(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	if v, ok := m.M[key]; !ok {
		return "", fmt.Errorf("Unknown envvar name")
	} else {
		return v, nil
	}
}

func (m *EnvVarMap) Check(vars ...string) error {
	m.RLock()
	defer m.RUnlock()
	for i := 0; i < len(vars); i++ {
		_, ok := m.M[vars[i]]
		if !ok {
			return fmt.Errorf("Missing %s from env map", vars[i])
		}
	}

	return nil
}
