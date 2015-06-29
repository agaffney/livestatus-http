package cmd

import (
	"flag"
	"fmt"
	"github.com/agaffney/livestatus-http/livestatus"
	"os"
)

type Options struct {
	ListenPort        int
	ListenAddr        string
	Xinetd            bool
	Fork              bool
	LivestatusSocket  string
	LivestatusAddress string
}

func Main() {
	// Parse commandline
	opts := Options{}
	flag.StringVar(&opts.ListenAddr, "listenaddr", "0.0.0.0", "address of interface to listen on")
	flag.IntVar(&opts.ListenPort, "listenport", 6558, "port to listen on")
	flag.BoolVar(&opts.Xinetd, "xinetd", false, "service is running under xinetd")
	flag.BoolVar(&opts.Fork, "fork", false, "fork into the background")
	flag.StringVar(&opts.LivestatusSocket, "livestatus-socket", "", "path to Livestatus UNIX socket")
	flag.StringVar(&opts.LivestatusAddress, "livestatus-address", "", "address (host:port) of livestatus TCP endpoint")
	flag.Parse()

	// Create livestatus endpoint
	var endpoint *livestatus.Endpoint
	switch {
	case opts.LivestatusAddress != "":
		endpoint = &livestatus.Endpoint{Type: livestatus.TYPE_TCP, Address: opts.LivestatusAddress}
	case opts.LivestatusSocket != "":
		endpoint = &livestatus.Endpoint{Type: livestatus.TYPE_UNIX, Address: opts.LivestatusSocket}
	default:
		fmt.Println("You must specify the livestatus endpoint information\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("%+v\n", endpoint)

}
