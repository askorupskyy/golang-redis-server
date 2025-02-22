package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/askorupskyy/golang-redis-server/app/cache"
	"github.com/askorupskyy/golang-redis-server/app/commands"
)

func handleConnection(conn net.Conn) {
	log.Printf("Client connected: %s\n", conn.LocalAddr().String())

	// each connection has a buffer
	buf := make([]byte, 512)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				log.Printf("Client ended stream: %s\n", conn.LocalAddr().String())
			}
			break
		}

		command := string(buf[:n])
		args := strings.Split(command, "\r\n")

		log.Printf("received args >>> %s\n", args)

		switch args[2] {
		case "PING":
			n, err = conn.Write([]byte("+PONG\r\n"))
		case "SET":
			commands.HandleSet(args, conn)
		case "TTL":
			commands.HandleTTL(args, conn)
		case "GET":
			{
				val, exists := cache.Cache.Get(args[4])
				if !exists {
					n, err = conn.Write([]byte("-0\r\n"))
				} else {
					n, err = conn.Write([]byte(fmt.Sprintf("+%s\r\n", val)))
				}
			}
		default:
			conn.Write([]byte("+Invalid command\r\n"))
		}
	}

	conn.Close()
	log.Printf("Client disconnected: %s\n", conn.LocalAddr().String())
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		log.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
