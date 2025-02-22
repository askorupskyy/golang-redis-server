package commands

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/askorupskyy/golang-redis-server/app/cache"
)

func HandleTTL(args []string, conn net.Conn) {
	var key string

	if len(args) >= 5 {
		key = args[4]
	}

	if key == "" {
		log.Printf("Key `TTL` not provided: %s\n", conn.LocalAddr())
		conn.Close()
	}

	obj, exists := cache.Cache.Get(key)

	if !exists {
		_, err := conn.Write([]byte("-0\r\n"))
		if err != nil {
			log.Printf("conn.Write() in `SET` failed: %s\n", conn.LocalAddr())
			conn.Close()
		}
	}

	expiryTimestamp := obj.Expiry.Local().UnixMilli()
	current := time.Now().UnixMilli()

	conn.Write([]byte(fmt.Sprintf("+%d\r\n", expiryTimestamp-current)))
}
