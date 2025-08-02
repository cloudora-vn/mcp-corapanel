package utils

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
)


func GetRandomStr(e int) string {
	const charset = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	var result strings.Builder
	for i := 0; i < e; i++ {
		index := rand.Intn(len(charset))
		result.WriteByte(charset[index])
	}
	return result.String()
}

func GetPortFromAddr(addr string) (string, error) {
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	hostPort := parsedURL.Host
	if strings.Contains(hostPort, ":") {
		parts := strings.Split(hostPort, ":")
		return parts[len(parts)-1], nil
	}

	return "", fmt.Errorf("port not found")
}
