package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"loranet20181205/database"
	"loranet20181205/exception"
	"loranet20181205/packet"
	"net"
	"time"
	//_ "github.com/go-sql-driver/mysql"
)

// UDPserver listen on UDP port for gateway and web
func UDPServer() {
	LoraServerIP, LoraServerPort := database.DbGetServerIP()
	/* Now listen at selected port */
	//	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("192.168.10.90"), Port: 7878})
	//	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("172.16.10.16"), Port: 7878})
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(LoraServerIP), Port: LoraServerPort})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Local server IP: <%s> \n", ServerConn.LocalAddr().String())
	exception.CheckError(err)
	exception.LogFile("Local server IP: " + ServerConn.LocalAddr().String() + "\n")
	defer ServerConn.Close()

	/* udp server loop */
	pullack := make([]byte, 4)
	pullack[0] = 0x01
	pullack[3] = 0x04

	pushack := make([]byte, 4)
	pushack[0] = 0x01
	pushack[3] = 0x01

	//PULL_RESP packet
	//pullresp := make([]byte, 4)
	//pullresp[0] = 0x01
	//pullresp[3] = 0x03

	//var NwKSKeyB = make([]byte, 16)
	//var AppSKeyB = make([]byte, 16)

	//var DefaultGateway string
	var DefaultGateway = database.DbGetGatewayIP()

	buf := make([]byte, 2048)
	for {
		n, remoteAddr, err := ServerConn.ReadFromUDP(buf)
		s := fmt.Sprintf("RX from %s :%X \n", remoteAddr, buf[0:n])
		log.Printf(s)        //to terminal
		exception.LogFile(s) //log to file

		if err != nil {
			log.Println("Error: ", err)

		}

		//manipulate different type of network data
		switch n {
		case 11: //Web command
			if buf[0] == 0xFD {
				log.Printf("Web: %d %X\n", n, buf[0:n])

				var DevAddrI = int(buf[2]) | int(buf[3])<<8 | int(buf[4])<<16 | int(buf[5])<<24

				if database.DbCheckEdInfo(DevAddrI) == true {
					// log.Print("-----------------------------------------------------------")
					//fmt.Printf("devAddr: %d \n", DevAddrI)
					NwkSKeyB, AppSKeyB, DownlinkCnt, InvokeID := database.DbParameters(DevAddrI)
					var InvokeIDB = make([]byte, 4)
					binary.LittleEndian.PutUint32(InvokeIDB, uint32(InvokeID-1)) //carefule the invokedId -1
					//coding type: ex get set action
					CodingType := uint32(buf[1])
					ObsIndex := uint32(buf[6]) | uint32(buf[7])<<8
					//deviceEUI
					//DevEUI_b := []byte{buf[2], buf[3], buf[4], 0x00, buf[5]}

					switch CodingType {
					case 0xC0: // get request
						switch ObsIndex {
						case 0: //FAN
							log.Printf("fan:%X \n", ObsIndex)
						case 41: //get current datetime
							var DataB = make([]byte, 0)
							log.Printf("getdate:%X \n", ObsIndex)
							FrameCmd := []byte{0xc0, 0x09, 0x00, 0x29, 0x02, 0x01, 0x46, 0x02}

							FrameCmd[1] = InvokeIDB[0]
							FrameB := packet.BuildFrame(DevAddrI, FrameCmd, DataB, DownlinkCnt)
							log.Printf("Webf    :%X \n", FrameB)
							//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
							//Pmsg = MHDR | MACPayload cmac = aes128_cmac(NwkSKey, B0 | Pmsg) page 45
							TxtGetPayload := packet.SecureFrame(FrameB, NwkSKeyB, AppSKeyB, DevAddrI, DownlinkCnt)
							log.Printf("Webplain:%X\n", TxtGetPayload)

							LoraPktPayloadB := packet.LoraPktEnc(TxtGetPayload, 0)
							log.Printf("Webframe:%X\n", LoraPktPayloadB)
							JSONPktPayloadB := []byte(packet.BuildTXJsonPkt(LoraPktPayloadB))
							log.Printf("webload json:%d  %s \n", len(JSONPktPayloadB), JSONPktPayloadB)

							//PULL_RESP packet load profile
							PullrespPayload := make([]byte, 4+len(JSONPktPayloadB))
							PullrespPayload[0] = 0x01
							PullrespPayload[3] = 0x03

							for i := 0; i < len(JSONPktPayloadB); i++ {
								PullrespPayload[4+i] = JSONPktPayloadB[i]
							}

							//DefaultGateway
							if DefaultGateway != "" {
								raddr41, err := net.ResolveUDPAddr("udp", DefaultGateway)
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr41)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr41, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							} else {
								raddr41, err := net.ResolveUDPAddr("udp", "172.16.10.56:7878")
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr41)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr41, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							}
						case 46: //get envent code
							var DataB = make([]byte, 18) //dmls_blue.pdf P.43,44
							DataB[0] = 0x02              // structure
							DataB[1] = 0x04              // 4 element
							DataB[2] = 0x06
							DataB[3] = 0x00
							DataB[4] = 0x00
							DataB[5] = 0x00
							DataB[6] = 0x01 // from entry =1
							DataB[7] = 0x06 // to entry =2 count =2
							DataB[8] = 0x00
							DataB[9] = 0x00
							DataB[10] = 0x00
							DataB[11] = 0x02
							DataB[12] = 0x12 // from selected value 1
							DataB[13] = 0x00
							DataB[14] = 0x02
							DataB[15] = 0x12
							DataB[16] = 0x00
							DataB[17] = 0x03 // to selected value-=5 or 0
							log.Printf("getevent:%X \n", ObsIndex)
							FrameCmd := []byte{0xc0, 0x09, 0x00, 0x2E, 0x02, 0x01, 0x46, 0x02}

							FrameCmd[1] = InvokeIDB[0]
							FrameB := packet.BuildFrame(DevAddrI, FrameCmd, DataB, DownlinkCnt)
							log.Printf("Webf    :%X \n", FrameB)
							//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
							//Pmsg = MHDR | MACPayload cmac = aes128_cmac(NwkSKey, B0 | Pmsg) page 45
							TxtGetPayload := packet.SecureFrame(FrameB, NwkSKeyB, AppSKeyB, DevAddrI, DownlinkCnt)
							log.Printf("Webplain:%X\n", TxtGetPayload)

							LoraPktPayloadB := packet.LoraPktEnc(TxtGetPayload, 0)
							log.Printf("Webframe:%X\n", LoraPktPayloadB)
							JSONPktPayloadB := []byte(packet.BuildTXJsonPkt(LoraPktPayloadB))
							log.Printf("webload json:%d  %s \n", len(JSONPktPayloadB), JSONPktPayloadB)

							//PULL_RESP packet load profile
							PullrespPayload := make([]byte, 4+len(JSONPktPayloadB))
							PullrespPayload[0] = 0x01
							PullrespPayload[3] = 0x03

							for i := 0; i < len(JSONPktPayloadB); i++ {
								PullrespPayload[4+i] = JSONPktPayloadB[i]
							}

							//DefaultGateway
							if DefaultGateway != "" {
								raddr46, err := net.ResolveUDPAddr("udp", DefaultGateway)
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr46)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr46, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							} else {
								raddr46, err := net.ResolveUDPAddr("udp", "172.16.10.56:7878")
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr46)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr46, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							}
						case 273: //load profile
							var DataB = make([]byte, 18) //dmls_blue.pdf P.43,44
							DataB[0] = 0x02              // structure
							DataB[1] = 0x04              // 4 element
							DataB[2] = 0x06              // from entry =1
							DataB[3] = 0x00
							DataB[4] = 0x00
							DataB[5] = 0x00
							DataB[6] = 0x01
							DataB[7] = 0x06 // to entry =2 count =2
							DataB[8] = 0x00
							DataB[9] = 0x00
							DataB[10] = 0x00
							DataB[11] = 0x02
							DataB[12] = 0x12 // from selected value 1
							DataB[13] = 0x00
							DataB[14] = 0x01
							DataB[15] = 0x12
							DataB[16] = 0x00
							DataB[17] = 0x00 // to selected value-=5 or 0
							FrameCmd := []byte{0xc0, 0x09, 0x01, 0x11, 0x02, 0x01, 0x46, 0x02}

							FrameCmd[1] = InvokeIDB[0]
							FrameB := packet.BuildFrame(DevAddrI, FrameCmd, DataB, DownlinkCnt)
							log.Printf("Webf    :%X \n", FrameB)
							//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
							//Pmsg = MHDR | MACPayload cmac = aes128_cmac(NwkSKey, B0 | Pmsg) page 45
							TxtGetPayload := packet.SecureFrame(FrameB, NwkSKeyB, AppSKeyB, DevAddrI, DownlinkCnt)
							log.Printf("Webplain:%X\n", TxtGetPayload)

							LoraPktPayloadB := packet.LoraPktEnc(TxtGetPayload, 0)
							log.Printf("Webframe:%X\n", LoraPktPayloadB)
							JSONPktPayloadB := []byte(packet.BuildTXJsonPkt(LoraPktPayloadB))
							log.Printf("webload json:%d  %s \n", len(JSONPktPayloadB), JSONPktPayloadB)

							//PULL_RESP packet load profile
							PullrespPayload := make([]byte, 4+len(JSONPktPayloadB))
							PullrespPayload[0] = 0x01
							PullrespPayload[3] = 0x03

							for i := 0; i < len(JSONPktPayloadB); i++ {
								PullrespPayload[4+i] = JSONPktPayloadB[i]
							}

							//DefaultGateway
							if DefaultGateway != "" {
								raddr273, err := net.ResolveUDPAddr("udp", DefaultGateway)
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr273)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr273, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							} else {
								raddr273, err := net.ResolveUDPAddr("udp", "172.16.10.56:7878")
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr273)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr273, PullrespPayload[:4], string(JSONPktPayloadB))
								exception.LogFile(s)
							}
							//fmt.Printf("GW273:%s\n", DefaultGateway)
							//time.Sleep(1 * time.Second)

						}
					case 0xC1:
						switch ObsIndex {
						case 0: //set FAN

						case 41: //set current datetime
							PullrespPayload := packet.SetDateTime(DevAddrI, NwkSKeyB, AppSKeyB, DownlinkCnt, InvokeIDB)

							// fmt.Printf("-----------Webplain:%X\n", TxtGetPayload)
							// LoraPktPayloadB := packet.LoraPktEnc(TxtGetPayload, 0)
							// fmt.Printf("-----------Webframe:%X\n", LoraPktPayloadB)
							// JSONPktPayloadB := []byte(packet.BuildTXJsonPkt(LoraPktPayloadB))
							// fmt.Printf("-----------webload json:%d  %s \n", len(JSONPktPayloadB), JSONPktPayloadB)

							//PULL_RESP packet set datetime
							// PullrespPayload := make([]byte, 4+len(JSONPktPayloadB))
							// PullrespPayload[0] = 0x01
							// PullrespPayload[3] = 0x03
							//
							// for i := 0; i < len(JSONPktPayloadB); i++ {
							// 	PullrespPayload[4+i] = JSONPktPayloadB[i]
							// }

							//DefaultGateway
							if DefaultGateway != "" {
								raddr41, err := net.ResolveUDPAddr("udp", DefaultGateway)
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr41)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr41, PullrespPayload[:4], string(PullrespPayload[4:]))
								exception.LogFile(s)
							} else {
								raddr41, err := net.ResolveUDPAddr("udp", "192.168.43.160:7878")
								exception.CheckError(err)
								_, err = ServerConn.WriteToUDP(PullrespPayload, raddr41)
								exception.CheckError(err)
								s := fmt.Sprintf("TX pull rsp to %s :%X %s \n", raddr41, PullrespPayload[:4], string(PullrespPayload[4:]))
								exception.LogFile(s)
							}
						case 46: //set envent code

						case 273: //set profile

						}
					case 0xC3: // rejoin
						log.Printf("rejoin!")
					}
				} else {
					log.Print("Invalid parameter: Dev", DevAddrI)
				}
			} else {
				log.Printf("Web command format error!")

			}

		case 24: //Pull data check
			if buf[0] == 0x01 && buf[3] == 0x02 {

				//fill pull ack
				pullack[1] = buf[1]
				pullack[2] = buf[2]

				//fmt.Printf("Pull ack: %d %X\n", n, pullack)
				//_, err = ServerConn.WriteToUDP([]byte(pullack[:4]), remoteAddr)
				DefaultGateway = remoteAddr.String()
				//q := net.ParseIP("192.168.0.1")
				raddr, err := net.ResolveUDPAddr("udp", DefaultGateway)
				exception.CheckError(err)
				//fmt.Printf("GW:%s\n", DefaultGateway)
				if err != nil {
					log.Printf(err.Error())

				}
				_, err = ServerConn.WriteToUDP([]byte(pullack[:4]), raddr)
				exception.CheckError(err)
				//s := fmt.Sprintf("TX pull ack to %s :%X \n", remoteAddr, pullack[:4])
				exception.LogFile(fmt.Sprintf("TX pull ack to %s :%X \n", remoteAddr, pullack[:4]))
			} else {
				log.Printf("Pull data format error!")

			}

		default: // base64 json
			//PUSH_DATA
			if buf[0] == 0x01 && buf[3] == 0x00 {
				//fill push ack
				pushack[1] = buf[1]
				pushack[2] = buf[2]

				//trim size and data s := buf[21:n-2]
				//fmt.Printf("data: %X\n",get_jsondata(buf[21:n-2]))

				LoraPkt := packet.GetJsonData(buf[21 : n-2])

				//extract lora information
				packet.LoraPktDec(LoraPkt)

				//extract payload content
				JoinJSONPktB, KeyJSONPktB, SetDateTimeJSONPktB := packet.PayloadPktDec(LoraPkt)

				//response
				//fmt.Printf("Pull ack: %d %X\n", n, pullack)
				_, err = ServerConn.WriteToUDP([]byte(pushack[:4]), remoteAddr)
				exception.CheckError(err)
				//s := fmt.Sprintf("TX push ack to %s :%X\n", remoteAddr.String(), pushack[:4])
				exception.LogFile(fmt.Sprintf("TX push ack to %s :%X\n", remoteAddr.String(), pushack[:4]))

				//PULL_RESP packet Join
				pullresp1 := make([]byte, 4+len(JoinJSONPktB))
				pullresp1[0] = 0x01
				pullresp1[3] = 0x03
				for i := 0; i < len(JoinJSONPktB); i++ {
					pullresp1[4+i] = JoinJSONPktB[i]
				}

				//PULL_RESP packet set guk ak
				pullresp2 := make([]byte, 4+len(KeyJSONPktB))
				pullresp2[0] = 0x01
				pullresp2[3] = 0x03
				for i := 0; i < len(KeyJSONPktB); i++ {
					pullresp2[4+i] = KeyJSONPktB[i]
				}

				//PULL_RESP packet set datetime
				pullresp3 := make([]byte, 4+len(SetDateTimeJSONPktB))
				pullresp3[0] = 0x01
				pullresp3[3] = 0x03
				for i := 0; i < len(SetDateTimeJSONPktB); i++ {
					pullresp3[4+i] = SetDateTimeJSONPktB[i]
				}
				//_, err = ServerConn.WriteToUDP(pullresp3, remoteAddr)

				//DefaultGateway
				if DefaultGateway != "" {
					raddr0, err := net.ResolveUDPAddr("udp", DefaultGateway)
					exception.CheckError(err)
					//join accept: delay 6 sec for rx2
					time.Sleep(6 * time.Second)
					_, err = ServerConn.WriteToUDP(pullresp1, raddr0)
					exception.CheckError(err)
					//s1 := fmt.Sprintf("TX join accept to %s :%X\n", raddr0.String(), pullresp1)
					exception.LogFile(fmt.Sprintf("TX join accept to %s :%X\n", raddr0.String(), pullresp1))

					//send key: delay 5 sec for fan update
					time.Sleep(5 * time.Second)
					_, err = ServerConn.WriteToUDP(pullresp2, raddr0)
					exception.CheckError(err)
					//s2 := fmt.Sprintf("TX join set key to %s :%X\n", raddr0.String(), pullresp2)
					exception.LogFile(fmt.Sprintf("TX join set key to %s :%X\n", raddr0.String(), pullresp2))

					//setdatetime: delay 9 sec for fan prepare key
					time.Sleep(9 * time.Second)
					_, err = ServerConn.WriteToUDP(pullresp3, raddr0)
					exception.CheckError(err)
					//s3 := fmt.Sprintf("TX join setdatetime to %s :%X\n", raddr0.String(), pullresp3)
					exception.LogFile(fmt.Sprintf("TX join setdatetime to %s :%X\n", raddr0.String(), pullresp3))

				} else {
					raddr0, err := net.ResolveUDPAddr("udp", "192.168.43.160:7878")
					exception.CheckError(err)
					_, err = ServerConn.WriteToUDP(pullresp1, raddr0)
					exception.CheckError(err)
					time.Sleep(5 * time.Second)

					_, err = ServerConn.WriteToUDP(pullresp2, raddr0)
					exception.CheckError(err)
					time.Sleep(3 * time.Second)

					_, err = ServerConn.WriteToUDP(pullresp3, raddr0)
					exception.CheckError(err)
					time.Sleep(2 * time.Second)
				}

			} else {

				fmt.Printf("Pull data format error!")

			}

		}

	}

}
