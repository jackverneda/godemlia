package network

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jackverneda/godemlia/internal/message"
)

type Broadcast struct {
	Port int
}

func (b *Broadcast) Recv(recv chan<- message.Message) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("0.0.0.0:%d", b.Port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 128)

	for {
		_, rAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err, rAddr)
		}

		recv <- message.Message{
			IP:     net.ParseIP(strings.Split(rAddr.String(), ":")[0]),
			Buffer: &buffer,
		}
	}
}

func (b *Broadcast) Send(msg *[]byte) {
	addr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: b.Port,
	}

	conn, err := net.DialUDP("udp4", nil, &addr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(*msg)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
