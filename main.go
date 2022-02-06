/**
TFTP client that reads and writes config files.

@author Aaron Chan, RIT Launch Initiative's Future Star Programmer 😎
*/

package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/rackn/tftp/v3"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var client *http.Client

// Struct representing config file entries
type config struct {
	srcIP    *widget.Entry
	dstIP    *widget.Entry
	gwIP     *widget.Entry
	subnetIP *widget.Entry

	srcUDP  *widget.Entry
	adc0UDP *widget.Entry
	adc1UDP *widget.Entry
	tcpUDP  *widget.Entry
}

/**
Error Handling
*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
Uploads files to given broadcast address
*/
func uploadFile(broadcastAddress string) {
	client, err := tftp.NewClient(broadcastAddress)
	check(err)

	file, err := os.Open("config")
	check(err)

	client.SetTimeout(5 * time.Second)
	readFrom, err := client.Send("config", "octet")
	check(err)

	bytesSent, err := readFrom.ReadFrom(file)
	check(err)

	fmt.Printf("%d bytes sent\n", bytesSent)
}

/**
Receives a file from the given broadcast address and inputs it into the text boxes
*/
func receiveFile(broadcastAddress string, configData *config) {
	client, err := tftp.NewClient(broadcastAddress)
	check(err)

	writeTo, err := client.Receive("config", "octet")
	check(err)

	file, err := os.Create("config")
	check(err)

	bytesReceived, err := writeTo.WriteTo(file)
	check(err)

	fileString, err := os.ReadFile("config")
	check(err)

	for _, line := range strings.Split(string(fileString), "\n") {
		if strings.Contains(line, "ip.src") {
			configData.srcIP.SetText(line[7:])
		} else if strings.Contains(line, "ip.dst") {
			configData.dstIP.SetText(line[7:])
		} else if strings.Contains(line, "ip.gw") {
			configData.gwIP.SetText(line[6:])
		} else if strings.Contains(line, "ip.subnet") {
			configData.subnetIP.SetText(line[10:])
		} else if strings.Contains(line, "udp.src") {
			configData.srcUDP.SetText(line[8:])
		} else if strings.Contains(line, "udp.adc0") {
			configData.adc0UDP.SetText(line[9:])
		} else if strings.Contains(line, "udp.adc1") {
			configData.adc1UDP.SetText(line[9:])
		} else if strings.Contains(line, "udp.tcp") {
			configData.tcpUDP.SetText(line[8:])
		}
	}

	fmt.Printf("File Recieved. %d bytes received\n", bytesReceived)
}

/**
Creates a config file in the application directory
*/
func createConfig(configData *config) {
	configString := ""

	// Convert config struct to commands for config file
	configSlice := []string{
		"ip.src=" + configData.srcIP.Text,
		"ip.dst=" + configData.dstIP.Text,
		"ip.gw=" + configData.gwIP.Text,
		"ip.subnet=" + configData.subnetIP.Text,
		"udp.src=" + configData.srcUDP.Text,
		"udp.adc0=" + configData.adc0UDP.Text,
		"udp.adc1=" + configData.adc1UDP.Text,
		"udp.tcp=" + configData.tcpUDP.Text,
	}

	// Iterate over config slice
	for configIndex, configLine := range configSlice {

		// Adds the line to config, if it is not blank and is a valid IP or port
		if configLine[len(configLine)-1] != '=' && validateEntry(configIndex, configData) {
			configString += configLine + "\n"
		}
	}

	err := os.WriteFile("config", []byte(configString), 0666)
	check(err)

	fmt.Println(configString)
}

/**
Validates if a string is an IP Address
*/
func ipAddrValidator(ip string) error {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	if re.Match([]byte(ip)) || ip == "" {
		return nil
	}

	return errors.New("invalid IP Address")

}

/**
Validates if a string is a port
*/
func portValidator(port string) error {
	re := regexp.MustCompile(`^[0-9]{1,5}$`)

	if re.Match([]byte(port)) || port == "" {
		return nil
	}

	return errors.New("invalid port")
}

/**
Factory method for creating a text field
*/
func makeEntryField(text string, validationType string) *widget.Entry {
	newField := widget.NewEntry()
	newField.SetPlaceHolder(text)

	// Set validator function
	if validationType == "ip" {
		newField.Validator = ipAddrValidator
	} else {
		newField.Validator = portValidator
	}

	return newField
}

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
		err = configData.tcpUDP.Validate()

	default:
		err = nil

	}

	if err == nil {
		validated = true
	}

	fmt.Println(err)

	return validated
}

// Main method
func main() {

	// Initialize GUI
	gui := app.New()

	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	// Text field representing address to broadcast to
	broadcastTo := widget.NewEntry()
	broadcastTo.SetText("localhost:69")

	// Entries for config file
	configData := config{
		srcIP:    makeEntryField("Source IP", "ip"),
		dstIP:    makeEntryField("Destination IP", "ip"),
		gwIP:     makeEntryField("Gateway IP", "ip"),
		subnetIP: makeEntryField("Subnet IP", "ip"),
		srcUDP:   makeEntryField("Source UDP Port", "port"),
		adc0UDP:  makeEntryField("ADC0 UDP Port", "port"),
		adc1UDP:  makeEntryField("ADC1 UDP Port", "port"),
		tcpUDP:   makeEntryField("TCP Port", "port"),
	}

	// Set and Display Content

	window.SetContent(
		container.NewVBox(
			broadcastTo,
			configData.srcIP,
			configData.dstIP,
			configData.gwIP,
			configData.subnetIP,
			configData.srcUDP,
			configData.adc0UDP,
			configData.adc1UDP,
			configData.tcpUDP,

			widget.NewButton("Create Config", func() {
				createConfig(&configData)
			}),

			widget.NewButton("Receive File", func() {
				receiveFile(broadcastTo.Text, &configData)
			}),

			widget.NewButton("Upload File", func() {
				uploadFile(broadcastTo.Text)
			}),
		),
	)

	window.ShowAndRun()
}
