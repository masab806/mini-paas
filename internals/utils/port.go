package utils

import (
	"fmt"
	"net"
)

func GetFreePort() (int, error) {

	listener, err := net.Listen(
		"tcp",
		"127.0.0.1:0",
	)

	if err != nil {
		return 0, fmt.Errorf(
			"failed to find free port: %w",
			err,
		)
	}

	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	return addr.Port, nil
}