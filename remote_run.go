package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: remote_run user@host")
		os.Exit(1)
	}

	parts := strings.Split(os.Args[1], "@")
	if len(parts) != 2 {
		fmt.Println("Usage: remote_run user@host")
		os.Exit(1)
	}
	user := parts[0]
	host := parts[1]

	fmt.Printf("%s's password:", os.Args[1])
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	pass := string(bytepw)

	fmt.Println("user", user)
	fmt.Println("host", host)
	fmt.Println("pass", pass)
}
