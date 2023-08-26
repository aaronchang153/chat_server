package main

import (
	"fmt"
	"net"
)

func main() {
    fmt.Println("Hello world!")

    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }

        fmt.Println("Received connection from", conn.RemoteAddr())
        HandleClient(conn)
    }
}
