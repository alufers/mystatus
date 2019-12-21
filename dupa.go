package main

import (
	"log"
	"net"
)

func main() {
	interfaces, _ := net.Interfaces()
	log.Printf("%+v", interfaces)
}
