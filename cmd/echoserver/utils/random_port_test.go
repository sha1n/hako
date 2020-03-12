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

func tryLister(port int) error {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
	if err == nil {
		listener, err := net.ListenTCP("tcp", addr)
		if err == nil {
			err = listener.Close()
		}
	}

	return err
}
