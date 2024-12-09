package main

import (
	"fmt"
	"log"
	"os"

	"github.com/debdutdeb/ipcalc-go/pkg"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("ip with subnet must be passed")
	}

	address := os.Args[1]

	ip, err := pkg.NewIP(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HostMin:", ip.HostMin())
	fmt.Println("HostMax:", ip.HostMax())
}
