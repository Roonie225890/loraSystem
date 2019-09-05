/* main function            */
/* file name: packet.go     */
/* 		link: packetAct.go */
/* 		                    */
/*    update: 20181122      */

package packet

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"loranet20181205/database"
	"loranet20181205/encrypt"
	"loranet20181205/exception"
	"math/rand"
	"time"
)

type pkt struct {
	Size int
	Data string
}

//uplink frequency center
var uplinkFC = [12]int{849357110, 849499970, 849642830, 849785690,
	849928550, 850071410, 850214270, 850357130,
	850499990, 850642850, 850785710, 850928570}

var downlinkFC = [8]int{839142860, 839428580, 839714300, 840000020,
	840285740, 840571460, 841142900, 840857180}

// GetJsonData transform JSON data to bytes
func GetJsonData(buff []byte) []byte {

	pkt1 := pkt{}
	//fmt.Printf("ldata:%s \n", string(buff))

	err := json.Unmarshal(buff, &pkt1)
	exception.CheckError(err)

	//fmt.Printf("pdata:%s \n", pkt1.Data)
	//fmt.Printf("length:%d data:%s \n", pkt1.Size, pkt1.Data)

	decoded, err := base64.StdEncoding.DecodeString(pkt1.Data)
	exception.CheckError(err)
	if len(decoded) != pkt1.Size {

		s := fmt.Sprintf("length:%d data:%s \n", pkt1.Size, pkt1.Data)
		log.Printf(s)        //to terminal
		exception.LogFile(s) //log to file
		//fmt.Printf("length:%d data:%s \n", pkt1.Size, pkt1.Data)

		return decoded
	} else {

		return decoded
	}
}

func BuildTXJsonPkt(buff []byte) string {
	encodedB64 := base64.StdEncoding.EncodeToString(buff)
	JsonPkt := fmt.Sprintf("{\"txpk\":[{\"size\":%d,\"data\":\"%s\"}]}", len(buff), encodedB64)

	return JsonPkt
}

//struct RX_PacketBurst {
//    uint32T    freq_hz;        // Central frequency in Hz, Sending sequency as: freq_hz[7:0], freq_hz[15:8], freq_hz[23:16], freq_hz[31:24]
//    uint32T    timestamp;      // Corrected RX timestamp in 1 uS resolution, Sending Sending sequency as: TS[7:0], TS[15:8], TS[23:16], TS[31:24]
//    uint8T     sf;             // RX datarate of the packet (SF for LoRa) like 7,8,...12
//    uint8T     coderate;       // 1: cr 4/5, 2: cr 4/6, 3: cr 4/7, 4:cr 4/8
//    uint8T     chain;          // bit(1:0) as RF chain(0,1,2,3),  bit(7:2) as IF chain(0,1,2,..19)
//    uint8T     status;         // status=1: Packet without CRC, status=5: Packet with CRC and CRC check ok.
//    uint8T     modulation;     // modulation=0x10 for Lora, all other value reserved
//    uint8T     bandwidth;      // modulation bandwidth, bandwidth=0x01 as 500KHZ, bandwidth=0x02 as 250KHZ, bandwidth=0x03 as 125KHZ
//    int16T         rssi;           // average packet RSSI in dB, format of Q13.2, resolution of 0.25 dB
//    int16T         snr;            // average packet SNR, in dB (LoRa only), format of Q13.2, resolution of 0.25 dB
//    int16T         snr_min;        // minimum packet SNR, in dB (LoRa only), format of Q13.2, resolution of 0.25 dB
//    int16T         snr_max;        // maximum packet SNR, in dB (LoRa only), format of Q13.2, resolution of 0.25 dB
//    uint16T    size;           // payload size in bytes
//    uint8T     payload[256];   // buffer containing the payload length of size  ( payload[size] )
//    };

func LoraPktDec(loraPkt []byte) {
	freq := uint32(loraPkt[0]) | uint32(loraPkt[1])<<8 | uint32(loraPkt[2])<<16 | uint32(loraPkt[3])<<24
	//timestamp := uint32(loraPkt[4]) | uint32(loraPkt[5])<<8 | uint32(loraPkt[6])<<16 | uint32(loraPkt[7])<<24
	sf := uint32(loraPkt[8])
	//coderate := uint32(loraPkt[9])
	//chain := uint32(loraPkt[10])
	//status := uint32(loraPkt[11])
	//modulation := uint32(loraPkt[12])
	//bandwidth := uint32(loraPkt[13]
	rssi := 0.25 * float64(uint32(loraPkt[14])|uint32(loraPkt[15])<<8)
	snr := 0.25 * float64(uint32(loraPkt[16])|uint32(loraPkt[17])<<8)
	//snr_min := float64((loraPkt[18]) | ((loraPkt[19])<<8)) * 0.25
	//snr_max := float64((loraPkt[20]) | ((loraPkt[21])<<8)) * 0.25
	size := uint32(loraPkt[22]) | uint32(loraPkt[23])<<8
	exception.Info.Print("freq:", freq, " sf:", sf, " rssi:", rssi, " snr:", snr, " size:", size)

	s := fmt.Sprintf("freq:%d  sf:%d rssi:%.2f snr:%.2f size:%d\n", freq, sf, rssi, snr, size)
	log.Printf(s)        //to terminal
	exception.LogFile(s) //log to file

	//fmt.Printf("freq: %d \n", freq)
	//fmt.Printf("sf: %d \n", sf)
	//fmt.Printf("rssi: %.2f \n", rssi)
	//fmt.Printf("snr: %.2f \n", snr)
	////fmt.Printf("snr_min: %.2f \n", snr_min)
	////fmt.Printf("snr_max: %.2f \n", snr_max)
	//fmt.Printf("size: %d \n", size)
}

//struct TX_PacketS {
//    uint16T   length;         // Number of bytes followed, !! 2 bytes this length is excluded, sending as: length[7:0], length[15:8]
//    uint32T   freq_hz;        // Center frequency of TX in Hz, sending as: freq_hz[7:0], freq_hz[15:8], freq_hz[23:16], freq_hz[31:24]
//    uint32T   timestamp;      // TX timestamp in 1 uS resolution (for tx_mode=1), sending as: TS[7:0], TS[15:8], TS[23:16], TS[31:24]
//    uint8T    tx_mode;        // tx_mode=0: TX service based on FIFO, tx_mode=1: TX service based on timestamp
//    int8T     rf_power;       // RF output power in dBm, Ex. rf_power=36 means +36dBm(4W) output
//    uint8T    sf;             // TX datarate of the packet (SF for LoRa) like 7,8,...12
//    uint8T    coderate;       // 1: cr 4/5, 2: cr 4/6, 3: cr 4/7, 4:cr 4/8
//    uint8T    CRC;            // CRC=1: sending with auto_generated CRC, CRC=0: sending without CRC (default)
//    uint8T    modulation;     // modulation=0x10 for Lora, all other value reserved
//    uint8T    bandwidth;      // modulation bandwidth, bandwidth=0x01 as 500KHZ, bandwidth=0x02 as 250KHZ, bandwidth=0x03 as 125KHZ
//    uint8T    preamble;       // Number of preamble count. (fixed to 8 in our application)
//    uint8T    invertIQ;      // invertIQ.0=0: normal(default), invertIQ.0=1: I and Q signals are inverted, invertIQ[7:1] are reserved
//    uint8T    no_header;      // no_header.0=0: normal header(default) no_header.0=1: Enable implicit header mode, no_header[7:1] are reserved
//    uint16T   size;           // payload size in bytes
//    uint8T    payload[256];   // buffer containing the payload length of size ( payload[size] )
//    };
func LoraPktEnc(loraPkt []byte, addrB int) []byte {
	var payloadLength int
	var pktLength int
	var loraPktB []byte
	payloadLength = len(loraPkt)
	pktLength = payloadLength + 22

	var pktLengthB = make([]byte, 4)
	binary.LittleEndian.PutUint32(pktLengthB, uint32(pktLength-2)) //the 2 lenght byte is not included

	var payloadLengthB = make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLengthB, uint32(payloadLength))

	loraPktB = make([]byte, 22+payloadLength)

	//downlinkFC[0]  can be random select
	var downlinkFCB = make([]byte, 4)
	binary.LittleEndian.PutUint32(downlinkFCB, uint32(downlinkFC[4]))

	var timestampB = make([]byte, 4)
	binary.LittleEndian.PutUint32(timestampB, uint32(time.Now().UnixNano()/1000))

	//fmt.Printf("PKT:%X\n", loraPkt)
	var DataRateSFB = make([]byte, 4)
	var devAddrI int
	if addrB == 0 {
		devAddrI = int(loraPkt[1]) | int(loraPkt[2])<<8 | int(loraPkt[3])<<16 | int(loraPkt[4])<<24
		// devAddrI = 23554699
		binary.LittleEndian.PutUint32(DataRateSFB, uint32(12-database.DbGetDataRate(devAddrI))) //default SF == 12
	} else {
		binary.LittleEndian.PutUint32(DataRateSFB, uint32(12)) //default SF == 12
	}

	loraPktB[0] = pktLengthB[0]
	loraPktB[1] = pktLengthB[1]
	loraPktB[2] = downlinkFCB[0]
	loraPktB[3] = downlinkFCB[1]
	loraPktB[4] = downlinkFCB[2]
	loraPktB[5] = downlinkFCB[3]
	loraPktB[6] = 0x40  //timestampB[0]
	loraPktB[7] = 0x42  //timestampB[1]
	loraPktB[8] = 0x0F  //timestampB[2]
	loraPktB[9] = 0x00  //timestampB[3]
	loraPktB[10] = 0x00 //tx_mode=0: TX service based on FIFO
	loraPktB[11] = 0x24 //rf_power=36

	loraPktB[12] = DataRateSFB[0]     //SF for LoRa
	loraPktB[13] = 0x01               //1: cr 4/5
	loraPktB[14] = 0x00               //CRC=0: sending without CRC (default)
	loraPktB[15] = 0x10               //modulation=0x10 for Lora
	loraPktB[16] = 0x02               //bandwidth=0x02 as 250KHZ
	loraPktB[17] = 0x08               //Number of preamble count. (fixed to 8
	loraPktB[18] = 0x00               //invertIQ.0=0: normal(default)
	loraPktB[19] = 0x00               //no_header.0=0: normal header(default)
	loraPktB[20] = payloadLengthB[0] //payload size L
	loraPktB[21] = payloadLengthB[1] //payload size H
	for i := 0; i < payloadLength; i++ {
		loraPktB[22+i] = loraPkt[i]
	}
	return loraPktB
}

func PayloadPktDec(loraPkt []byte) ([]byte, []byte, []byte) {
	MH := loraPkt[24]
	switch MH {
	case 0x00:
		//join request
		joinECB, set_keys, setDatetime := JoinRequest(loraPkt[24:])

		//fmt.Printf("AK:%s GUK:%s joinECB:%X\n", AKS, GUKS, joinECB)

		loraPktB := LoraPktEnc(joinECB, 1) //no address
		fmt.Printf("join Lora:%X\n", loraPktB)
		json_pkt := BuildTXJsonPkt(loraPktB)
		fmt.Printf("join json:%s\n", json_pkt)

		set_keys_loraPktB := LoraPktEnc(set_keys, 0)
		fmt.Printf("key Lora:%X\n", set_keys_loraPktB)
		set_keysJson_pkt := BuildTXJsonPkt(set_keys_loraPktB)
		fmt.Printf("key json:%s\n", set_keysJson_pkt)

		setDatetime_loraPktB := LoraPktEnc(setDatetime, 0)
		fmt.Printf("datetime Lora:%X\n", setDatetime_loraPktB)
		setDatetimeJson_pkt := BuildTXJsonPkt(setDatetime_loraPktB)
		fmt.Printf("datetime json:%s\n", setDatetimeJson_pkt)

		return []byte(json_pkt), []byte(set_keysJson_pkt), []byte(setDatetimeJson_pkt)

	case 0x10:
		return loraPkt, loraPkt, loraPkt
	default:
		fmt.Printf("packet payload error!\n")
		return loraPkt, loraPkt, loraPkt

	}
}

func BuildFrame(devAddr int, cmd []byte, data []byte, downlinkCnt int) []byte {
	//C0 ->mhdr
	//B3 00 00 02 ->devAddr
	//23 ->this.fctrl
	//0B 01 ->RequestProc.downlinkCnt.get(devAddr)
	//01 ->fport
	//C1 ->Coding   : below encoded
	//09 ->InvokeIdPriority
	//00 ->ObisIdx>>8
	//00 ->ObisIdx
	//02 ->AttriMethd
	//01 ->Flags
	//46 ->CntOrBlk
	//00 ->SelectiveAccess

	frameData := make([]byte, 17+len(data)) //9+8+lenght of data

	var devAddrB = make([]byte, 4)
	binary.LittleEndian.PutUint32(devAddrB, uint32(devAddr))

	frameData[0] = 0xC0        //polling request
	frameData[1] = devAddrB[0] //Low byte
	frameData[2] = devAddrB[1]
	frameData[3] = devAddrB[2]
	frameData[4] = devAddrB[3] //high byte
	frameData[5] = 0x21        //Pending bit 7, fctl bits 0~6

	var downlinkCntB = make([]byte, 4)
	binary.LittleEndian.PutUint32(downlinkCntB, uint32(downlinkCnt))

	frameData[6] = downlinkCntB[0]
	frameData[7] = downlinkCntB[1]
	frameData[8] = 0x01 // fport

	for i := 0; i < 8; i++ {
		frameData[9+i] = cmd[i]
	}

	for i := 0; i < len(data); i++ {
		frameData[17+i] = data[i]
	}

	return frameData
}

func SecureFrame(data []byte, NwkSKey []byte, AppSKey []byte, devAddr int, downlinkCnt int) []byte {
	//Polling Messages encryption

	frameDataCmac := make([]byte, 20+len(data)) // 16+len(data)+4
	blockNum := len(data)/len(NwkSKey) + 1
	frameSec := make([]byte, (blockNum * (len(NwkSKey))))
	var downlinkCntB = make([]byte, 4)
	binary.LittleEndian.PutUint32(downlinkCntB, uint32(downlinkCnt))
	var devAddrB = make([]byte, 4)
	binary.LittleEndian.PutUint32(devAddrB, uint32(devAddr))
	var counterB = make([]byte, 4)
	for i := 0; i < blockNum; i++ {
		frameSec[i*16] = 0x01
		frameSec[i*16+1] = 0x00
		frameSec[i*16+2] = 0x00
		frameSec[i*16+3] = 0x00
		frameSec[i*16+4] = 0x00
		frameSec[i*16+5] = 0x01 //Dir=0x00 for uplink, Dir=0x01 for downlink
		frameSec[i*16+6] = devAddrB[0]
		frameSec[i*16+7] = devAddrB[1]
		frameSec[i*16+8] = devAddrB[2]
		frameSec[i*16+9] = devAddrB[3]
		frameSec[i*16+10] = 0x00 //Pending bit
		frameSec[i*16+11] = 0x21 //fctl bits 0~6,Unum channel 33
		frameSec[i*16+12] = downlinkCntB[0]
		frameSec[i*16+13] = downlinkCntB[1]
		frameSec[i*16+14] = 0x00
		binary.LittleEndian.PutUint32(counterB, uint32(i+1))
		frameSec[i*16+15] = counterB[0]
	}

	frameSecAes := aesacc.GetAesECB(frameSec, AppSKey)

	for i := 9; i < len(data); i++ {
		data[i] = data[i] ^ frameSecAes[i-9]
	}
	//fmt.Printf("sec code:%X\n", frameSec)
	//fmt.Printf("enc code:%X\n", frameSecAes)
	//fmt.Printf("framS:%X\n", data)
	//for i:=9; i<len(data); i++{
	//	data[i] = data[i] ^ frameSecAes[i-9]
	//}
	//fmt.Printf("framR:%X\n",data)

	//Polling Messages MIC
	for i := 0; i < 16; i++ {
		frameDataCmac[i] = frameSec[i]
	}
	frameDataCmac[0] = 0x49
	frameDataCmac[5] = 0x01
	binary.LittleEndian.PutUint32(counterB, uint32(len(data)))
	frameDataCmac[15] = counterB[0]

	for i := 0; i < len(data); i++ {
		frameDataCmac[16+i] = data[i]
	}
	//MIC := getCmac(frameDataCmac[0:len(data)+17], NwkSKey)
	MIC := aesacc.GetCmac(frameDataCmac[0:len(data)+16], AppSKey)

	//fmt.Printf("framM:%X\n", MIC)

	for i := 0; i < 4; i++ {
		frameDataCmac[16+len(data)+i] = MIC[i]
	}

	//debug MIC

	// testmic := []byte{0xC0, 0xB3, 0x00, 0x00, 0x02, 0x21, 0x8F, 0x00, 0x01, 0x88, 0xF6, 0x1B, 0x66, 0x90, 0x9A, 0x20, 0x6D, 0x11, 0x15, 0x55, 0xF7, 0x8A, 0xE1, 0x30, 0x10, 0x05, 0x44, 0x99, 0x2D, 0x38, 0x04, 0x1A, 0x52, 0xAA, 0xCD}
	// fmt.Printf("Length frameDataCmac:%d test:%d\n", len(frameDataCmac), len(testmic))
	// for i := 0; i < len(testmic); i++ {
	// 	frameDataCmac[16+i] = testmic[i]
	// }
	// MIC1 := getCmac(frameDataCmac[0:len(testmic)+16], AppSKey)
	// fmt.Printf("testMIC:%X\n", MIC1)

	return frameDataCmac[16:]
}

func JoinRequest(loraPkt []byte) ([]byte, []byte, []byte) {

	MeterID := string(loraPkt[1:9]) // interesting 1:9 not 1:8
	DevEUI := fmt.Sprintf("%X", loraPkt[12:17])
	DevNonce := loraPkt[17:19]

	// var syscnt []database.TbSyscnt
	//var NwkSKey string
	//var AppSKey string

	// Create Connection
	db, err := database.DbConnect()
	exception.CheckError(err)

	// Close Connection
	defer db.Close()

	// Query
	meterList := database.GetMeterList(MeterID)
	// qstr1 := "SELECT meterUniqueID, meterID, customerId, AK, GUK FROM tb_meter_list WHERE meterID = '"
	// qstr1 += MeterID + "'"
	// fmt.Printf("qstr1: %s \n", qstr1)

	log.Print("AK:", meterList.AK, " GUK:", meterList.GUK)

	// Query
	//	qstr2 := "SELECT MeterID, DevEUI, AppKey, NwkSKey, AppSKey FROM tb_edInfo WHERE DevEUI = '"
	// qstr2 := "SELECT id, DevEUI, AppKey FROM tb_edInfo WHERE DevEUI = '"
	// qstr2 += DevEUI + "'"
	// fmt.Printf("qstr2: %s \n", qstr2)

	// rows2, err := db.Query(qstr2)
	// exception.CheckError(err)
	//
	// for rows2.Next() {
	// 	err = rows2.Scan(&id, &DevEUI, &AppKey)
	// 	exception.CheckError(err)
	//
	edInfo := database.GetEdInfo(DevEUI)
	log.Print("DevEUI:", edInfo.DevEUI, " AppKey:", edInfo.AppKey)
	//
	// }
	log.Print("DevNonce:", DevNonce)
	log.Print("LoraPkt:", loraPkt)

	//cmac check
	txt, _ := hex.DecodeString(edInfo.AppKey)
	cmac := aesacc.GetCmac(loraPkt[0:19], txt)

	//fmt.Printf("cmac: %X\n", cmac)
	if loraPkt[19] == cmac[0] && loraPkt[20] == cmac[1] && loraPkt[21] == cmac[2] && loraPkt[22] == cmac[3] {
		fmt.Printf("check sum success! cmac: %X\n", cmac)
	} else {
		fmt.Printf("check sum fail! cmac: %X\n", cmac)
	}

	//set key generation string
	netID := 21 //default netID 21
	AppNonce := rand.Intn(65535)
	temp := fmt.Sprintf("%02X%06X%06X%04X%014X", 1, AppNonce, netID, DevNonce, 0)
	nkTxt, _ := hex.DecodeString(temp)

	fmt.Printf("key string1 :%s\n", temp)

	//network session key
	NwkSKey := aesacc.GetAesECB(nkTxt, txt)
	fmt.Printf("NwkSKey:%X\n", NwkSKey)
	NwkSKeyB64 := base64.StdEncoding.EncodeToString(NwkSKey[0:16])
	fmt.Printf("NwkSKeyb:%s\n", NwkSKeyB64)

	//app session key
	nkTxt[0] = 0x02
	AppSKey := aesacc.GetAesECB(nkTxt, txt)
	fmt.Printf("AppSKey:%X\n", AppSKey)
	AppSKeyB64 := base64.StdEncoding.EncodeToString(AppSKey[0:16])
	fmt.Printf("AppSKeyb:%s\n", AppSKeyB64)

	DevNonce32 := (int(DevNonce[0]) << 8) + (int(DevNonce[1]))

	num := database.DbUpdate(edInfo.ID, edInfo.MeterID, NwkSKeyB64, AppSKeyB64, AppNonce, DevNonce32, DevEUI, edInfo.MeterUniqueID, edInfo.AK, edInfo.GUK)

	if num == 1 {
		fmt.Printf("updated Rows:%d\n", num)
	}

	//join accept packet
	//cmac = aes128Cmac(AppKey, MHDR | AppNonce | NetID | DevAddr |LinkParam)

	DevEUIB, err := hex.DecodeString(DevEUI)
	exception.CheckError(err)
	var devAddrI = int(DevEUIB[4]) | int(DevEUIB[3])<<8 | int(DevEUIB[2])<<16 | int(DevEUIB[0])<<24

	log.Print("DevEUI:", DevEUI)
	log.Print("Dev:", devAddrI)

	// Query
	//	qstr3 := "SELECT linkparam FROM tb_edConf WHERE devAddr = '"
	// qstr3 := "SELECT linkparam, devAddr FROM tb_edConf WHERE devAddr = " + strconv.Itoa(devAddrI)
	// //qstr3 += devAddr
	// fmt.Printf("qstr3: %s \n", qstr3)
	//
	// rows3, err := db.Query(qstr3)
	// exception.CheckError(err)
	//
	// for rows3.Next() {
	// 	err = rows3.Scan(&linkparam, &devAddr)
	// 	exception.CheckError(err)
	edConf := database.GetEdConf(devAddrI)
	log.Print("Linkparam:", edConf.Linkparam, " DevAddr:", edConf.DevAddr)
	// }
	log.Print("AppNonce:", AppNonce)

	var AppNonceB = make([]byte, 4)
	binary.LittleEndian.PutUint32(AppNonceB, uint32(AppNonce))

	var linkparamB = make([]byte, 4)
	binary.LittleEndian.PutUint32(linkparamB, uint32(edConf.Linkparam))

	TxtJoinAccept := make([]byte, 17)
	TxtJoinAccept[0] = 0x20
	TxtJoinAccept[1] = AppNonceB[2] //(AppNonce >>16)&0xff
	TxtJoinAccept[2] = AppNonceB[1] // (AppNonce >>8)&0xff
	TxtJoinAccept[3] = AppNonceB[0] // (AppNonce )&0xff
	TxtJoinAccept[4] = 0x00
	TxtJoinAccept[5] = 0x15
	TxtJoinAccept[6] = DevEUIB[4] //hiByte
	TxtJoinAccept[7] = DevEUIB[3]
	TxtJoinAccept[8] = DevEUIB[2]
	TxtJoinAccept[9] = DevEUIB[0]     //lowByte
	TxtJoinAccept[10] = linkparamB[2] //(linkparam >>16)&0xff
	TxtJoinAccept[11] = linkparamB[1] //(linkparam >>8)&0xff
	TxtJoinAccept[12] = linkparamB[0] //(linkparam )&0xff

	//generate cmac
	cmacJoinAccept := aesacc.GetCmac(TxtJoinAccept[0:13], txt)
	TxtJoinAccept[13] = cmacJoinAccept[0]
	TxtJoinAccept[14] = cmacJoinAccept[1]
	TxtJoinAccept[15] = cmacJoinAccept[2]
	TxtJoinAccept[16] = cmacJoinAccept[3]

	fmt.Printf("join accept: %x \n", TxtJoinAccept)

	//ECB encryption
	TxtJoinAcceptDec := aesacc.GetAesECBDec(TxtJoinAccept[1:17], txt)
	for i := 0; i < 16; i++ {
		TxtJoinAccept[1+i] = TxtJoinAcceptDec[i]
	}
	fmt.Printf("join acceptECB:%x \n", TxtJoinAccept)

	//build frame payload
	// Query
	// qstr4 := "SELECT downlinkCnt, InvokeId FROM tbSyscnt WHERE devAddr = " + strconv.Itoa(devAddrI)
	// fmt.Printf("qstr4: %s \n", qstr4)
	// rows4, err := db.Query(qstr4)
	// exception.CheckError(err)
	//
	// var rowCount int
	// for rows4.Next() {
	// 	err = rows4.Scan(&downlinkCnt, &InvokeId)
	// 	exception.CheckError(err)
	// 	fmt.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, InvokeId)
	// 	downlinkCnt++
	// 	InvokeId++
	// 	rowCount++
	// }
	//
	// //check counter table
	// if rowCount == 0 {
	// 	fmt.Printf("rowCount %d, Dnlnk: %d Ivk:%d\n", rowCount, downlinkCnt, InvokeId)
	// 	stmt, err := db.Prepare("insert into tbSyscnt(devAddr, downlinkCnt, uplinkCnt, InvokeId) values(?,?,?,?);")
	// 	result5, err := stmt.Exec(devAddrI, 1, 0, 0)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	defer stmt.Close()
	// 	lastInsertId, err := result5.LastInsertId()
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("insertid:%d\n", lastInsertId)
	// } else {
	// 	stmt, err := db.Prepare("update tbSyscnt set downlinkCnt=?, uplinkCnt=?, InvokeId=? where devAddr=?;")
	// 	result, err := stmt.Exec(downlinkCnt, uplinkCnt, InvokeId, devAddrI)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	defer stmt.Close() // Close the statement when we leave main() / the program terminates
	// 	rowsAffect, err := result.RowsAffected()
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("update:%d\n", rowsAffect)
	// }
	syscn := database.GetSyscn(devAddrI)
	var downlinkCounter = syscn.DownlinkCnt - 1
	//set GUK AK
	//0x09 string, 0x20 (32)length
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var CntBlk1 = make([]byte, 4)
	binary.LittleEndian.PutUint32(CntBlk1, r1.Uint32())
	gukakB, _ := hex.DecodeString("0920" + edInfo.GUK + edInfo.AK)
	frameCmd := []byte{0xc1, 0x09, 0x00, 0x00, 0x02, 0x01, CntBlk1[0], 0x00}

	var InvokeIDB = make([]byte, 4)
	binary.LittleEndian.PutUint32(InvokeIDB, uint32(syscn.InvokeID-1))
	frameCmd[1] = InvokeIDB[0]

	frameB := BuildFrame(devAddrI, frameCmd, gukakB, downlinkCounter)

	fmt.Printf("frame:%X \n", frameB)
	//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
	//Pmsg = MHDR | MACPayload cmac = aes128Cmac(NwkSKey, B0 | Pmsg) page 45
	TxtSetGukAk := SecureFrame(frameB, NwkSKey[0:16], AppSKey[0:16], devAddrI, downlinkCounter)

	fmt.Printf("fram :%X\n", TxtSetGukAk)

	//set date time
	//0x19 datetime
	setdatetimeB, _ := hex.DecodeString("1907E2060C0209112FFF80")
	tSet := time.Now()
	var tTempB = make([]byte, 4)
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Year()))
	setdatetimeB[1] = tTempB[0]
	setdatetimeB[2] = tTempB[1]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Month()))
	setdatetimeB[3] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Day()))
	setdatetimeB[4] = tTempB[0]
	setdatetimeB[5] = 0xff

	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Hour()))
	setdatetimeB[6] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Minute()))
	setdatetimeB[7] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Second()))
	setdatetimeB[8] = tTempB[0]

	r2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var CntBlk2 = make([]byte, 4)
	binary.LittleEndian.PutUint32(CntBlk2, r2.Uint32())
	setFrameCmd := []byte{0xc1, 0x09, 0x00, 0x00, 0x02, 0x01, CntBlk2[0], 0x00}

	//var InvokeIDB = make([]byte, 4)
	//binary.LittleEndian.PutUint32(InvokeIDB, uint32(InvokeId))
	setFrameCmd[1] = InvokeIDB[0]

	setFrameB := BuildFrame(devAddrI, setFrameCmd, setdatetimeB, downlinkCounter)

	fmt.Printf("frame:%X \n", setFrameB)

	//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
	//Pmsg = MHDR | MACPayload cmac = aes128Cmac(NwkSKey, B0 | Pmsg) page 45
	TxtSetDatetime := SecureFrame(setFrameB, NwkSKey[0:16], AppSKey[0:16], devAddrI, downlinkCounter)

	return TxtJoinAccept, TxtSetGukAk, TxtSetDatetime

}
