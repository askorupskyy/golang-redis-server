package helpers

func GetKey(args []string) string {
	var key string

	if len(args) >= 5 {
		key = args[4]
	}

	return key
}
