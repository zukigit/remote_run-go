package common

import (
	"golang.org/x/crypto/ssh"
)

type Auth struct {
	Username, Password, Hostname string
	Port                         int
	Ssh_client                   *ssh.Client
}

func New_auth(username, password, hostname string, port int, ssh_client *ssh.Client) *Auth {
	return &Auth{
		Username:   username,
		Password:   password,
		Hostname:   hostname,
		Port:       port,
		Ssh_client: ssh_client,
	}
}

// func Get_auth() *Auth {
// 	arg_len := len(os.Args)
// 	var host string
// 	port_int := 22
// 	var err error

// 	// check args
// 	if arg_len == 2 {
// 		host = os.Args[1]
// 	} else if arg_len == 3 {
// 		host = os.Args[1]
// 		port := os.Args[2]

// 		// check port
// 		if port_int, err = strconv.Atoi(port); err != nil {
// 			fmt.Printf("Error: Bad port '%s'\n", port)
// 			os.Exit(1)
// 		}
// 	} else {
// 		fmt.Println("Usage: remote_run user@host <port>")
// 		os.Exit(1)
// 	}

// 	// split user and host
// 	parts := strings.Split(host, "@")
// 	if len(parts) != 2 {
// 		fmt.Println("Usage: remote_run user@host <port>")
// 		os.Exit(1)
// 	}

// 	user := parts[0]
// 	host = parts[1]

// 	// get password
// 	fmt.Printf("%s's password:", os.Args[1])
// 	bytepw, err := term.ReadPassword(int(syscall.Stdin))
// 	fmt.Println() //new line
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 		os.Exit(1)
// 	}
// 	password := string(bytepw)

// 	config := lib.Get_config(user, password)
// 	client, err := lib.Get_client(host+":22", config)
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 		os.Exit(1)
// 	}

// 	return New_auth(user, password, host, port_int, client)
// }
