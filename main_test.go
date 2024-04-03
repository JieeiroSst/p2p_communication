package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func TestConnectHost(t *testing.T) {
	ctx := context.Background()
	port := 8080

	host1, err := ConnectHost(port, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	defer host1.Close()

	host2, err := ConnectHost(port+1, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	defer host2.Close()

	result := <-ping.Ping(ctx, host1, host2.ID())
	if result.Error != nil {
		t.Errorf("Failed to ping peer: %v", result.Error)
	}
}

func TestConnectSendReceive(t *testing.T) {
	ctx := context.Background()
	port1 := 8081 
	port2 := 8082

	host1, err := ConnectHost(port1, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	defer host1.Close()

	host2, err := ConnectHost(port2, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	defer host2.Close()

	rw, err := ConnectP2P(ctx, host2, fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/p2p/%s", port1, host1.ID()))
	if err != nil {
		t.Fatal(err)
	}

	message := "Hello from host2!"
	fmt.Fprintf(rw, "%s\n", message)
	rw.Flush()

	buf := make([]byte, 1024)
	n, err := rw.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	
	if !bytes.Equal(buf[:n], []byte(message+"\n")) {
		t.Errorf("Expected message: %s, Received: %s", message, string(buf[:n]))
	}
}
