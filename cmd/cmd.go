package cmd

import (
	"flag"
	"fmt"
)

type Options struct {
	ListenPort int
	ListenAddr string
	Xinetd     bool
	Fork       bool
	Socket     string
	Address    string
}

func Main() {
	opts := Options{}
	// Parse commandline
	flag.StringVar(&opts.ListenAddr, "listenaddr", "0.0.0.0", "address of interface to listen on")
	flag.IntVar(&opts.ListenPort, "listenport", 6558, "port to listen on")
	flag.BoolVar(&opts.Xinetd, "xinetd", false, "service is running under xinetd")
	flag.BoolVar(&opts.Fork, "fork", false, "fork into the background")
	flag.StringVar(&opts.Socket, "socket", "", "path to Livestatus UNIX socket")
	flag.StringVar(&opts.Address, "address", "", "host:port of livestatus endpoint")
	flag.Parse()
	fmt.Printf("%v\n", opts)
}
