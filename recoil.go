package recoil

import (
	"bufio"
	"github.com/s-rah/go-ricochet"
	"log"
	"os"
	"strings"
	"time"
)

type Recoil struct {
	goricochet.StandardRicochetService
	privateKey string
	hostname   string
	Ready      chan bool
}

func (r *Recoil) OnAuthenticationResult(channelID int32, serverHostname string, result bool) {
	if true {
		log.Printf("Successfully Authenticated to %s", serverHostname)
		r.Ready <- true
	} else {
		log.Fatal("Failed to authenticate to %s", serverHostname)
		r.Ready <- false
	}
}

func (r *Recoil) Ping(privateKey string, hostname string, target string) bool {
	r.privateKey = privateKey
	r.hostname = hostname
	r.Init(r.privateKey, r.hostname)
	err := r.Ricochet().Connect(r.hostname, target)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (r *Recoil) Authenticate(privateKey string, hostname string, target string) {
	if !r.Ping(privateKey, hostname, target) {
		r.Ready <- false
	}

	r.Init(r.privateKey, r.hostname)
	err := r.Ricochet().Connect(r.hostname, target)

	if err == nil {
		parts := strings.SplitAfterN(target, "|", 2)
		log.Printf("Recoil is connecting to %s", parts[1])

		r.OnConnect(parts[1])
		r.Ricochet().ListenAndWait(parts[1], r)
	} else {
		r.Ready <- false
	}
}

func (r *Recoil) SendContactRequest(name string, message string) {
	r.Ricochet().SendContactRequest(3, name, message)
}

func (r *Recoil) SpamChannel() {
	m := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		m[i] = 'a'
	}
	r.Ricochet().OpenChatChannel(5)
	for {
		r.Ricochet().SendMessage(5, string(m))
	}
}

func (r *Recoil) SendMessage(messageFile string) {
	r.Ricochet().OpenChatChannel(5)
	file, err := os.Open(messageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		message := scanner.Text()
		message = strings.Replace(message, "\\0", string(0x00), -1)
		if len(message) > 0 && message[0] != '#' {
			log.Printf("Sending message: %+q", message)
			r.Ricochet().SendMessage(5, message)
			time.Sleep(time.Second * 1)
		}

		if len(message) > 2 && message[0] == '#' {
			log.Printf("Sending %s", message[2:])
		}
	}
}
