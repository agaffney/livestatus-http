package livestatus

import (
	"fmt"
	"strings"
	"time"
)

func (e *Endpoint) Get(table string, headers []string) {
	req := fmt.Sprintf("GET %s\nOutputFormat: json\n%s\n\n", table, strings.Join(headers, "\n"))
	e.Send_request(req)
}

func (e *Endpoint) Command(cmd string) {
	req := fmt.Sprintf("COMMAND [%d] %s\n\n", time.Now().Unix(), cmd)
	e.Send_request(req)
}
