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

// ===== Global Menus =====
var startMenuObj = Menu{
	Title: "START MENU: type 'exit' to quit",
	Options: []string{
		"Login",
		"Create New Account",
	},
}

var mainMenuObj = Menu{
	Title: "MAIN MENU: type 'exit' to quit",
	Options: []string{
		"Delete Account",
		"Deposit",
		"Withdraw",
		"Logout",
	},
}

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

	// Start the program flow
	startMenu(os.Stdin, conn)
}

// ===== Start Menu =====
func startMenu(input io.Reader, conn *net.UDPConn) {
	for {
		choice, err := showMenu(input, startMenuObj)
		if err != nil {
			fmt.Println("Input error:", err)
			continue
		}

		if choice == "exit" {
			exit()
		}

		switch choice {
		case "1":
			handleLogin(input, conn)
			// After successful login, go to main menu
			mainMenu(input, conn)
		case "2":
			handleCreateAccount(input, conn)
		default:
			fmt.Println("Invalid option.")
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
			handleDelete(input, conn)
		case "2":
			handleDeposit(input, conn)
		case "3":
			handleWithdraw(input, conn)
		case "4":
			fmt.Println("Logging out...")
			return // Go back to start menu
		default:
			fmt.Println("Invalid option.")
		}
	}
}

// ===== UDP Send / Receive =====
func sendRequestReceiveReply(conn *net.UDPConn, message string) string {
	buffer := make([]byte, BUFFER_SIZE)

	for {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Network Error:", err)
			return ""
		}

		conn.SetReadDeadline(time.Now().Add(TIMEOUT_MS * time.Millisecond))

		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Timeout - retrying...")
				continue
			}
			fmt.Println("Network Error:", err)
			return ""
		}

		return string(buffer[:n])
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
	fmt.Println("Reply:", reply)
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
	fmt.Println("Reply:", reply)
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
	fmt.Println("Reply:", reply)
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
	fmt.Println("Reply:", reply)
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

func exit() {
	fmt.Println("\nThank you for using our application!")
	os.Exit(0)
}
