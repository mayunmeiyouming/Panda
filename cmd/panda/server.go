package panda

import (
	"log"
	"Panda/internal/panda"
)

// Server 是 Panda 的实际入口
func Server(listen string, socks string) {
	log.Println(listen)
	log.Println(socks)

	if listen != "" {
		log.Println("Starting to listen for clients")
		go panda.ListenForSocks(listen)
		log.Fatal(panda.ListenForClients(socks))
	}
}
