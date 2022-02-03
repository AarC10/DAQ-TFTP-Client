/**
TFTP client that reads and writes config files
*/

package main

type IPConfig struct {
	src string
	dst string
	gw string
	subnet string
}

type UDPConfig struct {
	src string
	adc0 string
	adc1 string
	tc string
}

import (
	"net"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfig() {


}

func writeConfig() {

}

func main() {

}
