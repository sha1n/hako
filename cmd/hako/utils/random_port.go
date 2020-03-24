package utils

import (
	"net"
)

// RandomFreePort attempts to find and return a free TCP port number.
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
