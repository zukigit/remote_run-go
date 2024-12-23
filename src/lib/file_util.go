package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zukigit/remote_run-go/src/common"
)

func Get_file_trunc(filepath string, flag int) *os.File {
	file, err := os.OpenFile("hosts.json", flag, 0644)
	if err != nil {
		fmt.Printf("Failed to open  hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}
	return file
}

func Get_hosts_from_jsonfile(jsonfilepath string) {
	var temp_hosts []common.Host_struct

	common.Host_pool = common.Host_pool[:0] // clean readed hosts

	// Open the JSON file
	host_jsonfile := Get_file_trunc(jsonfilepath, os.O_CREATE|os.O_RDONLY)
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

	// Iterate through temp_hosts and create appropriate host type (Linux_host or Windows_host)
	for _, temp_host := range temp_hosts {
		var host common.Host

		// If Host_type is Linux, create a Linux_host
		if *temp_host.Host_type == common.LINUX_AGENT || *temp_host.Host_type == common.LINUX_SERVER {
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
			common.Host_pool = append(common.Host_pool, host)
		} else if *temp_host.Host_type == common.WINDOWS_AGENT {
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
			common.Host_pool = append(common.Host_pool, host)
		}
	}
}

func Set_hosts_to_jsonfile(hosts *[]common.Host, json_filepath string) {
	host_jsonfile := Get_file_trunc(json_filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	defer host_jsonfile.Close()

	encoder := json.NewEncoder(host_jsonfile)

	if err := encoder.Encode(hosts); err != nil {
		fmt.Printf("Failed to encode  hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}
}
