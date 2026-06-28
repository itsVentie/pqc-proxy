package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pqc-proxy/internal/network"
)

func main() {
	mode := flag.String("mode", "", "Running mode: client or server")
	listenAddr := flag.String("listen", "", "Address to listen on")
	targetAddr := flag.String("target", "", "Target remote address")

	flag.Parse()

	if *mode != "client" && *mode != "server" {
		fmt.Println("Usage of pqc-proxy:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *listenAddr == "" || *targetAddr == "" {
		log.Fatalf("Error: both -listen and -target flags are required")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if *mode == "server" {
		srv := network.NewServer(*listenAddr, *targetAddr)
		fmt.Printf("[+] Starting PQC SERVER on %s -> Forwarding to %s\n", *listenAddr, *targetAddr)
		go func() {
			if err := srv.Start(); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}()
		<-sigChan
		fmt.Println("\n[-] Shutting down server...")
		srv.Stop()
	} else {
		cli := network.NewClient(*listenAddr, *targetAddr)
		fmt.Printf("[+] Starting PQC CLIENT on %s -> Tunneling to %s\n", *listenAddr, *targetAddr)
		go func() {
			if err := cli.Start(); err != nil {
				log.Fatalf("Client error: %v", err)
			}
		}()
		<-sigChan
		fmt.Println("\n[-] Shutting down client...")
		cli.Stop()
	}
	fmt.Println("[+] Off.")
}
