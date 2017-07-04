package amqprpc

import (
	"fmt"
	"github.com/btccom/rpctransport-go/rpc"
	"github.com/btccom/rpctransport-go/util"
	"os"
	"strconv"
)

var DefaultAMQPPort = 5672

type AMQPConfig struct {
	Username string
	Password string
	Hostname string
	Port     int
}

const (
	varPort = "port"
	varHost = "host"
	varUser = "user"
	varPassword = "password"
)

var DefaultAMQPEnvVars = rpc.NewEnvVarMap(map[string]string{
	"host":     "TRANSPORT_AMQP_HOST",
	"port":     "TRANSPORT_AMQP_PORT",
	"user":     "TRANSPORT_AMQP_USER",
	"password": "TRANSPORT_AMQP_PASSWORD",
})

func (c *AMQPConfig) LoadDefaultConfigFromEnv() error {
	return c.LoadConfigFromEnv(DefaultAMQPEnvVars)
}

func (c *AMQPConfig) LoadConfigFromEnv(varMap *rpc.EnvVarMap) error {
	err := varMap.Check([]string{varHost, varPort, varUser, varPassword}...)
	if err != nil {
		return err
	}

	hostVar, _ := varMap.Var(varHost)
	portVar, _ := varMap.Var(varPort)
	userVar, _ := varMap.Var(varUser)
	passVar, _ := varMap.Var(varPassword)

	portParam := os.Getenv(portVar)
	var port int
	if portParam == "" {
		port = DefaultAMQPPort
	} else {
		var err error
		port, err = strconv.Atoi(portParam)
		if err != nil {
			return err
		}
	}

	c.Hostname = util.GetEnv(hostVar, "localhost")
	c.Port = port
	c.Username = util.GetEnv(userVar, "guest")
	c.Password = util.GetEnv(passVar, "guest")

	return nil
}

func (c *AMQPConfig) Dsn() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Hostname, c.Port)
}
