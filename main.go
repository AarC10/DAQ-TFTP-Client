/**
TFTP client that reads and writes config files
*/

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

const broadcastAddress string = "127.0.0.1:69"

var client *http.Client

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createConfig() {
	var srcIP, dstIP, gswIP, subnetIP, srcUDP, adc0UDP, adc1UDP, tcpUDP string

	fmt.Scanln(&srcIP)
	fmt.Scanln(&dstIP)
	fmt.Scanln(&gswIP)
	fmt.Scanln(&subnetIP)

	fmt.Scanln(&srcUDP)
	fmt.Scanln(&adc0UDP)
	fmt.Scanln(&adc1UDP)
	fmt.Scanln(&tcpUDP)

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

func readBytes(conn *net.UDPConn, buffer []byte) {
	_, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Received ", string(buffer), " from ", addr)
}

func sendFile(fileName string) {

	// Open file
	//file, err := os.Open(fileName)
	file, err := ioutil.ReadFile(fileName)
	check(err)
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//
	//	}
	//}(file)

	// Create a buffer to store the file's contents
	//fileInfo, err := file.Stat()
	//check(err)

	// Get file size and create a buffer of that size
	//fileSize := fileInfo.Size()
	//fileData := make([]byte, fileSize)
	//_, err = file.Read(fileData)
	//check(err)

	// Create a UDP connection
	conn, err := net.Dial("udp", broadcastAddress)
	check(err)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Send the file
	_, err = conn.Write(file)
	check(err)
}

func main() {
	client = &http.Client{}
	gui := app.New()

	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	//createConfig()
	//sendFile("config")

	window.ShowAndRun()

}
