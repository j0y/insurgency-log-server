package main

import (
	"fmt"
	insurgencylog "my.com/insurgency-log"
	"net"
	"os"
	"strconv"
	"time"
)

var files = make(map[string]*os.File)

func main() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{192, 168, 33, 1}, Port: 10001, Zone: ""})
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		if n < 6 {
			continue
		}
		text := string(buf[5:n])

		write(addr.IP.String(), text)
	}
}

func write(ip string, text string) {
	message, err := insurgencylog.Parse(text)
	if err != nil {
		fmt.Print(err.Error() + ": " + text)
		return
	}

	if _, ok := files[ip]; !ok {
		createFileForNewIP(ip)
	}

	if message.GetType() == insurgencylog.LoadingMapType {
		mes, ok := message.(insurgencylog.LoadingMap)
		if !ok {
			return
		}

		startNewFile(mes, ip)
	}

	_, err = files[ip].WriteString(text + "\n")
	if err != nil {
		fmt.Print(err.Error())
	}
}

func createFileForNewIP(ip string) {
	var err error

	t := time.Now().Unix()
	timestamp := strconv.FormatInt(t, 10)

	files[ip], err = os.Create(ip + "_" + timestamp + ".log")
	if err != nil {
		panic(err)
	}
}

func startNewFile(event insurgencylog.LoadingMap, ip string) {
	err := files[ip].Close()
	if err != nil {
		fmt.Print(err.Error())
	}

	eventTime := event.Time.Unix()

	files[ip], err = os.Create(ip + "_" + strconv.FormatInt(eventTime, 10) + "_" + event.Map + ".log")
	if err != nil {
		files[ip].Close()
		panic(err)
	}
}
