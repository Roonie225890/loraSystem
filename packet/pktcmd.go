package packet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func SetDateTime(DevAddrI int, NwkSKeyB []byte, AppSKeyB []byte, DownlinkCounter int, InvokeIDB []byte) []byte {

	//set date time
	//0x19 datetime
	SetdatetimeB, _ := hex.DecodeString("1907E2060C0209112FFF80")
	tSet := time.Now()
	var tTempB = make([]byte, 4)
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Year()))
	SetdatetimeB[1] = tTempB[0]
	SetdatetimeB[2] = tTempB[1]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Month()))
	SetdatetimeB[3] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Day()))
	SetdatetimeB[4] = tTempB[0]
	SetdatetimeB[5] = 0xff

	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Hour()))
	SetdatetimeB[6] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Minute()))
	SetdatetimeB[7] = tTempB[0]
	binary.LittleEndian.PutUint32(tTempB, uint32(tSet.Second()))
	SetdatetimeB[8] = tTempB[0]

	r2 := rand.New(rand.NewSource(time.Now().UnixNano()))
	var CntBlk2 = make([]byte, 4)
	binary.LittleEndian.PutUint32(CntBlk2, r2.Uint32())
	SetFrameCmd := []byte{0xc1, 0x09, 0x00, 0x00, 0x02, 0x01, CntBlk2[0], 0x00}

	//var InvokeIdB = make([]byte, 4)
	//binary.LittleEndian.PutUint32(InvokeIdB, uint32(InvokeId))
	SetFrameCmd[1] = InvokeIDB[0]

	SetFrameB := BuildFrame(DevAddrI, SetFrameCmd, SetdatetimeB, DownlinkCounter)

	fmt.Printf("frame:%X \n", SetFrameB)

	//If a data frame carries a payload, FRMPayload must be encrypted before the message integrity code (MIC) is calculated.
	//Pmsg = MHDR | MACPayload cmac = aes128Cmac(NwkSKey, B0 | Pmsg) page 45
	// TxtSetDatetime := SecureFrame(SetFrameB, NwkSKeyB, AppSKeyB, DevAddrI, DownlinkCounter)
	TxtSetDatetime := SecureFrame(SetFrameB, NwkSKeyB, AppSKeyB, DevAddrI, DownlinkCounter)
	log.Print("-------------")

	LoraPktPayloadB := LoraPktEnc(TxtSetDatetime, 0)
	JSONPktPayloadB := []byte(BuildTXJsonPkt(LoraPktPayloadB))
	//fmt.Printf("webload json:%d  %s \n", len(JSONPktPayloadB), JSONPktPayloadB)

	//PULL_RESP packet set datetime
	PullrespPayload := make([]byte, 4+len(JSONPktPayloadB))
	PullrespPayload[0] = 0x01
	PullrespPayload[3] = 0x03

	for i := 0; i < len(JSONPktPayloadB); i++ {
		PullrespPayload[4+i] = JSONPktPayloadB[i]
	}

	return PullrespPayload
}
