package agentServer

import (
	"log"
	"net"
	"time"
)

// Client holds info about connection
type Client struct {
	conn     net.Conn
	host     string
	pingTime time.Time
}

// TCP server
type server struct {
	clients                  map[string]*Client
	address                  string // Address to open connection: localhost:9999
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

var gserver *server
var msgMap map[string]func(c *Client, a *AgentMsg)

const timeOutSec int64 = 20

// Start network server
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()
	//go s.CheckTimeout()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn: conn,
		}
		go client.OnMessage()
	}
}

// Creates new tcp server instance
func New(address string) {
	log.Println("Creating server with address", address)
	gserver = &server{
		address: address,
		clients: make(map[string]*Client),
	}
	msgMap = make(map[string]func(c *Client, a *AgentMsg))
	msgMap["token"] = TokenCheck
	msgMap["ping"] = Ping

	gserver.Listen()
}
