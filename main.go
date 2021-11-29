package main

import (
	"fmt"
	insurgencylog "my/insurgency-log"
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
		n, addr, _ := ServerConn.ReadFromUDP(buf)

		_, err = f.Write(buf[5:n])
		if err != nil {
			fmt.Print(err.Error())
		}
		message, err := insurgencylog.Parse(string(buf[5:n]))
		if err != nil {
			fmt.Print(err.Error())
		}

		if message.GetType() == insurgencylog.LoadingMapType {
			mes, ok := message.(*insurgencylog.LoadingMap)
			if !ok {
				continue
			}

			err = f.Close()
			if err != nil {
				fmt.Print(err.Error())
			}

			eventTime := mes.Time.Unix()
			f, err = os.Create(addr.IP.String() + "_" + strconv.FormatInt(eventTime, 10) + "_" + mes.Map + ".log")
			if err != nil {
				panic(err)
			}
		}
	}
}
