package commands

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/askorupskyy/golang-redis-server/app/cache"
	"github.com/askorupskyy/golang-redis-server/app/commands/helpers"
)

func HandleTTL(args []string, conn net.Conn) {
	key := helpers.GetKey(args)

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
		return
	}

	expiryTimestamp := obj.Expiry.Local().UnixMilli()
	current := time.Now().UnixMilli()

	conn.Write([]byte(fmt.Sprintf("+%d\r\n", expiryTimestamp-current)))
}
