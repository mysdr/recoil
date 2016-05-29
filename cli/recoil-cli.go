package main

import (
	"flag"
	"github.com/s-rah/recoil"
	"log"
)

func main() {

	target := flag.String("target", "", "the id of the ricochet client to use for testing")
	action := flag.String("action", "ping", "the action you want to take e.g. ping, contact-request")
	hostname := flag.String("hostname", "", "the hostname of a hidden service to use for authentication")
	privateKey := flag.String("key", "", "the private keyfile of a hidden service to use for authentication")
	messageFile := flag.String("messageFile", "", "a file containing a list of messages to send to the client")
	name := flag.String("name", "recoil", "a name to use when sending a contact request")
	message := flag.String("message", "I am the recoil testing tool", "a message to send during the contact request")
	flag.Parse()

	if *target == "" {
		log.Fatalf("target must be specified")
	}

	if *hostname == "" {
		log.Fatalf("hostname must be specified")
	}

	if *privateKey == "" {
		log.Fatalf("key must be specified")
	}

	recoil := new(recoil.Recoil)
	recoil.Ready = make(chan bool)

	if *action == "ping" {
		online := recoil.Ping(*privateKey, *hostname, *target)
		if online == true {
			log.Printf("%s is online", *target)
		} else {
			log.Printf("%s is offline", *target)
		}
	} else {
		go recoil.Authenticate(*privateKey, *hostname, *target)
		log.Printf("Running Recoil...")
		ready := <-recoil.Ready
		log.Printf("Received Authentication Result %v", ready)
		if ready == true {
			if *action == "contact-request" {
				recoil.SendContactRequest(*name, *message)
			}

			if *action == "spamchannel" {
				recoil.SpamChannel()
			}

			if *action == "sendmessage" {
				recoil.SendMessage(*messageFile)
			}
		}
	}
}
