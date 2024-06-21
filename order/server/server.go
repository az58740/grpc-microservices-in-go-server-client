package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/az58740/grpc-microservices-proto/golang/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.UnimplementedOrderServer
}

func getTSLCredentisls() (credentials.TransportCredentials, error) {
	serverCert, serverCertErr := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if serverCertErr != nil {
		return nil, fmt.Errorf("could not load server key pairs: %s", serverCertErr)
	}
	certPool := x509.NewCertPool()
	caCert, caCertErr := os.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not read CA cert: %s", caCertErr)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append the CA certs")
	}
	return credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequestClientCert,
			Certificates: []tls.Certificate{serverCert},
			ClientCAs:    certPool,
		}), nil
}

func (s *Server) Create(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return &pb.CreateOrderResponse{
		OrderId: 1234,
	}, nil
}
func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("Failed to listen on port %d . Error %v", 8080, err)
	}
	tlsCredentials, tlsCredentialsErr := getTSLCredentisls()
	if tlsCredentialsErr != nil {
		log.Fatal("cannot load server TLS credentials: ", tlsCredentialsErr)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(tlsCredentials))
	grpServer := grpc.NewServer(opts...)
	pb.RegisterOrderServer(grpServer, &Server{})
	log.Println("gRPC Server starting ...")
	if grpServer.Serve(listen); err != nil {
		log.Fatalf("Failed to start gprc server. Error:%v", err)
	}

}
