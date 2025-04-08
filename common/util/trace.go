package util

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"
)

func GenerateSpanID(addr string) string {
	strAddr := strings.Split(addr, ":")
	ip := strAddr[0]

	ipLong, _ := Ip2Long(ip)

	times := uint64(time.Now().UnixNano())

	// 使用 crypto/rand 生成更安全的随机数
	var randBytes [4]byte
	if _, err := rand.Read(randBytes[:]); err != nil {
		panic(err)
	}
	randInt := binary.BigEndian.Uint32(randBytes[:])

	spanId := ((times ^ uint64(ipLong)) << 32) | uint64(randInt)
	return strconv.FormatUint(spanId, 16)
}

func Ip2Long(ip string) (uint32, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(ipAddr.IP.To4()), nil
}
