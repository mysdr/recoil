package main

import (
	"bufio"
	"flag"
	"github.com/s-rah/go-ricochet"
	"log"
	"os"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "[Recoil]: ", log.Ltime|log.Lmicroseconds)

	target := flag.String("target", "", "the id of the ricochet client to use for testing")
	action := flag.String("action", "ping", "the action you want to take e.g. ping, connect")
	hostname := flag.String("hostname", "", "the hostname of a hidden service to use for authentication")
	privateKey := flag.String("key", "", "the private keyfile of a hidden service to use for authentication")
	messageFile := flag.String("messageFile", "", "a file containing a list of messages to send to the client")
	debug := flag.Bool("debug", false, "print the ricochet debug log to stdout")
	name := flag.String("name", "recoil", "a name to use when sending a contact request")
	message := flag.String("message", "I am the recoil testing tool", "a message to send during the contact request")
	flag.Parse()

	if *target == "" {
		logger.Fatalf("target must be specified")
	}

	if *hostname == "" {
		logger.Fatalf("hostname must be specified")
	}

	if *privateKey == "" {
		logger.Fatalf("key must be specified")
	}

	ricochet := new(goricochet.Ricochet)
	ricochet.Init(*privateKey, *debug)
	err := ricochet.Connect(*hostname, *target)
	if err != nil {
		logger.Printf("%s appears to be offline", *target)
	} else {
		logger.Printf("%s appears to be online", *target)
	}

	if *action == "contact-request" {
		ricochet := new(goricochet.Ricochet)
		ricochet.Init(*privateKey, *debug)
		ricochet.Connect(*hostname, *target)
		ricochet.SendContactRequest(*name, *message)
		logger.Printf("Sent contact request to [%s]", *target)
	}

	if *action == "send-messages" {
		file, err := os.Open(*messageFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		ricochet.OpenChannel("im.ricochet.chat", 5)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			message := scanner.Text()
			if len(message) > 0 && message[0] != '#' {
				logger.Printf("Sending message: %+q", message)
				ricochet.SendMessage(scanner.Text(), 5)
				time.Sleep(time.Second * 1)
			}

			if len(message) > 2 && message[0] == '#' {
				logger.Printf("Sending %s", message[2:])
			}
		}
	}

}
