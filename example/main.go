package main

import (
	"bufio"
	"fmt"
	"net"
	"parse-redis/protocol"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handle(conn)
	}

}

func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		command := protocol.Command{}
		data, err := command.Receiver(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("raw byte:", data)
		fmt.Println("command:", command.Value)
		replay := command.Error("error")
		conn.Write(replay)
	}
}
