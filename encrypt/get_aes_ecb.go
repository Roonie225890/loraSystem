package aesacc

import (
	"crypto/aes"
	"crypto/cmac"
	//"crypto/cmac"
	"crypto/cipher"
	//"encoding/hex"
	//"fmt"
)

func GetAesECB(txt []byte, key []byte) []byte {
	//txt, _ := hex.DecodeString("4800F21500B300000204414000000000")
	//key, _ := hex.DecodeString("A7D8942966B2B7BF6109829AEE3EEAA9")
	//enc, _ := hex.DecodeString("7EF38A5E3B2CCB587DCDC48117C62FD2")
	//dest_d := DecryptAes128Ecb(txt, key)
	dest_e := make([]byte, (len(txt)/len(key)+1)*len(key))

	aesCipher, _ := aes.NewCipher([]byte(key))
	encrypter := cipher.NewECBEncrypter(aesCipher)
	//decrypter := cipher.NewECBDecrypter(aesCipher)

	encrypter.CryptBlocks(dest_e, []byte(txt))

	return dest_e
}

func GetAesECBDec(data []byte, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func GetCmac(txt []byte, key []byte) []byte {
	//txt, _ := hex.DecodeString("00303734393438303648415202000000B302BB00000000")
	//key, _ := hex.DecodeString("A7D8942966B2B7BF6109829AEE3EEAA9")
	//enc, _ := hex.DecodeString("7EF38A5E3B2CCB587DCDC48117C62FD2")

	block, _ := aes.NewCipher(key)

	h, _ := cmac.Sum(txt, block, 16)

	//j := 0
	//for  i:=len(txt)-4 ;i<len(txt);i++ {
	//	txt[i] = h[j]
	//	j ++
	//}

	return h[0:4]

	//fmt.Printf("plain: %s\n" , hex.EncodeToString([]byte(txt)))
	//fmt.Printf("key:   %s\n" , hex.EncodeToString([]byte(key)))
	//fmt.Printf("dest:  %x\n" , h)

}
