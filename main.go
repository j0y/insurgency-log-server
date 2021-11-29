package main

import (
	"fmt"
	insurgencylog "my.com/insurgency-log"
	"net"
	"os"
	"strconv"
	"time"
)

var fileWriters = make(map[string]*os.File)

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

	if _, ok := fileWriters[ip]; !ok {
		t := time.Now().Unix()
		timestamp := strconv.FormatInt(t, 10)

		fileWriters[ip], err = os.Create(ip + "_" + timestamp + ".log")
		if err != nil {
			panic(err)
		}
	}

	if message.GetType() == insurgencylog.LoadingMapType {
		mes, ok := message.(insurgencylog.LoadingMap)
		if !ok {
			return
		}

		err = fileWriters[ip].Close()
		if err != nil {
			fmt.Print(err.Error())
		}

		eventTime := mes.Time.Unix()

		fileWriters[ip], err = os.Create(ip + "_" + strconv.FormatInt(eventTime, 10) + "_" + mes.Map + ".log")
		if err != nil {
			fileWriters[ip].Close()
			panic(err)
		}
	}

	_, err = fileWriters[ip].WriteString(text + "\n")
	if err != nil {
		fmt.Print(err.Error())
	}
}
