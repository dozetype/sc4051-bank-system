package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
===== CreateAccount Handler =====
Request Format:
	13:CREATEACCOUNT
	<userLength>:<username>
	<passLength>:<password>
	<currencyLength>:<currency>
	<depositLength>:<initialDeposit>
*/
func handleCreateAccount(input *bufio.Reader, conn *net.UDPConn) {
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

	parseReply(reply, nil)
}

/*
===== CloseAccount Handler =====
Request Format:
	12:CLOSEACCOUNT
	<nameLength>:<name>
	<acctLength>:<accountNumber>
	<passLength>:<password>
*/
func handleCloseAccount(input *bufio.Reader, conn *net.UDPConn) {
	name, accountNumber, password, err := promptVerification(input)
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

	parseReply(reply, nil)
}

/*
===== Deposit Handler =====
Request Format:
	7:DEPOSIT
	<nameLength>:<name>
	<acctLength>:<accountNumber>
	<passLength>:<password>
	<currencyLength>:<currency>
	<amountLength>:<amount>
*/
func handleDeposit(input *bufio.Reader, conn *net.UDPConn) {
	name, accountNumber, password, err := promptVerification(input)
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
	
	context := map[string]string{
		"currency": currency,
	}

	parseReply(reply, context)
}

/*
===== Withdraw Handler =====
Request Format:
	7:DEPOSIT
	<nameLength>:<name>
	<acctLength>:<accountNumber>
	<passLength>:<password>
	<currencyLength>:<currency>
	<amountLength>:<amount>
*/
func handleWithdraw(input *bufio.Reader, conn *net.UDPConn) {
	name, accountNumber, password, err := promptVerification(input)
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
	amount = "-" + amount

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

	context := map[string]string{
		"currency": currency,
	}

	parseReply(reply, context)
}

/*
===== View Balance Handler =====
Request Format:
	4:VIEW
	<nameLength>:<name>
	<acctLength>:<accountNumber>
	<passLength>:<password>
*/
func handleViewBalance(input *bufio.Reader, conn *net.UDPConn) {
	name, accountNumber, password, err := promptVerification(input)
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

	parseReply(reply, nil)
}

/*
===== Transfer Handler =====
Request Format:
	8:TRANSFER
	<nameLength>:<name>
	<acctLength>:<accountNumber>
	<passLength>:<password>
	<currencyLength>:<currency>
	<amountLength>:<amount>
	<targetLength>:<targetAccountNumber>
*/
func handleTransfer(input *bufio.Reader, conn *net.UDPConn) {
	name, accountNumber, password, err := promptVerification(input)
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

	fmt.Print("Amount to Transfer: ")
	amount, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	fmt.Print("Target Account Number: ")
	targetAccountNumber, err := readLine(input)
	if err != nil {
		fmt.Println("Input error:", err)
		return
	}

	requestProtocol := fmt.Sprintf(
		"8:TRANSFER%d:%s%d:%s%d:%s%d:%s%d:%s%d:%s",
		len(name), name,
		len(accountNumber), accountNumber,
		len(password), password,
		len(currency), currency,
		len(amount), amount,
		len(targetAccountNumber), targetAccountNumber,
	)

	reply, err := sendRequestReceiveReply(conn, requestProtocol)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	parseReply(reply, nil)
}

/*
===== Callback Register =====
Request Format:
	7:MONITOR
	<timeLength>:<timeSeconds>
*/
func handleRegister(input *bufio.Reader, conn *net.UDPConn) {
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

	context := map[string]string{
	"seconds": timeSeconds,
	}

	parseReply(reply, context)
}

// ===== Helper Functions =====
func promptInput(input *bufio.Reader, label string) (string, error) {
	fmt.Print(label)

	value, err := readLine(input)
	if err != nil {
		return "", err
	}

	return value, nil
}

func promptVerification(input *bufio.Reader) (string, string, string, error) {

	name, err := promptInput(input, "Account Name: ")
	if err != nil {
		return "", "", "", err
	}

	accountNumber, err := promptInput(input, "Account Number: ")
	if err != nil {
		return "", "", "", err
	}

	password, err := promptInput(input, "Password: ")
	if err != nil {
		return "", "", "", err
	}

	return name, accountNumber, password, nil
}

func parseReply(reply string, context map[string]string) {

	if context == nil {
		context = map[string]string{}
	}

	fields, err := parseFields(reply)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	if len(fields) == 0 {
		fmt.Println("parseReply: No reply.")
		return
	}

	status := fields[0]

	switch status {

	case "FAIL":

		if len(fields) >= 2 {
			fmt.Println("Error:", fields[1])
		} else {
			fmt.Println("Operation Failed")
		}

	case "CREATEACCOUNTSUCCESS":

		if len(fields) >= 2 {
			fmt.Println("Created Successfully.")
			fmt.Println("Account Number:", fields[1])
		} else {
			fmt.Println("Invalid reply format.", reply)
		}

	case "CLOSESUCCESS":

		fmt.Println("Closed Successfully.")

	case "DEPOSITSUCCESS":

		currency := context["currency"]
		if len(fields) >= 2 {
			fmt.Println("Deposit Successful.")
			fmt.Println("New Balance:", fields[1])
		} else {
			fmt.Println("Invalid reply format.", reply, currency)
		}

	case "WITHDRAWSUCCESS":

		currency := context["currency"]
		if len(fields) >= 2 {
			fmt.Println("Withdraw Successful.")
			fmt.Println("New Balance:", fields[1])
		} else {
			fmt.Println("Invalid reply format.", reply, currency)
		}

	case "VIEWSUCCESS":

		if len(fields) < 3 {
			fmt.Println("Balance List Empty")
			return
		}

		fmt.Println("Current Balance:")

		for i := 1; i+1 < len(fields); i += 2 {
			fmt.Printf("%s %s\n", fields[i], fields[i+1])
		}

	case "TRANSFERSUCCESS":
		if len(fields) >= 2 {
			fmt.Println("Transfer Successful.")
		}

		fmt.Println("Current Balance:")

		for i := 1; i+1 < len(fields); i += 2 {
			fmt.Printf("%s %s\n", fields[i], fields[i+1])
		}

	case "MONITORSUCCESS":

		seconds := context["seconds"]
		fmt.Println("Callback Registered Successfully for", seconds, "seconds.")

	case "CALLBACK":

		// Rebuild remaining message after CALLBACK prefix
		if len(fields) > 1 {

			var b strings.Builder
			for _, f := range fields[1:] {
				b.WriteString(strconv.Itoa(len(f)))
				b.WriteByte(':')
				b.WriteString(f)
			}
			remaining := b.String()

			fmt.Println("Callback:", remaining)

			// Parse normally if callback contains protocol reply
			parseReply(remaining, context)

		} else {
			fmt.Println("Invalid reply format.", reply)
		}
	
	default:
		fmt.Println("Server Reply:", status)

		if len(fields) > 1 {
			for _, v := range fields[1:] {
				fmt.Println(v)
			}
		}
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

		index++

		if index+length > len(data) {
			return nil, fmt.Errorf("field length exceeds message size")
		}

		field := data[index : index+length]
		fields = append(fields, field)

		index += length
	}

	return fields, nil
}
