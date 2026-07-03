package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				out <- str
				str = string(data[i+1:])
			} else {
				str += string(data)
			}
		}
		if len(str) != 0 {
			out <- str
		}

	}()
	return out
}

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Println("listening on 42069")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("error", err)
			continue
		}

		fmt.Println("connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("connection closed")
	}
}
