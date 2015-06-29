package http

import (
	"errors"
	"fmt"
	"github.com/agaffney/livestatus-http/livestatus"
	"net/http"
	"strings"
)

type Options struct {
	ListenPort        int
	ListenAddr        string
	Xinetd            bool
	LivestatusSocket  string
	LivestatusAddress string
}

var options *Options
var ls_endpoint *livestatus.Endpoint

func Start(opts *Options) error {
	options = opts
	// Create livestatus endpoint
	switch {
	case opts.LivestatusAddress != "":
		ls_endpoint = &livestatus.Endpoint{Type: livestatus.TYPE_TCP, Address: opts.LivestatusAddress}
	case opts.LivestatusSocket != "":
		ls_endpoint = &livestatus.Endpoint{Type: livestatus.TYPE_UNIX, Address: opts.LivestatusSocket}
	default:
		return errors.New("You must specify the livestatus endpoint information")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", opts.ListenAddr, opts.ListenPort), nil)

	return nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	table := strings.TrimPrefix(req.URL.Path, "/")
	foo, _ := ls_endpoint.Get(table, []string{})
	w.Write(foo.Body.Bytes())
}
