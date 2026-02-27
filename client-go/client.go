package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// ===== Constants =====
const (
	SERVER_IP   = "192.168.18.15"
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

	fmt.Println("Client local address:", conn.LocalAddr())
	fmt.Println("Server remote address:", conn.RemoteAddr())

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
		fmt.Println("DEBUG RECEIVED:", msg) // debug

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
func sendRequestReceiveReply(conn *net.UDPConn, request string) (string, error) {
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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	parseReply(reply)
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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	parseReply(reply)
}

// ===== Deletion Handler =====
// Format: 12:CLOSEACCOUNT<nameLength>:<name><acctLength>:<accountNumber><passLength>:<password>
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
		"12:CLOSEACCOUNT%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
	)

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

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

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

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
	fields, err := parseFields(reply)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	if len(fields) < 2 {
		fmt.Println("Invalid reply format:", reply)
		return
	}

	status := fields[0]
	message := fields[1]

	if status == "FAIL" {
		fmt.Println("Error:", message)
	} else {
		fmt.Println("Success:", message)
	}
}

func parseFields(data string) ([]string, error) {
	var fields []string
	index := 0

	for index < len(data) {

		start := index
		for index < len(data) && data[index] != ':' {
			index++
		}

		if index >= len(data) {
			return nil, fmt.Errorf("missing colon in length prefix")
		}

		lengthStr := data[start:index]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid length: %s", lengthStr)
		}

		// 2️⃣ Skip colon
		index++

		// 3️⃣ Ensure enough data remains
		if index+length > len(data) {
			return nil, fmt.Errorf("field length exceeds message size")
		}

		// 4️⃣ Extract string
		field := data[index : index+length]
		fields = append(fields, field)

		// 5️⃣ Move index forward
		index += length
	}

	return fields, nil
}

func exit() {
	fmt.Println("\nThank you for using our application!")
	os.Exit(0)
}
