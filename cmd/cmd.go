package cmd

import (
	"flag"
	"fmt"
	"github.com/agaffney/livestatus-http/http"
	"os"
)

func Main() {
	// Parse commandline
	opts := &http.Options{}
	var fork bool
	flag.BoolVar(&fork, "fork", false, "fork into the background")
	flag.StringVar(&opts.ListenAddr, "listenaddr", "0.0.0.0", "address of interface to listen on")
	flag.IntVar(&opts.ListenPort, "listenport", 6558, "port to listen on")
	flag.BoolVar(&opts.Xinetd, "xinetd", false, "service is running under xinetd")
	flag.StringVar(&opts.LivestatusSocket, "livestatus-socket", "", "path to Livestatus UNIX socket")
	flag.StringVar(&opts.LivestatusAddress, "livestatus-address", "", "address (host:port) of livestatus TCP endpoint")
	flag.Parse()

	if opts.Xinetd && fork {
		fmt.Println("You cannot specify both 'fork' and 'xinetd' options at the same time\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if opts.LivestatusSocket == "" && opts.LivestatusAddress == "" {
		fmt.Println("You must specify the livestatus endpoint information\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := http.Start(opts)
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
}
