package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	arg_len := len(os.Args)
	var host string
	port := "22"

	// check args
	if arg_len == 2 {
		host = os.Args[1]
	} else if arg_len == 3 {
		host = os.Args[1]
		port = os.Args[2]

		// check port
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Printf("Error: Bad port '%s'\n", port)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: remote_run user@host <port>")
		os.Exit(1)
	}

	// split user and host
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		fmt.Println("Usage: remote_run user@host <port>")
		os.Exit(1)
	}

	user := parts[0]
	host = parts[1] + ":" + port

	// get password
	fmt.Printf("%s's password:", os.Args[1])
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() //new line
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	pass := string(bytepw)

	fmt.Println("user", user)
	fmt.Println("host", host)
	fmt.Println("pass", pass)
}
