package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"patterns/circuit-breaker/protobuf"
)

const grpcListen = ":4004"

type SmsServerResource struct {
	protobuf.UnsafeSmsNotifierServer
}

func NewSmsServerResource() *SmsServerResource {
	return &SmsServerResource{}
}

func (sr *SmsServerResource) Send(ctx context.Context, req *protobuf.SmsRequest) (*protobuf.SmsResponse, error) {
	val := rand.Float64()
	if val < 0.2 {
		log.Printf("got request, fail")
		return &protobuf.SmsResponse{Status: "error"}, fmt.Errorf("error")
	}

	log.Printf("got request,success")
	return &protobuf.SmsResponse{Status: fmt.Sprintf("successfully sent from %s to %s", req.From, req.To)}, nil
}

func Run() error {
	lis, err := net.Listen("tcp", grpcListen)
	if err != nil {
		log.Fatalf(err.Error())
	}

	srv := grpc.NewServer()

	resource := NewSmsServerResource()
	protobuf.RegisterSmsNotifierServer(srv, resource)

	log.Printf("serving grpc on %s", grpcListen)
	return srv.Serve(lis)
}
