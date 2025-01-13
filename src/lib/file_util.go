package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zukigit/remote_run-go/src/common"
)

func Get_file_trunc(filepath string, flag int, permission os.FileMode) *os.File {
	file, err := os.OpenFile("hosts.json", flag, permission)
	if err != nil {
		fmt.Printf("Failed to open  hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}
	return file
}

func Get_hosts_from_jsonfile(jsonfilepath string) {
	var temp_hosts []common.Host_struct

	common.Hosts = common.Hosts[:0] // clean readed hosts

	// Open the JSON file
	host_jsonfile := Get_file_trunc(jsonfilepath, os.O_CREATE|os.O_RDONLY, 0644)
	defer host_jsonfile.Close()

	// Decode the JSON file into the temp_hosts slice
	decoder := json.NewDecoder(host_jsonfile)
	if err := decoder.Decode(&temp_hosts); err != nil {
		if err == io.EOF {
			return
		}
		fmt.Printf("Failed to decode hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(temp_hosts) <= 0 {
		fmt.Println("error: no hosts to run, use 'register_hosts' command to register.")
		os.Exit(1)
	}

	// Iterate through temp_hosts and create appropriate host type (Linux_host or Windows_host)
	for _, temp_host := range temp_hosts {
		var host common.Host

		// If Host_type is Linux, create a Linux_host
		if *temp_host.Host_type == common.LA_HOST_TYPE || *temp_host.Host_type == common.LS_HOST_TYPE {
			host = &common.Linux_host{
				Host_name:         temp_host.Host_name,
				Host_run_username: temp_host.Host_run_username,
				Host_ip:           temp_host.Host_ip,
				Host_dns:          temp_host.Host_dns,
				Host_connect_port: temp_host.Host_connect_port,
				Host_use_ip:       temp_host.Host_use_ip,
				Host_type:         temp_host.Host_type,
			}

			// Append the created host to the hosts slice
			common.Hosts = append(common.Hosts, host)
		} else if *temp_host.Host_type == common.WA_HOST_TYPE {
			// If Host_type is Windows, create a Windows_host (assuming you have such a struct)
			host = &common.Windows_host{
				Host_name:         temp_host.Host_name,
				Host_run_username: temp_host.Host_run_username,
				Host_ip:           temp_host.Host_ip,
				Host_dns:          temp_host.Host_dns,
				Host_connect_port: temp_host.Host_connect_port,
				Host_use_ip:       temp_host.Host_use_ip,
				Host_type:         temp_host.Host_type, // Set specific Host_type
			}

			// Append the created host to the hosts slice
			common.Hosts = append(common.Hosts, host)
		}
	}
}

func Set_hosts_to_jsonfile(hosts *[]common.Host, json_filepath string) {
	host_jsonfile := Get_file_trunc(json_filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer host_jsonfile.Close()

	encoder := json.NewEncoder(host_jsonfile)

	if err := encoder.Encode(hosts); err != nil {
		fmt.Printf("Failed to encode  hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}
}
