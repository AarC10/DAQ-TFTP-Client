/**
Functions for each of the buttons in the GUI
@Author Aaron Chan
*/

package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/tatsushid/go-fastping"
	"gitlab.com/rackn/tftp/v3"
	"net"
	"os"
	"strings"
	"time"
)

/**
Test ping an IP before receiving/uploading
*/
func pingCheck(broadcastAddr *widget.Entry, messageLabel **widget.Label) bool {
	pingSuccess := false
	pinger := fastping.NewPinger()

	resolvedIPAddr, err := net.ResolveIPAddr("ip4:icmp", broadcastAddr.Text)
	check(err, messageLabel)

	pinger.AddIPAddr(resolvedIPAddr)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		pingSuccess = true
	}
	pinger.OnIdle = func() {
		fmt.Println("Can't ping current IP")
	}

	err = pinger.Run()
	check(err, messageLabel)

	return pingSuccess
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
		"udp.tc=" + configData.tcUDP.Text,
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
	extras.inputResponse.SetText("Config successfully created.")
}

/**
Receives a file from the given broadcast address and inputs it into the text boxes
*/
func receiveFile(configData *config, extras *extras) {
	extras.loadingBar.SetValue(0)

	err := extras.broadcastAddr.Validate()
	if err != nil {
		check(err, &extras.inputResponse)
		return
	}

	if !pingCheck(extras.broadcastAddr, &extras.inputResponse) {
		extras.inputResponse.SetText("Unable to ping current DAQ IP")
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

	// TODO: Figure out a cleaner way of implementing this
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
		} else if strings.Contains(line, "udp.tc") {
			configData.tcUDP.SetText(line[7:])
		} else if strings.Contains(line, "rate.adc0") {
			if (line[10:]) == "slow" {
				configData.adc0Rate.SetSelectedIndex(0)
			} else if (line[10:]) == "fast" {
				configData.adc0Rate.SetSelectedIndex(1)
			}

			configData.adc0Rate.SetSelected("")
		} else if strings.Contains(line, "rate.adc1") {
			if (line[10:]) == "slow" {
				configData.adc1Rate.SetSelectedIndex(0)
			} else if (line[10:]) == "fast" {
				configData.adc1Rate.SetSelectedIndex(1)
			}
		}

	}
	fmt.Printf("File Recieved. %d bytes received\n", bytesReceived)

	extras.loadingBar.SetValue(100)
	extras.inputResponse.SetText("File successfully received.")
}

/**
Uploads files to given broadcast address
*/
func uploadFile(extras *extras) {
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
	extras.inputResponse.SetText("File successfully uploaded.")

}
