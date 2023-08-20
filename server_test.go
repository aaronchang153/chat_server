package main

import (
    "encoding/binary"
    "net"
    "testing"
)

func TestServer(t *testing.T) {
    serverSock, clientSock := net.Pipe()

    go HandleClient(serverSock)

    msgStr := "Hello world!"
    data := []byte(msgStr)

    var msgType MessageType = MSG_TEXT
    WriteInt32(clientSock, int32(msgType), t)

    msgLen := int32(len(data))
    WriteInt32(clientSock, msgLen, t)

    clientSock.Write(data)
    var err error
    var n int
    n, err = clientSock.Read(data)
    if err != nil {
        t.Error(err)
    }

    recvStr := string(data[:n])
    if recvStr != msgStr {
        t.Fail()
    }

    msgType = MSG_CLOSE
    WriteInt32(clientSock, int32(msgType), t)
}

func WriteInt32(conn net.Conn, data int32, t *testing.T) {
    err := binary.Write(conn, binary.BigEndian, data)
    if err != nil {
        t.Error(err)
    }
}
