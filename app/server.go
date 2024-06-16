package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/askorupskyy/golang-redis-server/app/cache"
)

func handleConnection(conn net.Conn) {
	fmt.Printf("Client connected: %s\n", conn.LocalAddr().String())

	// each connection has a buffer
	buf := make([]byte, 512)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client ended stream: %s\n", conn.LocalAddr().String())
			}
			break
		}

		command := string(buf[:n])
		args := strings.Split(command, "\r\n")

		log.Printf("received args >>> %s\n", args)

		if args[2] == "PING" {
			cache.Cache.Set("ping", "pong", cache.SetCacheValueArgs{Expiry: 10000})
			n, err = conn.Write([]byte("+PONG\r\n"))
		} else if args[2] == "SET" {
			expiry, _ := strconv.Atoi(args[8])
			cache.Cache.Set(args[4], args[6], cache.SetCacheValueArgs{Expiry: int64(expiry)})
			n, err = conn.Write([]byte("+OK\r\n"))

		} else if args[2] == "GET" {
			val, exists := cache.Cache.Get(args[4])
			if !exists {
				n, err = conn.Write([]byte("-0\r\n"))
			}
			n, err = conn.Write([]byte(fmt.Sprintf("+%s\r\n", val)))

		} else {
			conn.Write([]byte("+Invalid command\r\n"))
		}

		if err != nil {
			fmt.Printf("conn.Write() failed: %s\n", err)
			break
		}
	}

	conn.Close()
	fmt.Printf("Client disconnected: %s\n", conn.LocalAddr().String())
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
