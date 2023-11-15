package main

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"patterns/circuit-breaker/protobuf"
	"time"
)

var (
	serverAddress = "localhost:4004"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 2 && failureRatio >= 0.3
	}
	st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		fmt.Printf("Circuit breaker %v state change from %v to %v\n", name, from, to)
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := protobuf.NewSmsNotifierClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {

		status, execErr := cb.Execute(func() (interface{}, error) {

			resp, sendErr := client.Send(ctx, &protobuf.SmsRequest{
				From: "87073929044",
				To:   "87073929045",
				Text: "some text",
			})
			if sendErr != nil {
				return nil, sendErr
			}

			return resp.Status, nil
		})

		fmt.Println(cb.Counts())

		if execErr != nil {
			log.Printf("[ERROR] send sms: %s\n", execErr)
			return
		}

		fmt.Println(status)
	}
}
