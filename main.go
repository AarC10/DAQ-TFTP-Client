/**
TFTP client that reads and writes config files
*/

package main

import (
	"fmt"
	"net"
	"os"
)

const broadcastAddress string = "255.255.255.255.69"

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

func sendFile(fileName string) {

	// Open file
	file, err := os.Open(fileName)
	check(err)
	defer file.Close()

	// Create a buffer to store the file's contents
	fileInfo, err := file.Stat()
	check(err)

	// Get file size and create a buffer of that size
	fileSize := fileInfo.Size()
	fileData := make([]byte, fileSize)
	_, err = file.Read(fileData)
	check(err)

	// Create a UDP connection
	conn, err := net.Dial("udp", broadcastAddress)
	check(err)
	defer conn.Close()

	// Send the file
	_, err = conn.Write(fileData)
	check(err)
}

func main() {
	sendFile("config")

}
