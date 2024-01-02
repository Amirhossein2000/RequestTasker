package test

import "net"

func GetAvailablePort() (int, error) {
	// Binding to port 0 will let the system choose a random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	// Get the port from the listener's address
	address := listener.Addr().(*net.TCPAddr)
	return address.Port, nil
}