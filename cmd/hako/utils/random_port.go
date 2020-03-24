package utils

import (
	"net"
)

func RandomFreePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")

	if err == nil {
		listener, err := net.ListenTCP("tcp", addr)
		if err == nil {
			port = listener.Addr().(*net.TCPAddr).Port
			_ = listener.Close()
		}
	}

	return port, err
}
