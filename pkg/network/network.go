package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func GetIpFromHost() string {
	cmd := exec.Command("hostname", "-i")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		//fmt.Println("Error running docker inspect:", err)
		return ""
	}
	ip := strings.TrimSpace(out.String())
	return ip
}

func Send(addr net.Addr, buffer *[]byte) {
	conn, err := net.Dial("tcp", addr.String())
	for err != nil {
		conn, err = net.Dial("tcp", addr.String())
	}
	defer conn.Close()
	// fmt.Println("Stablish connection")

	conn.Write(*buffer)
}

func Recv(recv chan *net.TCPConn, port int) {
	address := fmt.Sprintf("%s:%d", GetIpFromHost(), port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address)

	lis, err := net.ListenTCP("tcp", tcpAddr)
	for err != nil {
		lis, err = net.ListenTCP("tcp", tcpAddr)
	}
	defer lis.Close()

	tcpConn, _ := lis.AcceptTCP()

	select {
	case conn, closed := <-recv:
		if !closed {
			recv <- conn
		}
	case <-time.After(100 * time.Millisecond):
		recv <- tcpConn
	}
}
