package agentServer

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"time"
)

// Client holds info about connection
type Client struct {
	conn     net.Conn
	Server   *server
	host     string
	tmpCid   int
	pingTime time.Time
}

// TCP server
type server struct {
	clients                  map[string]*Client
	tmpClients               map[int]*Client
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

var gserver *server
var tmpCid int

const timeOutSec int64 = 20

// Read client data from channel
func (c *Client) listen() {
	buffer := make([]byte, 1024)
	for {
		reader := bufio.NewReader(c.conn)
		len, err := reader.Read(buffer)
		if err != nil {
			log.Println("msg error:", err.Error())
			DisConnect(c)
			return
		}
		dataLength := binary.LittleEndian.Uint32(buffer)
		if dataLength <= 0 || dataLength > 1020 {
			continue
		}
		a := AgentMsg{}
		json.Unmarshal(buffer[4:dataLength+4], &a)
		log.Println("recv agent msg, msg: ", a, dataLength, len)
		if a.Cmd == "ping" {
			c.pingTime = time.Now()
		}
	}
}

// Send text message to client
func (c *Client) Send(message string) error {
	log.Println("send msg:", message)
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Start network server
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()
	go s.CheckTimeout()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

func (s *server) CheckTimeout() {
	var diffSec int64
	for {
		now := time.Now().Unix()
		for _, v := range s.clients {
			diffSec = now - v.pingTime.Unix()
			log.Println("host now:", now)
			log.Println("host pingtime:", v.pingTime.Unix())

			if diffSec > timeOutSec {
				DisConnect(v)
			}
		}
		for _, v := range s.tmpClients {
			diffSec = now - v.pingTime.Unix()
			log.Println("now:", now)
			log.Println("pingtime:", v.pingTime.Unix())
			if diffSec > timeOutSec {
				DisConnect(v)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

// Creates new tcp server instance
func New(address string) {
	log.Println("Creating server with address", address)
	gserver = &server{
		address:    address,
		clients:    make(map[string]*Client),
		tmpClients: make(map[int]*Client),
	}

	gserver.OnNewClient(NewClient)
	//gserver.OnNewMessage(OnMessage)
	//gserver.OnClientConnectionClosed(DisConnect)

	gserver.Listen()
}
