/**
TFTP client that reads and writes config files
*/

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/rackn/tftp/v3"
	"net/http"
	"os"
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

	for _, v := range configSlice {
		if v[len(v)-1] != '=' {
			configString += v + "\n"
		}
	}

	err := os.WriteFile("config", []byte(configString), 0666)
	check(err)

	fmt.Println(configString)
}

func main() {

	gui := app.New()
	//
	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	broadcastTo := widget.NewEntry()
	broadcastTo.SetText("localhost:69")

	srcIP := widget.NewEntry()
	srcIP.SetPlaceHolder("Source IP")

	dstIP := widget.NewEntry()
	dstIP.SetPlaceHolder("Destination IP")

	gwIP := widget.NewEntry()
	gwIP.SetPlaceHolder("Gateway IP")

	subnetIP := widget.NewEntry()
	subnetIP.SetPlaceHolder("Subnet IP")

	srcUDP := widget.NewEntry()
	srcUDP.SetPlaceHolder("Source UDP Port")

	adc0UDP := widget.NewEntry()
	adc0UDP.SetPlaceHolder("ADC0 UDP Port")

	adc1UDP := widget.NewEntry()
	adc1UDP.SetPlaceHolder("ADC1 UDP Port")

	tcpUDP := widget.NewEntry()
	tcpUDP.SetPlaceHolder("TCP UDP Port")

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
