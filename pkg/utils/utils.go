package utils

import (
	"crypto/sha1"
	"strconv"
)

func NewID(ip string, port int) ([]byte, error) {
	dumb := []byte(ip + ":" + strconv.FormatInt(int64(port), 10))
	hashValue := sha1.Sum(dumb)
	return []byte(hashValue[:]), nil
}
