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

const broadcastAddress string = "localhost:69"

var client *http.Client

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func uploadFile() {
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

func getFile() {
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

func createConfig(srcIP, dstIP, gswIP, subnetIP, srcUDP, adc0UDP, adc1UDP, tcpUDP string) {

	configString := "ip.src=" + srcIP + "\n" +
		"ip.dst=" + dstIP + "\n" +
		"ip.gw=" + gswIP + "\n" +
		"ip.subnet=" + subnetIP + "\n" +
		"udp.src=" + srcUDP + "\n" +
		"udp.adc0=" + adc0UDP + "\n" +
		"udp.adc1=" + adc1UDP + "\n" +
		"udp.tc=" + tcpUDP + "\n"

	err := os.WriteFile("config", []byte(configString), 0666)

	if err != nil {
		return
	}

}

func main() {

	gui := app.New()
	//
	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	srcIP := widget.NewEntry()
	dstIP := widget.NewEntry()
	gwIP := widget.NewEntry()
	subnetIP := widget.NewEntry()
	srcUDP := widget.NewEntry()
	adc0UDP := widget.NewEntry()
	adc1UDP := widget.NewEntry()
	tcpUDP := widget.NewEntry()

	window.SetContent(
		container.NewVBox(
			srcIP,
			dstIP,
			gwIP,
			subnetIP,
			srcUDP,
			adc0UDP,
			adc1UDP,
			tcpUDP,
			widget.NewButton("Create Config", func() {
				createConfig(srcIP.Text, dstIP.Text, gwIP.Text, subnetIP.Text, srcUDP.Text, adc0UDP.Text, adc1UDP.Text, tcpUDP.Text)
			}),

			widget.NewButton("Get File", func() {
				getFile()
			}),

			widget.NewButton("Upload File", func() {
				uploadFile()
			}),
		),
	)
	window.ShowAndRun()
}
