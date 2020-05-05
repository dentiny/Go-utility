package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
			case msg := <-messages:
				// Broadcast incoming message to all
				// clients' outgoing message channels.
				for cli := range clients {
					cli <- msg
				}

			case cli := <-entering:
				clients[cli] = true

			case cli := <-leaving:
				delete(clients, cli)
				close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close();
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	// keep receiving from connection with client
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func main() {
	// check argument
	N := len(os.Args)
	if N != 2 {
		fmt.Println("Invalid arguments, should be ./netcat ip:port")
		return
	}

	// set up as server
	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
