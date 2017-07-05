package rpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvVarNotFound(t *testing.T) {
	m := NewEnvVarMap(nil)
	assert.NotNil(t, m)

	_, err := m.Var("schrodinger")
	assert.Error(t, err)
}

func TestEnvCheckNotFound(t *testing.T) {
	m := NewEnvVarMap(nil)
	assert.NotNil(t, m)

	err := m.Check("check")
	assert.Error(t, err)
}

func TestEnvMap(t *testing.T) {
	fixtures := make([][2]string, 2)
	fixtures[0] = [2]string{"key", "value"}
	fixtures[1] = [2]string{"user", "ENV_KEY_FOR_USER"}

	nVars := len(fixtures)

	envmap := make(map[string]string)
	for i := 0; i < nVars; i++ {
		envmap[fixtures[i][0]] = fixtures[i][1]
	}

	descCreate := "test new EnvVarMap"
	t.Run(descCreate, func(t *testing.T) {
		m := NewEnvVarMap(envmap)
		assert.NotNil(t, m)
		assert.IsType(t, EnvVarMap{}, *m)
	})

	for i := 0; i < nVars; i++ {
		desc := fmt.Sprintf("Fetch EnvVar name for %s", fixtures[i][0])
		t.Run(desc, func(t *testing.T) {
			m := NewEnvVarMap(envmap)
			assert.NotNil(t, m)

			v, err := m.Var(fixtures[i][0])
			assert.NoError(t, err)
			assert.Equal(t, fixtures[i][1], v)
		})
	}

	allKeys := make([]string, nVars)
	for i := 0; i < nVars; i++ {
		allKeys[i] = fixtures[i][0]
	}

	t.Run("Check finds all keys", func(t *testing.T) {
		m := NewEnvVarMap(envmap)
		assert.NotNil(t, m)

		res := m.Check(allKeys)
		assert.Nil(t, res)
		assert.NoError(t, res)
	})
}
