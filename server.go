package main

import (
    "fmt"
    "net"
    "encoding/binary"
)

type MessageType int32
const (
    MSG_CLOSE MessageType = 0
    MSG_TEXT MessageType = 1
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

func HandleClient(conn net.Conn) {
    var err error
    running := true
    for running {
        //First 4 bytes tell us the message type
        msgType := ReadInt32(conn)

        switch msgType {
        case int32(MSG_CLOSE):
            running = false
        case int32(MSG_TEXT):
            //Next 4 bytes tell us the length of the message
            msgLen := ReadInt32(conn)

            //Read the message string
            msg := make([]byte, msgLen)
            bytesRead := 0
            for bytesRead < int(msgLen) {
                n, err := conn.Read(msg[bytesRead:])
                if err != nil {
                    //TODO: this probably isn't good enough
                    fmt.Println(err)
                    running = false
                }
                bytesRead += n
            }

            ProcessMessage(conn.RemoteAddr().String(), string(msg[:msgLen]))

            //Echo it back
            _, err = conn.Write(msg)
            if err != nil {
                fmt.Println(err)
                running = false
            }
        }
    }
    conn.Close()
}

func ReadInt32(conn net.Conn) int32 {
    buffer := make([]byte, 4)
    err := binary.Read(conn, binary.BigEndian, buffer)
    if err != nil {
        fmt.Println(err)
        return 0
    } else {
        return int32(binary.BigEndian.Uint32(buffer))
    }
}

func ProcessMessage(sender string, msg string) {
    fmt.Printf("Received message from %s: %s\n", sender, msg)
}
