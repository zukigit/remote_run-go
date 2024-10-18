package lib

import (
	"errors"
	"fmt"
	"regexp"
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
