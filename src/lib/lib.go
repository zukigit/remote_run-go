package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"syscall"

	"github.com/zukigit/remote_run-go/src/common"
	"golang.org/x/term"
)

func Get_res_no(stdout string) (string, error) {
	regex := regexp.MustCompile(`Registry number\s*:\s*\[\d+\]`)

	match := regex.FindString(stdout)
	if match == "" {
		return "", errors.New("no registry number found")
	}

	numberRegex := regexp.MustCompile(`\[(\d+)\]`)
	numberMatch := numberRegex.FindStringSubmatch(match)
	if len(numberMatch) > 1 {
		registryNumber := numberMatch[1]
		return registryNumber, nil
	} else {
		return "", errors.New("could not extract the registry number")
	}
}

// Get_str_str_map_single() gives [string]string map.
func Get_str_str_map(keysAndValues ...string) (map[string]string, error) {
	if len(keysAndValues)%2 != 0 {
		return nil, fmt.Errorf("error: Must provide an even number of arguments (key-value pairs)")
	}
	envs := map[string]string{}
	var key, value string

	for i := 0; i < len(keysAndValues); i += 2 {
		key = keysAndValues[i]
		value = keysAndValues[i+1]
		envs[key] = value
	}

	return envs, nil
}

func Ask_usrinput_string(message_to_show string) string {
	var temp_string string
	fmt.Printf("%s: ", message_to_show)
	fmt.Scan(&temp_string)

	return temp_string
}

func Ask_usrinput_int(message_to_show string) (int, error) {
	var temp_int int

	fmt.Printf("%s: ", message_to_show)

	_, err := fmt.Scan(&temp_int)
	if temp_int == 0 {
		fmt.Println("port is zero")
	}
	if err != nil {
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	if temp_int == 0 {
		fmt.Println("port is zero")
	}

	return temp_int, err
}

func Ask_usrinput_passwd_string(message_to_show string) string {
	fmt.Printf("%s: ", message_to_show)
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Failed in getting password, Error:", err.Error())
		os.Exit(1)
	}

	return string(bytepw)
}

// Get_host returns found host otherwise returns "no such host" error
func Get_host(hosts []common.Host, host_type common.Host_type) (common.Host, error) {
	for _, host := range hosts {
		if host.Get_Host_type() == host_type {
			return host, nil
		}
	}

	return nil, fmt.Errorf("no such host: %s", host_type)
}

// Get_host returns found hosts otherwise returns "no such host" error
func Get_hosts(hosts []common.Host, host_types ...common.Host_type) ([]common.Host, error) {
	var chosen_hosts []common.Host
	found_host := make(map[int]bool)
	var found_host_type bool

	for _, host_type := range host_types {
		found_host_type = false

		for host_idx, host := range hosts {
			if found_host[host_idx] {
				continue
			}

			if host.Get_Host_type() == host_type {
				chosen_hosts = append(chosen_hosts, host)
				found_host[host_idx] = true
				found_host_type = true

				break
			}
		}

		if !found_host_type {
			return nil, fmt.Errorf("no such host: %s", host_type)
		}
	}

	return chosen_hosts, nil
}
