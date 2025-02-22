package commands

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/askorupskyy/golang-redis-server/app/cache"
)

// TODO: need to figure out how to extend this
type SetArgs struct {
	Key    string
	Value  any
	Expiry int64
}

func validateSetArgs(args []string) (SetArgs, error) {
	var value any = nil
	var key string = ""

	if len(args) >= 5 {
		key = args[4]
	}

	if key == "" {
		return SetArgs{}, errors.New("Key not specified, returning error")
	}

	if len(args) >= 7 {
		value = args[6]
	}

	if value == nil {
		return SetArgs{}, errors.New("Value not specified, returning error")
	}

	setArgs := SetArgs{
		Key:   key,
		Value: value,
	}

	for _, arg := range args {
		// check if we pass the `exp` arg
		if strings.Index(arg, "EXP") > -1 {
			post := strings.SplitN(arg, "=", 2)
			if len(post) > 1 {
				expiry, err := strconv.Atoi(post[1])

				if err != nil {
					return SetArgs{}, errors.New("Expiry value is an invalid integer, returning error")
				}

				setArgs.Expiry = int64(expiry)
			}
		}
	}

	return setArgs, nil
}

func HandleSet(args []string, conn net.Conn) {
	setArgs, err := validateSetArgs(args)

	if err != nil {
		log.Printf("conn.Write() in `SET` failed: %s\n", conn.LocalAddr())
		conn.Close()
	}

	cache.Cache.Set(setArgs.Key, setArgs.Value, cache.SetCacheValueArgs{Expiry: setArgs.Expiry})
	_, err = conn.Write([]byte("+OK\r\n"))

	if err != nil {
		log.Printf("conn.Write() in `SET` failed: %s\n", conn.LocalAddr())
		conn.Close()
	}
}
