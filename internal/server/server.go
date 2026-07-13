package server

import (
	"fmt"
	"io"
	"net"

	"github.com/abh-i-navv/httpfromtcp/internal/request"
	"github.com/abh-i-navv/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	closed  bool
	handler Handler
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}

	// writer := bytes.NewBuffer([]byte{})
	s.handler(responseWriter, r)

	// var body []byte = nil
	// status := response.StatusOK

	// if handlerError != nil {
	// 	status = handlerError.StatusCode
	// 	body = []byte(handlerError.Message)
	// } else {
	// 	body = writer.Bytes()
	// }

	// headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

	// response.WriteStatusLine(conn, status)
	// response.WriteHeaders(conn, headers)
	// conn.Write(body)

}

func runServer(s *Server, listener net.Listener) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go runConnection(s, conn)
	}

}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{closed: false, handler: handler}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
