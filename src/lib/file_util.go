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

func Get_hosts_from_jsonfile(jsonfilepath string) *[]common.Host {
	temp_hosts := make([]common.Host, 0)
	host_jsonfile := Get_file_trunc(jsonfilepath, os.O_CREATE|os.O_RDONLY)
	defer host_jsonfile.Close()

	decoder := json.NewDecoder(host_jsonfile)

	if err := decoder.Decode(&temp_hosts); err != nil {
		if err == io.EOF {
			fmt.Println("it is eof error")
			return &temp_hosts
		}
		fmt.Printf("Failed to decode  hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}
	return &temp_hosts
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
