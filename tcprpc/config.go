package tcprpc

import (
	"fmt"
	"github.com/btccom/rpctransport-go/rpc"
	"github.com/btccom/rpctransport-go/util"
	"strconv"
)

type TCPConfig struct {
	Host string
	Port int
	Type string
}

func (c *TCPConfig) Dsn() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var (
	msgEOF = []byte{0x0d, 0x0a}
)

const (
	varPort = "port"
	varHost = "host"
	varTLS  = "tls"
)

var DefaultTCPEnvVars = rpc.NewEnvVarMap(map[string]string{
	varHost: "TRANSPORT_TCP_HOST",
	varPort: "TRANSPORT_TCP_PORT",
	varTLS:  "TRANSPORT_TCP_TLS",
})

func (t *TCPConfig) LoadDefaultConfigFromEnv() error {
	return t.LoadConfigFromEnv(DefaultTCPEnvVars)
}

func (t *TCPConfig) LoadConfigFromEnv(varMap *rpc.EnvVarMap) error {
	err := varMap.Check([]string{varHost, varPort, varTLS}...)
	if err != nil {
		return err
	}

	hostVar, _ := varMap.Var(varHost)
	portVar, _ := varMap.Var(varPort)
	tlsVar, _ := varMap.Var(varTLS)

	port, err := strconv.Atoi(util.GetEnv(portVar, "6969"))
	if err != nil {
		return err
	}

	tls, err := strconv.ParseBool(util.GetEnv(tlsVar, "false"))
	if err != nil {
		return err
	}

	tcpType := "tcp"
	if tls {
		tcpType = "tls"
	}

	host := util.GetEnv(hostVar, "127.0.0.1")

	t.Port = port
	t.Host = host
	t.Type = tcpType

	return nil
}
