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
}

func Main() {
	opts := Options{}
	// Parse commandline
	flag.StringVar(&opts.ListenAddr, "listenaddr", "0.0.0.0", "address of interface to listen on")
	flag.IntVar(&opts.ListenPort, "listenport", 6558, "port to listen on")
	flag.BoolVar(&opts.Xinetd, "xinetd", false, "service is running under xinetd")
	flag.BoolVar(&opts.Fork, "fork", false, "fork into the background")
	flag.Parse()
	fmt.Printf("%v\n", opts)
}
