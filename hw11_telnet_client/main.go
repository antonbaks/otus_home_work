package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
	timeout := flag.Duration("timeout", 10*time.Second, "connect timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		return errors.New("incorrect port or hostname")
	}

	address := net.JoinHostPort(args[0], args[1])

	c := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM)
	defer cancel()

	go func() {
		err := c.Send()
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stderr, "...EOF")
		cancel()
	}()

	go func() {
		err := c.Receive()
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stderr, "...Connection was closed by peer")
		cancel()
	}()

	<-ctx.Done()

	return nil
}
