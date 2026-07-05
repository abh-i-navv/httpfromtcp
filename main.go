package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abh-i-navv/httpfromtcp/request"
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

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.Method)
		fmt.Printf("- Target: %s\n", r.RequestTarget)
		fmt.Printf("- Version: %s\n", r.HttpVersion)
		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})

		fmt.Println("connection closed")
	}
}
