package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abh-i-navv/httpfromtcp/internal/request"
	"github.com/abh-i-navv/httpfromtcp/internal/response"
	"github.com/abh-i-navv/httpfromtcp/internal/server"
)

const port = 42069

func response400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func response500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func response200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := response200()
		status := response.StatusOK

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			body = response400()
			status = response.StatusBadRequest

		case "/myproblem":
			body = response500()
			status = response.StatusInternalServerError
		}
		w.WriteStatusLine(status)

		h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
