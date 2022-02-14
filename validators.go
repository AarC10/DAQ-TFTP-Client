/**
Validator methods for user entries

@Author Aaron Chan
*/

package main

import (
	"errors"
	"fmt"
	"regexp"
)

/**
Ensures text input will be valid before writing it to a file
*/
func validateEntry(configIndex int, configData *config) bool {
	err := errors.New("validateEntry never reaches switch case")
	validated := false

	/**
	Cases represent the "index" of the configDa struct
	*/

	switch configIndex {

	case 0:
		err = configData.srcIP.Validate()

	case 1:
		err = configData.dstIP.Validate()

	case 2:
		err = configData.gwIP.Validate()

	case 3:
		err = configData.subnetIP.Validate()

	case 4:
		err = configData.srcUDP.Validate()

	case 5:
		err = configData.adc0UDP.Validate()

	case 6:
		err = configData.adc1UDP.Validate()

	case 7:
		err = configData.tcUDP.Validate()

	default:
		err = nil

	}

	if err == nil {
		validated = true
	}

	fmt.Println(err)

	return validated
}

/**
Validates if a string is an IP Address
*/
func ipAddrValidator(ip string) error {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	if re.Match([]byte(ip)) || ip == "" || ip == "localhost" {
		return nil
	}

	return errors.New("error: Invalid IP Address")
}

/**
Validates DAQ IP Address
*/
func daqIPAddrValidator(ip string) error {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	if re.Match([]byte(ip)) || ip == "localhost" {
		return nil
	}

	return errors.New("error: DAQ IP Address must be filled in")
}

/**
Validates if a string is a port
*/
func portValidator(port string) error {
	re := regexp.MustCompile(`^[0-9]{1,5}$`)

	if re.Match([]byte(port)) || port == "" {
		return nil
	}

	return errors.New("error: Invalid port")
}
