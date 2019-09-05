package main

import (
	"log"
	"loranet20181205/exception"
	"loranet20181205/schedule"
	"loranet20181205/server"
)

func main() {

	log.Printf("LoRa Net Server Version 0.1 Start...\n")
	exception.Info.Print("LoRa Net Server Version 0.1 Start...\n")
	/* schedule go routines */
	/* file : schedule.go   */
	schedule.Schedule()

	/* start a UDP server   */
	/* file : udp_server.go */
	server.UDPServer()
	// log.Print(database.DbCheckEdInfo(33554699))

	// ed := database.ActiveEdInfo()
	// for i := range ed {
	// 	log.Print(ed[i])
	// }

}
