//  mdcliapi class - Majordomo Protocol Client API
//  Implements the MDP/Worker spec at http://rfc.zeromq.org/spec:7
//
//  Author: iano <scaly.iano@gmail.com>
//  Based on C & Python example

package mdp

import (
	zmq "github.com/alecthomas/gozmq"
	"log"
	"time"
)

type Client interface {
	Close()
	Send([]byte, [][]byte) [][]byte
}

type MDClient struct {
	broker  string
	client  *zmq.Socket
	context *zmq.Context
	Retries int
	Timeout time.Duration
	verbose bool
}

func NewClient(broker string, verbose bool) *MDClient {
	context, _ := zmq.NewContext()
	self := &MDClient{
		broker:  broker,
		context: context,
		Retries: 3,
		Timeout: 2500 * time.Millisecond,
		verbose: verbose,
	}
	self.reconnect()
	return self
}

func (self *MDClient) reconnect() {
	if self.client != nil {
		self.client.Close()
	}

	self.client, _ = self.context.NewSocket(zmq.DEALER)
	self.client.SetLinger(0)
	self.client.Connect(self.broker)
	if self.verbose {
		log.Printf("I: connecting to broker at %s...\n", self.broker)
	}
}

func (self *MDClient) Close() {
	if self.client != nil {
		self.client.Close()
	}
	self.context.Close()
}

func (self *MDClient) Send(service []byte, request [][]byte) (reply [][]byte) {
	//  Prefix request with protocol frames
	//  Frame 1: "MDPCxy" (six bytes, MDP/Client x.y)
	//  Frame 2: Service name (printable string)
	frame := append([][]byte{[]byte(""), []byte(MDPC_CLIENT), service}, request...)
	if self.verbose {
		log.Printf("I: send request to '%s' service:", service)
		Dump(request)
	}

	for retries := self.Retries; retries > 0; {
		self.client.SendMultipart(frame, 0)
		items := zmq.PollItems{
			zmq.PollItem{Socket: self.client, Events: zmq.POLLIN},
		}

		_, err := zmq.Poll(items, self.Timeout)
		if err != nil {
			panic(err) //  Interrupted
		}

		if item := items[0]; item.REvents&zmq.POLLIN != 0 {
			msg, _ := self.client.RecvMultipart(0)
			if self.verbose {
				log.Println("I: received reply: ")
				Dump(msg)
			}

			//  We would handle malformed replies better in real code
			if len(msg) < 4 {
				panic("Error msg len")
			}

			header := msg[1]
			if string(header) != MDPC_CLIENT {
				panic("Error header")
			}

			replyService := msg[2]
			if string(service) != string(replyService) {
				panic("Error reply service")
			}

			reply = msg[3:]
			break
		} else if retries--; retries > 0 {
			if self.verbose {
				log.Println("W: no reply, reconnecting...")
			}
			self.reconnect()
		} else {
			if self.verbose {
				log.Println("W: permanent error, abandoning")
			}
			break
		}
	}
	return
}
