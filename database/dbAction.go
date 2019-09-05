/* main function            */
/* file name: db_action.go  */
/* 		link:               */
/* 		                    */
/*    update: 20181122      */

package database

import (
	//"fmt"

	"encoding/base64"
	"fmt"
	"log"
	"loranet20181205/exception"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	//"encoding/hex"
	//"net"
	//"Strings"
	//"math/rand"
	//"encoding/base64"
)

func DbUpdate(IDI int, MeterIDs string, NwkSKeyB64 string, AppSKeyB64 string, AppNonceI int, DevNonce32 int, DevEUIs string, MeterUniqueID string, AK string, GUK string) int64 {

	var edInfo []TbEdInfo

	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()

	//verify connection
	// err = db.Ping()
	// if err != nil {
	// 	panic(err.Error())
	// }

	//fmt.Printf("M:%s N:%s D:%s\n",MeterID_s, NwkSKeyB64, DevEUI_s)

	//update
	// stmt, err := db.Prepare("update tb_ed_info set meterUniqueID=?, AK=?, GUK=?, MeterID=?, CustomerId=?, AppEUI=?, NwkSKey=?, AppSKey=?, AppNonce=?, DevNonce=?, NetID=?, activeStatus=? where DevEUI=?;")
	// result, err := stmt.Exec(meterUniqueID, AK, GUK, MeterID_s, MeterID_s, MeterID_s, NwkSKeyB64, AppSKeyB64, AppNonceI, DevNonce_32, 21, 1, DevEUI_s)
	//
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer stmt.Close() // Close the statement when we leave this function the program terminates

	rowsAffect := db.Model(&edInfo).Where("DevEUI = ?", DevEUIs).Updates(map[string]interface{}{"meterUniqueID": IDI, "AK": AK, "GUK": GUK, "MeterID": MeterIDs, "CustomerId": MeterIDs, "AppEUI": MeterIDs, "NwkSKey": NwkSKeyB64, "AppSKey": AppSKeyB64, "AppNonce": AppNonceI, "DevNonce": DevNonce32, "NetID": 21, "activeStatus": 1}).RowsAffected
	//get updated row number
	// rowsAffect, err := result.RowsAffected()
	// if err != nil {
	// 	panic(err.Error())
	// }
	return rowsAffect
	//fmt.Println(rowsAffect)

}

func DbParameters(devAddrI int) ([]byte, []byte, int, int) {
	// var NwKSKey string
	// var AppSKey string
	var downlinkCnt int
	var uplinkCnt int
	var invokeID int
	// var uplink_cnt int
	var edInfo []TbEdInfo
	var syscnt []TbSyscnt
	// Create Connection
	db, err := DbConnect()
	exception.CheckError(err)
	// Close Connection
	defer db.Close()
	// Query
	// qstr3 := "SELECT NwKSKey, AppSKey FROM tb_ed_info WHERE devAddr = " + strconv.Itoa(devAddrI)
	db.Where("devAddr = ?", strconv.Itoa(devAddrI)).Find(&edInfo)
	//fmt.Printf("qstr3: %s \n", qstr3)
	if len(edInfo) == 0 {
		log.Print("ERROR: EdInfo")
	}
	// log.Print(len(edInfo[0].NwkSKey))
	nwkSKeyB, err := base64.StdEncoding.DecodeString(edInfo[0].NwkSKey)
	exception.CheckError(err)
	appSKeyB, err := base64.StdEncoding.DecodeString(edInfo[0].AppSKey)
	exception.CheckError(err)

	//build frame payload
	// Query
	db.Where("devAddr = ?", strconv.Itoa(devAddrI)).Find(&syscnt)

	if len(syscnt) == 0 {
		downlinkCnt = 1
		invokeID = 0
		db.Create(&TbSyscnt{DevAddr: devAddrI, DownlinkCnt: downlinkCnt, UplinkCnt: uplinkCnt, InvokeID: invokeID})
		db.Last(&syscnt)
		log.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, invokeID)

		// log.Print("ERROR: Syscnt")
		log.Printf("insertid:%d\n", syscnt)
	} else {
		downlinkCnt = syscnt[0].DownlinkCnt
		invokeID = syscnt[0].InvokeID
		downlinkCnt++
		invokeID++
		rowsAffect := db.Model(&syscnt).Where("devAddr = ?", devAddrI).Updates(map[string]interface{}{"downlink_cnt": downlinkCnt, "Invoke_Id": invokeID}).RowsAffected
		log.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, invokeID)
		log.Printf("update:%d\n", rowsAffect)

	}
	// rows4, err := db.Query(qstr4)
	// exception.CheckError(err)

	// for rows4.Next() {
	// 	err = rows4.Scan(&downlink_cnt, &Invoke_Id)
	// 	exception.CheckError(err)
	// 	log.Printf("Dnlnk: %d Ivk:%d\n", downlinkCnt, invokeID)
	// 	DownlinkCnt++
	// 	InvokeID++
	// 	rowCount++
	// }
	//check counter table

	return nwkSKeyB, appSKeyB, downlinkCnt, invokeID
}

func DbGetDataRate(devAddrI int) int {

	var edConf []TbEdConf
	// Create Connection
	db, err := DbConnect()
	exception.CheckError(err)
	// Close Connection
	defer db.Close()
	// Query
	db.Where("devAddr = ?", strconv.Itoa(devAddrI)).Find(&edConf)

	if len(edConf) == 0 {
		log.Print("rocord not found")
		return 0
	} else {
		log.Printf("rate: %d fc:%d\n", edConf[0].DataRate, edConf[0].Uplinkfc)
	}
	return edConf[0].DataRate
}

func DbGetServerIP() (string, int) {

	var serverConfig []TbServerConfig
	// Create Connection
	db, err := DbConnect()
	exception.CheckError(err)
	// Close Connection
	defer db.Close()
	//get local IP
	ip, err := ExternalIP()
	if err != nil {
		log.Println(err)
	}
	// Query
	db.Model(&serverConfig).Where("id = ?", 1).Updates(map[string]interface{}{"ServerIP": ip})
	db.Where("id = ?", 1).Find(&serverConfig)
	if len(serverConfig) == 0 {
		log.Print("ERROR: ServerConfig")
	}
	return serverConfig[0].ServerIP, serverConfig[0].ServerPort
}

func DbGetGatewayIP() string {

	var gwInfo []TbGwInfo
	// Create Connection
	db, err := DbConnect()
	exception.CheckError(err)
	// Close Connection
	defer db.Close()
	// Query
	db.Where("gateway_id = ?", 1).Find(&gwInfo)
	if len(gwInfo) == 0 {
		log.Print("ERROR: GwInfo")
	}

	log.Printf("Gateway IP %s : %d\n", gwInfo[0].IPAddr, gwInfo[0].SocketPort)
	DefaultGateway := fmt.Sprintf("%s:%d", gwInfo[0].IPAddr, gwInfo[0].SocketPort)
	return DefaultGateway
}

func DbCheckEdInfo(devAddrI int) bool {

	var edInfo []TbEdInfo

	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Where("devAddr = ?", strconv.Itoa(devAddrI)).Find(&edInfo)
	if len(edInfo) == 0 {
		log.Print("ERROR: EdInfo:", devAddrI)
		return false
	} else {
		if len(edInfo[0].AppSKey) == 24 && len(edInfo[0].NwkSKey) == 24 {
			// log.Print("OK")
			return true
		} else {
			return false
		}
	}
}

func ActiveEdInfo(activeStatus int) []TbEdInfo {

	var edInfo []TbEdInfo

	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Where("activeStatus = ?", activeStatus).Find(&edInfo)

	return edInfo
}

func SetAcStatus2(activeStatus int) {
	var edInfo []TbEdInfo
	db, err := DbConnect()
	exception.CheckError(err)
	defer db.Close()
	db.Model(&edInfo).Where("activeStatus = ?", 1).Update("activeStatus", activeStatus)
}
