package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var host string
	var destination string

	fmt.Scanln(&host)
	fmt.Scanln(&destination)

	fmt.Println(host)
	fmt.Println(destination)

	configText := "ip.src=" + host + "\nip.dst=" + destination
	err := os.WriteFile("config", []byte(configText), 0666)
	check(err)

}
