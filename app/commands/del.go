package commands

import (
	"net"

	"github.com/askorupskyy/golang-redis-server/app/cache"
	"github.com/askorupskyy/golang-redis-server/app/commands/helpers"
)

func HandleDel(args []string, conn net.Conn) {
	key := helpers.GetKey(args)

	success := cache.Cache.Del(key)

	if !success {
		conn.Write([]byte("-0\r\n"))
		return
	}

	conn.Write([]byte("+OK\r\n"))
	return
}
