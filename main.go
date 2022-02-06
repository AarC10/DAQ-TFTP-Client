/**
TFTP client that reads and writes config files
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
	"time"
)

var client *http.Client

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func uploadFile(broadcastAddress string) {
	println(broadcastAddress)
	c, err := tftp.NewClient(broadcastAddress)
	check(err)

	file, err := os.Open("config")
	check(err)

	c.SetTimeout(5 * time.Second)
	rf, err := c.Send("config", "octet")
	check(err)

	n, err := rf.ReadFrom(file)
	check(err)

	fmt.Printf("%d bytes sent\n", n)
}

func getFile(broadcastAddress string) {
	c, err := tftp.NewClient(broadcastAddress)
	check(err)

	wt, err := c.Receive("config", "octet")
	check(err)

	file, err := os.Create("config")
	check(err)

	n, err := wt.WriteTo(file)
	check(err)

	fmt.Printf("File Recieved. %d bytes received\n", n)
}

func createConfig(configData *config) {
	configString := ""

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

	for configIndex, configLine := range configSlice {
		fmt.Println(configIndex)
		if configLine[len(configLine)-1] != '=' && validateEntry(configIndex, configData) {
			configString += configLine + "\n"
		}
	}

	err := os.WriteFile("config", []byte(configString), 0666)
	check(err)

	fmt.Println(configString)
}

func ipAddrValidator(ip string) error {
	re := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.Match([]byte(ip)) || ip == "" {
		return nil
	}

	return errors.New("invalid IP Address")

}

func portValidator(port string) error {
	re := regexp.MustCompile(`^[0-9]{1,5}$`)
	if re.Match([]byte(port)) {
		return nil
	}

	return errors.New("invalid port")
}

func makeEntryField(text string, validationType string) *widget.Entry {
	newField := widget.NewEntry()
	newField.SetPlaceHolder(text)

	if validationType == "ip" {
		newField.Validator = ipAddrValidator
	} else {
		newField.Validator = portValidator
	}

	return newField
}

func validateEntry(configIndex int, configData *config) bool {
	err := errors.New("validateEntry never reaches switch case")
	validated := false

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

func main() {

	gui := app.New()

	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	broadcastTo := widget.NewEntry()
	broadcastTo.SetText("localhost:69")

	srcIP := makeEntryField("Source IP", "ip")
	dstIP := makeEntryField("Destination IP", "ip")
	gwIP := makeEntryField("Gateway IP", "ip")
	subnetIP := makeEntryField("Subnet IP", "ip")
	srcUDP := makeEntryField("Source UDP Port", "port")
	adc0UDP := makeEntryField("ADC0 UDP Port", "port")
	adc1UDP := makeEntryField("ADC1 UDP Port", "port")
	tcpUDP := makeEntryField("TCP Port", "port")

	configData := config{
		srcIP:    srcIP,
		dstIP:    dstIP,
		gwIP:     gwIP,
		subnetIP: subnetIP,
		srcUDP:   srcUDP,
		adc0UDP:  adc0UDP,
		adc1UDP:  adc1UDP,
		tcpUDP:   tcpUDP,
	}

	window.SetContent(
		container.NewVBox(
			broadcastTo,
			srcIP,
			dstIP,
			gwIP,
			subnetIP,
			srcUDP,
			adc0UDP,
			adc1UDP,
			tcpUDP,
			widget.NewButton("Create Config", func() {
				createConfig(&configData)
			}),

			widget.NewButton("Get File", func() {
				getFile(broadcastTo.Text)
			}),

			widget.NewButton("Upload File", func() {
				uploadFile(broadcastTo.Text)
			}),
		),
	)

	window.ShowAndRun()
}
