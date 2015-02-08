// Author: Marcin Matlaszek https://github.com/emate

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ActiveState/tail"
)

func ReadFile(filename string, stream chan<- string) {
	seekinfo := tail.SeekInfo{Whence: os.SEEK_END}
	t, _ := tail.TailFile(filename, tail.Config{Follow: true, Location: &seekinfo})
	for line := range t.Lines {
		stream <- line.Text + "\n"
	}
}

func client(conn net.Conn, c <-chan string) {
	for msg := range c {
		conn.Write([]byte(msg))
	}
}
func server(channels map[chan string]bool, input <-chan string) {
	for {
		fmt.Println("Servering!")
		for msg := range input {
			for k, _ := range channels {
				fmt.Println("Sending")
				k <- msg
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage %s <FILENAME>", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	address := flag.String("l", "localhost:8080", "Listen address")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Filepath argument missing")
		usage()
	}
	var filename string = args[0]

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no such file : %s", filename)
		usage()
	}

	connections := make(map[net.Conn]bool)
	channels := make(map[chan string]bool)
	filestream := make(chan string)

	fmt.Println(*address)
	l, err := net.Listen("tcp", *address)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	go ReadFile(filename, filestream)
	go server(channels, filestream)
	for {
		c, _ := l.Accept()
		defer c.Close()
		channel := make(chan string)
		channels[channel] = true
		connections[c] = true
		go client(c, channel)
	}
}
