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
	var body bytes.Buffer
	io.Copy(&body, req.Body)
	var req_params map[string]interface{}
	json.Unmarshal(body.Bytes(), &req_params)
	fmt.Printf("%+v\n", req_params)
	var headers []string
	if foo, ok := req_params["headers"]; ok {
		for _, header := range foo.([]interface{}) {
			headers = append(headers, header.(string))
		}
	}
	resp, _ := ls_endpoint.Get(table, headers)
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
