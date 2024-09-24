package lib

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/ssh"
)

var Common_client *ssh.Client

func Set_common_client(c *ssh.Client) {
	Common_client = c
}

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
