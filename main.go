package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{192, 168, 33, 1}, Port: 10001, Zone: ""})
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	t := time.Now().Unix()
	timestamp := strconv.FormatInt(t, 10)

	f, err := os.Create(timestamp + ".log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		n, _, _ := ServerConn.ReadFromUDP(buf)

		_, err = f.Write(buf[5:n])
		if err != nil {
			fmt.Print(err)
		}
		//if buf == loading map { create new file }
	}
}
