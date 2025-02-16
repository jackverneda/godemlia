package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackverneda/godemlia/internal/network"
	"github.com/jackverneda/godemlia/pb"
	"github.com/jackverneda/godemlia/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ip   = network.GetIpFromHost()
	port = 8080
)

func main() {
	if ip == "" {
		ip = "0.0.0.0"
	}
	fmt.Println("IP: ", ip)
	// peer := test.InitPeer(ip, port, 32140)
	peer := test.NewPeer(ip, port, 32140, true)

	grpcServer := grpc.NewServer()

	pb.RegisterNodeServer(grpcServer, &peer.Node)
	reflection.Register(grpcServer)

	go func() {
		<-time.After(time.Second * 5)
		fmt.Println("WORKING")
		payload := []byte("{\"username\":\"jackverneda\",\"email\":\"email\",\"password\":\"\"}")
		peer.Store("user", &payload)
	}()

	go peer.Node.Republish()

	grpcAddr := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

}
