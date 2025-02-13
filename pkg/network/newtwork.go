package network

import "net"

func Send(addr net.Addr, buffer *[]byte) {
	conn, err := net.Dial("tcp", addr.String())
	for err != nil {
		conn, err = net.Dial("tcp", addr.String())
	}
	defer conn.Close()
	// fmt.Println("Stablish connection")

	conn.Write(*buffer)
}
