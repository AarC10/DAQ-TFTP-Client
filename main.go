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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/rackn/tftp/v3"
	"image/color"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

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

	adc0Rate *widget.Select
	adc1Rate *widget.Select

	resetCheck *widget.Check
}

// Struct representing user settings and extra widgets.
type extras struct { // TODO: Brainstorm a better struct name
	broadcastAddr *widget.Entry
	guiResponses  *widget.Entry
	loadingBar    *widget.ProgressBar
	inputResponse *widget.Label
}

var TFTP_PORT = ":69"

/**
Error Handling
*/
func check(e error, messageLabel **widget.Label) {
	if e != nil {
		(*messageLabel).SetText(e.Error())
		log.Println(e)
	}
}

/**
Test ping an IP before receiving/uploading
*/
func pingCheck(broadcastAddr *widget.Entry, messageLabel **widget.Label) bool {
	pingSuccess := true

	_, err := net.Dial("udp", broadcastAddr.Text)

	if err != nil {
		log.Println("Can't connect")
		(*messageLabel).SetText("Can't connect to DAQ IP Address")
		pingSuccess = false
	}

	return pingSuccess
}

/**
Uploads files to given broadcast address
*/
func uploadFile(extras *extras) {
	extras.loadingBar.SetValue(0)

	if extras.broadcastAddr.Text == "" {
		return
	}

	if !pingCheck(extras.broadcastAddr, &extras.inputResponse) {
		return
	}

	client, err := tftp.NewClient(extras.broadcastAddr.Text + TFTP_PORT)
	if err != nil {
		return
	}

	extras.loadingBar.SetValue(25)

	file, err := os.Open("config")
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(50)

	client.SetTimeout(5 * time.Second)
	readFrom, err := client.Send("config", "octet")
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(75)

	bytesSent, err := readFrom.ReadFrom(file)
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(100)

	fmt.Printf("%d bytes sent\n", bytesSent)

}

/**
Receives a file from the given broadcast address and inputs it into the text boxes
*/
func receiveFile(extras *extras, configData *config) {
	extras.loadingBar.SetValue(0)

	err := extras.broadcastAddr.Validate()
	if err != nil {
		check(err, &extras.inputResponse)
		return
	}

	if !pingCheck(extras.broadcastAddr, &extras.inputResponse) {
		return
	}

	client, err := tftp.NewClient(extras.broadcastAddr.Text + TFTP_PORT)
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(20)

	writeTo, err := client.Receive("config", "octet")
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(40)

	file, err := os.Create("config")
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(60)

	bytesReceived, err := writeTo.WriteTo(file)
	check(err, &extras.inputResponse)

	extras.loadingBar.SetValue(80)

	fileString, err := os.ReadFile("config")
	check(err, &extras.inputResponse)

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

		fmt.Printf("File Recieved. %d bytes received\n", bytesReceived)

	}

	extras.loadingBar.SetValue(100)

}

/**
Creates a config file in the application directory
*/
func createConfig(configData *config, extras *extras) {
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
		"rate.adc0=" + adcRateSelected(configData.adc0Rate),
		"rate.adc1=" + adcRateSelected(configData.adc1Rate),
	}

	// Iterate over config slice
	for configIndex, configLine := range configSlice {

		// Adds the line to config, if it is not blank and is a valid IP or port
		if configLine[len(configLine)-1] != '=' && (validateEntry(configIndex, configData)) {
			configString += configLine + "\n"
		}
	}

	if configData.resetCheck.Checked {
		configString += "reset\n"
	}

	err := os.WriteFile("config", []byte(configString), 0666)
	check(err, &extras.inputResponse)

	fmt.Println(configString)
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

/**
Checks if selection box has a selection and return a string representation
*/
func adcRateSelected(adcRate *widget.Select) string {
	out := ""

	switch adcRate.Selected {
	case "Slow (8 kHz)":
		out = "slow"

	case "Fast (43 kHz)":
		out = "fast"
	}

	return out
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
Factory method for creating texts
*/
func makeNewText(text string) *canvas.Text {
	newText := canvas.NewText(text, color.White)
	newText.Alignment = fyne.TextAlignCenter

	return newText
}

/**
Factory method for creating selection box
*/
func makeNewSelection(text string, options []string) *widget.Select {
	newSelection := widget.NewSelect(options, func(value string) {
		log.Println(text, "set to", value)
	})

	return newSelection
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

	extras := extras{
		makeEntryField("Current DAQ IP", "ip"),
		widget.NewEntry(),
		widget.NewProgressBar(),
		widget.NewLabel(""),
	}

	extras.broadcastAddr.Validator = daqIPAddrValidator

	// Entries for config file
	configData := config{
		srcIP:    makeEntryField("Source IP", "ip"),
		dstIP:    makeEntryField("Destination IP", "ip"),
		gwIP:     makeEntryField("Gateway IP", "ip"),
		subnetIP: makeEntryField("Subnet IP", "ip"),

		srcUDP:  makeEntryField("Source UDP Port", "port"),
		adc0UDP: makeEntryField("ADC0 UDP Port", "port"),
		adc1UDP: makeEntryField("ADC1 UDP Port", "port"),
		tcpUDP:  makeEntryField("TCP Port", "port"),

		adc0Rate: makeNewSelection("ADC0 Rate", []string{"Slow (8 kHz)", "Fast (43 kHz)"}),
		adc1Rate: makeNewSelection("ADC1 Rate", []string{"Slow (8 kHz)", "Fast (43 kHz)"}),

		resetCheck: widget.NewCheck("Reset", func(b bool) {}),
	}

	// Instructions message
	instructionsOne := makeNewText("The first entry should be a valid IP address, otherwise the program will not attempt to send/receive files")
	instructionsTwo := makeNewText("All other entries, if invalid will be marked red and will not be written to the config file")
	instructionsThree := makeNewText("The write file button will write all valid entries as commands into a config file")
	instructionsFour := makeNewText("The upload button will send the config file that is currently in the application directory to the DAQ")
	instructionsFive := makeNewText("The receive button will overwrite the config file in the application directory and input its fields into the application")

	// Set and Display Content
	window.SetContent(
		container.NewVBox(
			extras.broadcastAddr,

			canvas.NewLine(color.White),

			configData.srcIP,
			configData.dstIP,
			configData.gwIP,
			configData.subnetIP,
			configData.srcUDP,
			configData.adc0UDP,
			configData.adc1UDP,
			configData.tcpUDP,

			container.NewHBox(
				widget.NewLabel("ADC0 Rate: "),
				configData.adc0Rate,

				widget.NewLabel("ADC1 Rate: "),
				configData.adc1Rate,

				configData.resetCheck,
			),

			canvas.NewLine(color.White),

			widget.NewButton("Create Config", func() {
				createConfig(&configData, &extras)
			}),

			widget.NewButton("Receive File", func() {
				receiveFile(&extras, &configData)
			}),

			widget.NewButton("Upload File", func() {
				uploadFile(&extras)
			}),

			canvas.NewLine(color.White),

			instructionsOne,
			instructionsTwo,
			instructionsThree,
			instructionsFour,
			instructionsFive,
			extras.loadingBar,
			canvas.NewLine(color.White),
			extras.inputResponse,
		),
	)

	window.ShowAndRun()
}
