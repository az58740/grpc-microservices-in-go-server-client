package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/az58740/grpc-microservices-proto/golang/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	serverAddress = flag.String("add", "localhost:8080", "The server address in the format of host:port")
)

func CreateOrder(client pb.OrderClient) {
	log.Println("Creating order...")
	orderResponse, errCreate := client.Create(context.Background(), &pb.CreateOrderRequest{
		UserId: 2,
		OrderItems: []*pb.OrderItem{
			{
				ProductCode: "1",
				UnitPrice:   10,
				Quantity:    2,
			},
			{
				ProductCode: "2",
				UnitPrice:   20,
				Quantity:    2,
			},
		},
	})
	if errCreate != nil {
		log.Fatalf("Failed to create order Error:%v", errCreate)

	} else {
		log.Printf("The order created successfuly with ID:%d", orderResponse.GetOrderId())
	}
}
func getTSLCredentisls() (credentials.TransportCredentials, error) {
	clientCert, clientCertErr := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if clientCertErr != nil {
		return nil, fmt.Errorf("could not load client key pairs:%s", clientCertErr)
	}
	certPool := x509.NewCertPool()
	caCert, caCertErr := os.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not load CA cert:%s", caCertErr)
	}
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to aapend CA cert")
	}
	return credentials.NewTLS(&tls.Config{
		ServerName:   "*.microservices.dev",
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}), nil

}
func main() {
	flag.Parse()
	tlsCredentials, tlsCredentialsErr := getTSLCredentisls()
	if tlsCredentialsErr != nil {
		log.Fatal("cannot load client TLS credentials: ", tlsCredentialsErr)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	conn, err := grpc.NewClient(*serverAddress, opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}
	defer conn.Close()
	client := pb.NewOrderClient(conn)
	CreateOrder(client)
}
