package livestatus

import (
	"fmt"
	"net"
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

func (e *Endpoint) Send_request(req string) (string, error) {
	var socktype string
	switch {
	case e.Type == TYPE_UNIX:
		socktype = "unix"
	case e.Type == TYPE_TCP:
		socktype = "tcp"
	}
	fmt.Println("Dialing...")
	c, err := net.Dial(socktype, e.Address)
	if err != nil {
		return "", err
	}
	fmt.Println("Writing...")
	fmt.Printf("Sending request:\n%s", req)
	_, err = c.Write([]byte(req))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 1024)
	c.Read(buf[:])
	c.Close()
	fmt.Printf("Response:\n%s\n", buf)
	return string(buf), nil
}

func (e *Endpoint) Get(table string, headers []string) {
	req := fmt.Sprintf("GET %s\nOutputFormat: json\n%s\n\n", table, strings.Join(headers, "\n"))
	e.Send_request(req)
}

func (e *Endpoint) Command(cmd string) {
	req := fmt.Sprintf("COMMAND [%d] %s\n\n", time.Now().Unix(), cmd)
	e.Send_request(req)
}