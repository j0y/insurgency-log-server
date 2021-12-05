package main

import (
	"fmt"
	insurgencylog "my.com/insurgency-log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var files = make(map[string]*os.File, 0)

func main() {
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 16842, Zone: ""})
	if err != nil {
		panic(err)
	}
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		if n < 6 {
			continue
		}
		text := string(buf[5 : n-1])

		write(addr.IP.String(), text)
	}
}

func write(ip string, text string) {
	message, err := insurgencylog.Parse(text)
	if err != nil {
		fmt.Print(err.Error() + ": " + text)
		return
	}

	if message.GetType() == insurgencylog.LoadingMapType {
		mes, ok := message.(insurgencylog.LoadingMap)
		if !ok {
			return
		}

		startNewFile(mes, ip)
	} else if message.GetType() == insurgencylog.ServerMessageType {
		mes, ok := message.(insurgencylog.ServerMessage)
		if !ok {
			return
		} else if mes.Text == "quit" {
			event := insurgencylog.LoadingMap{
				Meta: insurgencylog.Meta{
					Time: mes.Time,
					Type: insurgencylog.LoadingMapType,
				},
				Map: "ministry_coop",
			}
			startNewFile(event, ip)
		}
	} else {
		if _, ok := files[ip]; !ok {
			createFileForNewIP(ip)
		}
	}

	_, err = files[ip].WriteString(text)
	if err != nil {
		fmt.Print(err.Error())
	}
}

func createFileForNewIP(ip string) {
	var err error

	t := time.Now().Unix()
	timestamp := strconv.FormatInt(t, 10)

	newDir := filepath.Join(".", time.Now().Format("2006-01-02"))
	err = os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	files[ip], err = os.Create(newDir + "/" + ip + "_" + timestamp + ".log")
	if err != nil {
		panic(err)
	}
}

func startNewFile(event insurgencylog.LoadingMap, ip string) {
	var err error

	if _, ok := files[ip]; ok {
		err = files[ip].Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}

	eventTime := event.Time.Unix()

	newDir := filepath.Join(".", time.Now().Format("2006-01-02"))
	err = os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	files[ip], err = os.Create(newDir + "/" + ip + "_" + strconv.FormatInt(eventTime, 10) + "_" + event.Map + ".log")
	if err != nil {
		files[ip].Close()
		panic(err)
	}
}
