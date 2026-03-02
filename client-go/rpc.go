package main

import (
	"fmt"
	"net"
	"strings"
	"time"
	"strconv"
	"math/rand"
)

// ===== UDP Listener =====
func udpListener(conn *net.UDPConn, replyChan, callbackChan chan string) {
	buffer := make([]byte, 2048)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}

		msg := string(buffer[:n])
		fmt.Println("DEBUG RECEIVED:", msg) // debug

		if strings.HasPrefix(msg, "8:CALLBACK") {
			callbackChan <- msg
		} else {
			replyChan <- msg
		}
	}
}

// ===== UDP Send / Receive =====
func sendRequestReceiveReply(conn *net.UDPConn, request string) (string, error) {

	retriesLimit := 10
	fmt.Println("\nChoose Invocation Semantics:")
	fmt.Println("Default: 0")
	fmt.Println("At-Least-Once: 1")
	fmt.Println("At-Most-Once: 2")
	fmt.Print("Enter choice: ")

	var mode int
	fmt.Scanln(&mode)

	switch mode {

	case 0:
		return defaultInvocation(conn, request)

	case 1:
		return atLeastOnce(conn, request, retriesLimit)

	case 2:
		return atMostOnce(conn, request, retriesLimit)

	default:
		return "", fmt.Errorf("invalid invocation mode")
	}
}

// ===== Invocation Semantics =====
func defaultInvocation(conn *net.UDPConn, request string) (string, error) {
	packetLossProbability := 0.3

	if rand.Float64() < packetLossProbability {
		fmt.Println("Simulating packet loss for request:", request)
		return "", fmt.Errorf("Simulated packet loss.")
	}
	
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Send error:", err)
		return "", err
	}

	select {
	case reply := <-replyChan:
		return reply, nil
	case <-time.After(TIMEOUT_MS * time.Millisecond):
		return "Timeout waiting for reply", fmt.Errorf("timeout")
	}
}

func atLeastOnce(conn *net.UDPConn, request string, retries int) (string, error) {

	baseDelay := 100 * time.Millisecond

	fmt.Println("\nInitial request sent.")

	for i := 0; i < retries; i++ {

		reply, err := defaultInvocation(conn, request)
		if err == nil {
			return reply, nil
		}

		// If more retries are allowed, wait before next retransmission
		if i < retries-1 {

			delay := baseDelay * time.Duration(1<<i)

			fmt.Printf("Retransmission %d will be sent in %v\n", i+1, delay)

			time.Sleep(delay)

			fmt.Printf("\nRetransmission %d sent.\n", i+1)
		}
	}

	return "", fmt.Errorf("Max retries reached")
}

// <length>:<requestID> appended to the front of request
func atMostOnce(conn *net.UDPConn, request string, retries int) (string, error) {

	requestID := strconv.FormatInt(time.Now().UnixNano(), 10)

	fullRequest := fmt.Sprintf("%d:%s%s",
		len(requestID), requestID, request)

	fmt.Println("DEBUG: " + fullRequest) // debug

	baseDelay := 100 * time.Millisecond

	fmt.Println("\nInitial request sent.")

	for i := 0; i < retries; i++ {

		reply, err := defaultInvocation(conn, request)
		if err == nil {
			return reply, nil
		}

		// If more retries are allowed, wait before next retransmission
		if i < retries-1 {

			delay := baseDelay * time.Duration(1<<i)

			fmt.Printf("Retransmission %d will be sent in %v\n", i+1, delay)

			time.Sleep(delay)

			fmt.Printf("\nRetransmission %d sent.\n", i+1)
		}
	}

	return "", fmt.Errorf("Max retries reached")
}