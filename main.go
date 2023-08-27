package main

import (
    "context"
	"log"
	"net"
	"net/http"
	"time"
    "os"
    "os/signal"
)

func main() {
    log.SetOutput(os.Stdout)
    err := run()
    if err != nil {
        log.Fatal(err)
    }
}

func run() error {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        return err
    }

    server := &http.Server{
        Handler: ChatServer{logf: log.Printf},
        ReadTimeout: time.Second * 10,
        WriteTimeout: time.Second * 10,
    }
    errc := make(chan error, 1)
    go func() {
        errc <- server.Serve(listener)
    }()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt)
    select {
    case err := <-errc:
        log.Printf("failed to serve: %v", err)
    case sig := <-sigs:
        log.Printf("terminating: %v", sig)
    }

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    return server.Shutdown(ctx)
}
