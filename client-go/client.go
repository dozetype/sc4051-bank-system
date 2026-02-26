package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// ===== Constants =====
const (
	SERVER_IP   = "localhost"
	SERVER_PORT = 2222
	TIMEOUT_MS  = 3000
	BUFFER_SIZE = 512
)

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
		"Register for Updates",
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

	// Start Central UDP listener
	go udpListener(conn, replyChan, callbackChan)

	// Print callbacks if registered
	go func() {
		for cb := range callbackChan {
			fmt.Println("\n📢 CALLBACK:", cb)
		}
	}()

	// Start the program flow
	mainMenu(os.Stdin, conn)
}

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

		if strings.HasPrefix(msg, "8:CALLBACK") {
			callbackChan <- msg
		} else {
			replyChan <- msg
		}
	}
}

// ===== Main Menu =====
func mainMenu(input io.Reader, conn *net.UDPConn) {
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
			handleDelete(input, conn)
		case "3":
			handleDeposit(input, conn)
		case "4":
			handleWithdraw(input, conn)
		case "5":
			handleViewBalance(input, conn)
		case "6":
			handleRegister(input, conn)
		default:
			fmt.Println("Invalid option.")
		}
	}
}

// ===== UDP Send / Receive =====
func sendRequestReceiveReply(conn *net.UDPConn, request string) string {
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Send error:", err)
		return ""
	}

	select {
	case reply := <-replyChan:
		return reply
	case <-time.After(TIMEOUT_MS * time.Millisecond):
		return "Timeout waiting for reply"
	}
}

// ===== Login Handler =====
// Format: 5:LOGIN<userLength>:<username><passLength>:<password>
func handleLogin(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Account Number: ")
	accountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Password: ")
	password, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf("5:LOGIN%d:%s%d:%s",
		len(accountNumber), accountNumber,
		len(password), password)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	fmt.Println("Reply:", reply)
}

// ===== CreateAccount Handler =====
// Format: 13:CREATEACCOUNT<userLength>:<username><passLength>:<password><currencyLength>:<currency><depositLength>:<initialDeposit>
func handleCreateAccount(input io.Reader, conn *net.UDPConn) {
	fmt.Print("New Username: ")
	user, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("New Password: ")
	pass, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Currency: ")
	currency, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Initial Deposit: ")
	deposit, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf("13:CREATEACCOUNT%d:%s%d:%s%d:%s%d:%s",
		len(user), user,
		len(pass), pass,
		len(currency), currency,
		len(deposit), deposit)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== Deletion Handler =====
// Format: 6:DELETE<nameLength>:<name><acctLength>:<accountNumber><passLength>:<password>
func handleDelete(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Account Name: ")
	name, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Account Number: ")
	accountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Password: ")
	password, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf(
		"6:DELETE%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
	)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== Deposit Handler =====
// Format: 7:DEPOSIT<nameLength>:<name><acctLength>:<accountNumber><passLength>:<password><currencyLength>:<currency><amountLength>:<amount>
func handleDeposit(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Account Name: ")
	name, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Account Number: ")
	accountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Password: ")
	password, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Currency: ")
	currency, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Amount to Deposit: ")
	amount, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf("7:DEPOSIT%d:%s%d:%s%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
		len(currency), currency,
		len(amount), amount,
	)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== Withdraw Handler =====
// Format: 8:WITHDRAW<nameLength>:<name><acctLength>:<accountNumber><passLength>:<password><currencyLength>:<currency><amountLength>:<amount>
func handleWithdraw(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Account Name: ")
	name, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Account Number: ")
	accountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Password: ")
	password, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Currency: ")
	currency, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Amount to Withdraw: ")
	amount, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf("8:WITHDRAW%d:%s%d:%s%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
		len(currency), currency,
		len(amount), amount,
	)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== View Balance Handler =====
// Format: 4:VIEW<nameLength>:<name><acctLength>:<accountNumber><passLength>:<password>
func handleViewBalance(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Account Name: ")
	name, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Account Number: ")
	accountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Password: ")
	password, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf(
		"4:VIEW%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
	)

	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== Callback Register =====
// Format: 7:MONITOR<timeLength>:<timeSeconds>
func handleRegister(input io.Reader, conn *net.UDPConn) {
	fmt.Print("Time in Seconds: ")
	timeSeconds, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf("7:MONITOR%d:%s", len(timeSeconds), timeSeconds)
	reply := sendRequestReceiveReply(conn, requestProtocol)
	parseReply(reply)
}

// ===== Helper Functions =====
func readLine(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "Failed to read input", err
	}
	return strings.TrimSpace(input), nil
}

func showMenu(input io.Reader, menu Menu) (string, error) {
	fmt.Println("\n---", menu.Title, "---")
	for i, option := range menu.Options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	fmt.Print("Select: ")

	return readLine(input)
}

func parseReply(reply string) {
	parts := strings.SplitN(reply, ":", 3)

	if len(parts) < 3 {
		fmt.Println("Invalid reply format:", reply)
	}

	status := parts[1]
	message := strings.TrimSpace(parts[2])

	if status == "FAIL" {
		fmt.Println("Error:", message)
	} else {
		fmt.Println("Success:", message)
	}
}

func exit() {
	fmt.Println("\nThank you for using our application!")
	os.Exit(0)
}
