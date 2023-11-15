package main

import (
	"log"
	"patterns/circuit-breaker/sms-provider/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
