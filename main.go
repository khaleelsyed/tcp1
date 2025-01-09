package main

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
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
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func main() {
	server := NewServer(":3000")
	server.Start()
}
