package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// check argument
	N := len(os.Args)
	if N != 2 {
		fmt.Println("Invalid arguments, should be ./netcat ip:port")
		return
	}

	// try connect server
	conn, err := net.Dial("tcp", os.Args[1])
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // read from server and print it on stdout, NOTE: ignoring errors
		fmt.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin) // read from stdin, and send it to server
	<-done // wait for background goroutine to finish
}

// read from src and send it to dst
// keep reading, until a EOF occurs or an error happens
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}