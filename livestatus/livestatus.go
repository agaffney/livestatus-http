package livestatus

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	TYPE_UNIX = 0
	TYPE_TCP  = 1
)

type Endpoint struct {
	Type    int
	Address string
}

type Response struct {
	Code   int64
	Length int64
	Body   bytes.Buffer
}

func process_response(c net.Conn) *Response {
	resp := &Response{}
	var header = make([]byte, 16)
	c.Read(header)
	resp.Code, _ = strconv.ParseInt(string(header[0:3]), 10, 16)
	resp.Length, _ = strconv.ParseInt(strings.TrimSpace(string(header[4:15])), 10, 16)
	io.Copy(&resp.Body, c)
	return resp
}

func (e *Endpoint) Send_request(req string) (*Response, error) {
	var socktype string
	switch {
	case e.Type == TYPE_UNIX:
		socktype = "unix"
	case e.Type == TYPE_TCP:
		socktype = "tcp"
	}
	c, err := net.Dial(socktype, e.Address)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	fmt.Printf("Sending request:\n%s", req)
	_, err = c.Write([]byte(req))
	if err != nil {
		return nil, err
	}
	resp := process_response(c)
	return resp, nil
}

func (e *Endpoint) Get(table string, headers []string) (*Response, error) {
	req := fmt.Sprintf("GET %s\nResponseHeader: fixed16\nOutputFormat: json\n%s\n\n", table, strings.Join(headers, "\n"))
	resp, err := e.Send_request(req)
	return resp, err
}

func (e *Endpoint) Command(cmd string) error {
	req := fmt.Sprintf("COMMAND [%d] %s\nResponseHeader: fixed16\n\n", time.Now().Unix(), cmd)
	_, err := e.Send_request(req)
	return err
}
