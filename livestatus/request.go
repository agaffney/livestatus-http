package livestatus

import (
	"fmt"
	"strings"
	"time"
)

func send_request(req string) {
	fmt.Println(req)
}

func Get(table string, headers []string) {
	req := fmt.Sprintf("GET %s\nOutputFormat: json\n%s\n", table, strings.Join(headers, "\n"))
	send_request(req)
}

func Command(cmd string) {
	req := fmt.Sprintf("COMMAND [%d] %s\n", time.Now().Unix(), cmd)
	send_request(req)
}

func Test() {
	Get("hosts", []string{"Header1: foo", "Header2: bar", "Header3: baz"})
	Get("services", []string{"Header4: foo", "Header5: bar", "Header6: baz"})
	Command("SCHEDULE_SERVICE_DOWNTIME;foo;1;2;some comment")
}
