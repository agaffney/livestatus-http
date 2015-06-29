package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agaffney/livestatus-http/livestatus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Options struct {
	ListenPort        int
	ListenAddr        string
	Xinetd            bool
	LivestatusSocket  string
	LivestatusAddress string
}

type RequestBody struct {
	Headers []string `json:"headers"`
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

	http.HandleFunc("/command", command_handler)
	http.HandleFunc("/", query_handler)
	http.ListenAndServe(fmt.Sprintf("%s:%d", opts.ListenAddr, opts.ListenPort), nil)

	return nil
}

func command_handler(w http.ResponseWriter, req *http.Request) {
	var body bytes.Buffer
	io.Copy(&body, req.Body)
	err := ls_endpoint.Command(body.String())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
	}
}

func query_handler(w http.ResponseWriter, req *http.Request) {
	table := strings.TrimPrefix(req.URL.Path, "/")
	var body bytes.Buffer
	io.Copy(&body, req.Body)
	var req_body RequestBody
	json.Unmarshal(body.Bytes(), &req_body)
	resp, _ := ls_endpoint.Get(table, req_body.Headers)
	w.Header().Set("Content-Length", strconv.FormatInt(resp.Length, 10))
	switch {
	case resp.Code == 200:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	case resp.Code == 404:
		w.WriteHeader(404)
	default:
		w.WriteHeader(400)
	}
	w.Write(resp.Body.Bytes())
}
