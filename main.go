package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	msgCh      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		msgCh:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()
	s.listener = listener

	go s.acceptLoop()

	<-s.quitCh
	close(s.msgCh)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
			// Continue used as if return was used, it will be the end of the loop
		}

		fmt.Println("new connection recieved:", conn.RemoteAddr())

		go s.readLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		s.msgCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgCh {
			fmt.Printf("Message from connection (%s): %s", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
