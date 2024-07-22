package main

import (
	"cust-mobile-eu/rmqsreceiver"
	//log "github.com/sirupsen/logrus"
)

func main() {

	//log.Println("Starting Receiver Service")
	//receiver := rmqsreceiver.Connect()
	
	receiver := rmqsreceiver.Connect("ProducedEU")
	receiver.ReceiveMessage()
}
