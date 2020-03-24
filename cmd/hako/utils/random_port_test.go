package utils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
)

func Test_RandomFreePortShouldFindAFreePort(t *testing.T) {
	freePort, err := RandomFreePort()

	assert.NoError(t, err)
	assert.NoError(t, tryLister(freePort))
}

func tryLister(port int) (err error) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
	if err == nil {
		listener, err := net.ListenTCP("tcp", addr)
		if err == nil {
			_ = listener.Close()
		}
	}

	return err
}
