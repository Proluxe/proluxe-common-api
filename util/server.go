package util

import (
	"log"
	"net"
	"strconv"

	u "github.com/scottraio/go-utils"
)

func Port() int {
	port := u.GetDotEnvVariable("PORT")
	if port == "" {
		// Pick a random open port
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatalf("Failed to listen on random port: %v", err)
		}
		randomPort := listener.Addr().(*net.TCPAddr).Port
		listener.Close()
		port = strconv.Itoa(randomPort)
	}
	portInt, _ := strconv.Atoi(port)
	return portInt
}
