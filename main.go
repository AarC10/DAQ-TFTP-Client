/**
TFTP client that reads and writes config files
*/

package main

import (
	"fmt"
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

	file, err := os.Open("get_this_file.txt")
	check(err)

	c.SetTimeout(5 * time.Second) // optional
	rf, err := c.Send("get_this_file.txt", "octet")
	check(err)

	n, err := rf.ReadFrom(file)
	check(err)

	fmt.Printf("%d bytes sent\n", n)
}

func getFile() {
	c, err := tftp.NewClient(broadcastAddress)
	check(err)

	wt, err := c.Receive("get_this_file.txt", "octet")
	check(err)

	file, err := os.Create("get_this_file.txt")
	check(err)

	n, err := wt.WriteTo(file)
	check(err)

	fmt.Printf("File Recieved. %d bytes received\n", n)
}

func createConfig() {
	//var srcIP, dstIP, gswIP, subnetIP, srcUDP, adc0UDP, adc1UDP, tcpUDP string
	//
	//fmt.Scanln(&srcIP)
	//fmt.Scanln(&dstIP)
	//fmt.Scanln(&gswIP)
	//fmt.Scanln(&subnetIP)
	//
	//fmt.Scanln(&srcUDP)
	//fmt.Scanln(&adc0UDP)
	//fmt.Scanln(&adc1UDP)
	//fmt.Scanln(&tcpUDP)
	//
	//configString := "ip.src=" + srcIP + "\n" +
	//	"ip.dst=" + dstIP + "\n" +
	//	"ip.gw=" + gswIP + "\n" +
	//	"ip.subnet=" + subnetIP + "\n" +
	//	"udp.src=" + srcUDP + "\n" +
	//	"udp.adc0=" + adc0UDP + "\n" +
	//	"udp.adc1=" + adc1UDP + "\n" +
	//	"udp.tc=" + tcpUDP + "\n"
	//
	//err := os.WriteFile("config", []byte(configString), 0666)
	//if err != nil {
	//	return
	//}

}

func sendFile(fileName string) {

}

func printButton(config widget.Form) {
	fmt.Println(config.Items)

}

func main() {

	//gui := app.New()
	//
	//window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	//window.Resize(fyne.NewSize(1920, 500))
	//
	//config := widget.NewForm(
	//	widget.NewFormItem("IP Source", widget.NewEntry()),
	//	widget.NewFormItem("IP Destination", widget.NewEntry()),
	//	widget.NewFormItem("IP Gateway", widget.NewEntry()),
	//	widget.NewFormItem("IP Subnet", widget.NewEntry()),
	//	widget.NewFormItem("UDP Source", widget.NewEntry()),
	//	widget.NewFormItem("UDP ADC0", widget.NewEntry()),
	//	widget.NewFormItem("UDP ADC1", widget.NewEntry()),
	//	widget.NewFormItem("UDP TCP", widget.NewEntry()),
	//)
	//
	//writeButton := widget.NewButton("Write", func() {
	//	fmt.Println(config)
	//})
	//
	//readButton := widget.NewButton("Read", func() {
	//
	//})
	//
	////createConfig()
	////sendFile("config")
	//
	//window.SetContent(
	//	container.NewVBox(
	//		config,
	//		writeButton,
	//		readButton,
	//	),
	//)
	//window.ShowAndRun()

	//getFile()

	uploadFile()
}
