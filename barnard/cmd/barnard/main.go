package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/bontibon/gumble/barnard"
	"github.com/bontibon/gumble/barnard/uiterm"
	"github.com/bontibon/gumble/gumble"
	"github.com/bontibon/gumble/gumble_openal"
)

func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")

	flag.Parse()

	// Initialize
	b := barnard.Barnard{}
	b.Ui = uiterm.New(&b)

	// Gumble
	b.Config = gumble.Config{
		Username: *username,
		Address:  *server,
		Listener: &b,
	}
	if *insecure {
		b.Config.TlsConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		if cert, err := tls.LoadX509KeyPair(*certificate, *certificate); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		} else {
			b.Config.TlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	b.Client = gumble.NewClient(&b.Config)
	// Audio
	if stream, err := gumble_openal.New(b.Client); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		b.Config.AudioListener = stream
		b.Stream = stream
	}

	if err := b.Client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	os.Stderr.Close()
	b.Ui.Run()
}
