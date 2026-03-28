package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// ===== Constants =====
const (
	SERVER_IP   = "localhost"
	SERVER_PORT = 2222
	TIMEOUT_MS  = 3000
	BUFFER_SIZE = 512
	PACKET_LOSS_PROBABILITY = 0.1
)

type InvocationMode int

const (
	Base InvocationMode = iota
	AtLeastOnce
	AtMostOnce
)

var currentMode InvocationMode = Base

// ===== Menu Struct =====
type Menu struct {
	Title   string
	Options []string
}

// ===== Menu Display =====
var mainMenuObj = Menu{
	Title: "MAIN MENU: type 'exit' to quit",
	Options: []string{
		"Create Account",
		"Delete Account",
		"Deposit",
		"Withdraw",
		"View Balance",
		"Transfer",
		"Register for Updates",
		"Exit",
	},
}

var invocationMenuObj = Menu{
	Title: "Choose Invocation Semantics:",
	Options: []string{
		"Default",
		"At-Least-Once",
		"At-Most-Once",
	},
}

// ===== UDP Channels =====
var replyChan = make(chan string)
var callbackChan = make(chan string)

// ===== Entry Point =====
func main() {

	// Resolve server address
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", SERVER_IP, SERVER_PORT))
	if err != nil {
		panic(err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Client local address:", conn.LocalAddr())
	fmt.Println("Server remote address:", conn.RemoteAddr())

	// Start Central UDP listener
	go udpListener(conn, replyChan, callbackChan)

	// Print callbacks if registered
	go func() {
		for cb := range callbackChan {
			fmt.Println("\nCALLBACK:", cb)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	// Start the program flow
	invocationMenu(reader)
	mainMenu(reader, conn)
}

// ===== Invocation Menu =====
func invocationMenu(input *bufio.Reader) {
	mode, err := showMenu(input, invocationMenuObj)
	if err != nil {
		fmt.Println("Input error:", err)
	}

	switch mode {

	case "1":
		currentMode = Base

	case "2":
		currentMode = AtLeastOnce

	case "3":
		currentMode = AtMostOnce

	default:
		fmt.Println("Invalid choice, using Default mode")
		currentMode = Base
	}
}

// ===== Main Menu =====
func mainMenu(input *bufio.Reader, conn *net.UDPConn) {
	for {
		choice, err := showMenu(input, mainMenuObj)
		if err != nil {
			fmt.Println("Input error:", err)
			continue
		}

		if choice == "exit" {
			exit()
		}

		switch choice {
		case "1":
			handleCreateAccount(input, conn)
		case "2":
			handleCloseAccount(input, conn)
		case "3":
			handleDeposit(input, conn)
		case "4":
			handleWithdraw(input, conn)
		case "5":
			handleViewBalance(input, conn)
		case "6":
			handleTransfer(input, conn)
		case "7":
			handleRegister(input, conn)
		case "8":
			exit()
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func exit() {
	fmt.Println("\nThank you for using our application!")
	os.Exit(0)
}

// ===== Helper Functions =====
func readLine(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "Failed to read input", err
	}
	return strings.TrimSpace(input), nil
}

func showMenu(input *bufio.Reader, menu Menu) (string, error) {
	fmt.Println("\n---", menu.Title, "---")
	for i, option := range menu.Options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	fmt.Print("Select: ")

	return readLine(input)
}
