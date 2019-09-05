/* main function            */
/* file name: webcmd.go     */
/* 		link:               */
/* 		                    */
/*    update: 20181122      */

package webdata

import (
	"encoding/hex"
	"fmt"
	"log"
	"loranet20181205/database"
	"loranet20181205/exception"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func WebCmd(activeStatus int) {

	// var id int
	// var activeStatus int
	// var DevEUI string

	webcmddata := make([]byte, 11)
	webcmddata[0] = 0xFD
	webcmddata[1] = 0xC0
	webcmddata[6] = 0x11
	webcmddata[7] = 0x01
	webcmddata[8] = 0x02

	LoraServerIP, LoraServerPort := database.DbGetServerIP()

	//conn, err := net.Dial("udp", "192.168.10.90:7878")
	//conn, err := net.Dial("udp", "172.16.10.16:7878")
	var IPstring string
	IPstring = fmt.Sprintf("%s:%d", LoraServerIP, LoraServerPort)

	conn, err := net.Dial("udp", IPstring)
	defer conn.Close()
	exception.CheckError(err)

	// Create Connection
	edInfo := database.ActiveEdInfo(activeStatus)
	// Query
	//init delay 29 seconds (26+3)
	webcmddata[9] = 26
	// for rows.Next() {
	// 	err = rows.Scan(&id, &activeStatus, &DevEUI)
	// 	exception.CheckError(err)
	for r := range edInfo {
		log.Print(edInfo[r])
		decoded, err := hex.DecodeString(edInfo[r].DevEUI)
		exception.CheckError(err)

		//dev_addr
		webcmddata[2] = decoded[4]
		webcmddata[3] = decoded[3]
		webcmddata[4] = decoded[2]
		webcmddata[5] = decoded[0]

		//depand on communication speed 3 second per command
		webcmddata[9] += 3

		//check sum
		webcmddata[10] = 0xFD
		for i := 1; i < 10; i++ {
			webcmddata[10] += webcmddata[i]
		}

		//fmt.Printf("HEX: %X", webcmddata)

		// schedule the data
		conn.Write([]byte(webcmddata))
		// fmt.Println(edInfo[r].ID)
		// fmt.Println(edInfo[r].ActiveStatus)
		// fmt.Println(edInfo[r].DevEUI)
		// log.Println("---------------")
		// log.Println(webcmddata)
	}

}
