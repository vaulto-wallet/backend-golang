package main

// #cgo CFLAGS: -I../wallet-core/include
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWCoinType.h>
// #include <TrustWalletCore/TWHDWallet.h>
import "C"

import (
	"fmt"
)
import "unsafe"
import "encoding/hex"
import "math/rand"

// C.TWString -> Go string
func TWStringGoString(s unsafe.Pointer) string {
	return C.GoString(C.TWStringUTF8Bytes(s))
}

// Go string -> C.TWString
func TWStringCreateWithGoString(s string) unsafe.Pointer {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	str := C.TWStringCreateWithUTF8Bytes(cStr)
	return str
}

// C.TWData -> Go byte[]
func TWDataGoBytes(d unsafe.Pointer) []byte {
	cBytes := C.TWDataBytes(d)
	cSize := C.TWDataSize(d)
	return C.GoBytes(unsafe.Pointer(cBytes), C.int(cSize))
}

// Go byte[] -> C.TWData
func TWDataCreateWithGoBytes(d []byte) unsafe.Pointer {
	cBytes := C.CBytes(d)
	defer C.free(unsafe.Pointer(cBytes))
	data := C.TWDataCreateWithBytes((*C.uchar)(cBytes), C.ulong(len(d)))
	return data
}

type CWallet *C.struct_TWHDWallet

func CreateBIP44(gate string, wallet *C.struct_TWHDWallet, Account uint32, Change uint32, Address uint32) (key *C.struct_TWPrivateKey) {
	var plus C.uint = 0
	/*
		gatesWithEd25519 := []string{"tezos","nimiq", "stellar", "aion", "kin", "nano", "waves", "aeternity", "solana", "near", "algorand", "ton", "kusama", "polkadot", "cardano"}
		if ok, index := in_array(gate, gatesWithEd25519); ok {
			log.Println("using index from Ed25519", index)
			plus = 0x80000000
		}
		gateObj := c.GateMap[gate]*/
	return C.TWHDWalletGetKeyBIP44(wallet, C.TWCoinTypeBitcoin, 0x80000000+C.uint(Account), plus+C.uint(Change), plus+C.uint(Address))
}

func main() {
	fmt.Println("==> calling wallet core from go")
	str := TWStringCreateWithGoString("confirm bleak useless tail chalk destroy horn step bulb genuine attract split")
	str2 := TWStringCreateWithGoString("orange okay much equip pond cushion ask hover bar shove ceiling have")
	empty := TWStringCreateWithGoString("")
	//empty2 := TWStringCreateWithGoString("1")
	defer C.TWStringDelete(str)
	defer C.TWStringDelete(str2)
	defer C.TWStringDelete(empty)

	fmt.Println("<== mnemonic is valid: ", C.TWHDWalletIsValid(str))

	//wallet := C.TWHDWalletCreateWithMnemonic(str2, empty)
	wallet := C.TWHDWalletCreate(256, empty)

	defer C.TWHDWalletDelete(wallet)

	//wallet2 := C.TWHDWalletCreateWithMnemonic(str2, empty)
	//wallet := C.TWHDWalletCreate(256, empty)
	//defer C.TWHDWalletDelete(wallet2)

	walletSeed := C.TWHDWalletSeed(wallet)
	walletSeedHex := hex.EncodeToString(TWDataGoBytes(walletSeed))
	fmt.Println("wallet seed: ", walletSeedHex)

	walletMnemonic := C.TWHDWalletMnemonic(wallet)
	fmt.Println("wallet mnemonic  : ", TWStringGoString(walletMnemonic))

	hex2, _ := hex.DecodeString(walletSeedHex)
	hex3 := TWDataCreateWithGoBytes(hex2)
	hex4 := C.TWDataCreateWithHexString(TWStringCreateWithGoString(walletSeedHex))
	//hex5 := C.CBytes( hex2 )
	hex5 := make([]byte, 32)
	rand.Read(hex5)

	defer C.TWDataDelete(hex3)
	defer C.TWDataDelete(hex4)
	//defer C.free(unsafe.Pointer(hex5))

	wallet2 := C.TWHDWalletCreateWithData(TWDataCreateWithGoBytes(hex5), empty)
	//wallet2 := C.TWHDWallet( hex5 , empty)
	wallet2Mnemonic := C.TWHDWalletMnemonic(wallet2)
	fmt.Println("wallet mnemonic 2: ", TWStringGoString(wallet2Mnemonic))

	walletSeed2 := C.TWHDWalletSeed(wallet2)
	walletSeedHex2 := hex.EncodeToString(TWDataGoBytes(walletSeed2))
	fmt.Println("wallet seed: ", walletSeedHex2)

	key := C.TWHDWalletGetKeyForCoin(wallet, C.TWCoinTypeBitcoin)
	keyData := C.TWPrivateKeyData(key)
	keyHex := hex.EncodeToString(TWDataGoBytes(keyData))
	fmt.Println("<== private key: ", keyHex)

	cointype := (uint32)(C.TWCoinTypeEthereum)

	key2 := C.TWHDWalletGetKeyBIP44(wallet2, cointype, 0x80000000, 0x80000000, 0x00000000)
	key2Data := C.TWPrivateKeyData(key2)
	key2Hex := hex.EncodeToString(TWDataGoBytes(key2Data))
	fmt.Println("<== private key: ", key2Hex)

	//path := "m/44'/1729'/0'/0'"

	key3 := C.TWHDWalletGetKeyBIP44(wallet2, cointype, 0x80000000, 0x80000000, 0x80000000)
	key3Data := C.TWPrivateKeyData(key3)
	key3Hex := hex.EncodeToString(TWDataGoBytes(key3Data))
	fmt.Println("<== private key: ", key3Hex)

	pubKey, _ := hex.DecodeString("0288be7586c41a0498c1f931a0aaf08c15811ee2651a5fe0fa213167dcaba59ae8")
	pubKeyData := TWDataCreateWithGoBytes(pubKey)
	defer C.TWDataDelete(pubKeyData)

	fmt.Println("<== bitcoin public key is valid: ", C.TWPublicKeyIsValid(pubKeyData, C.TWPublicKeyTypeSECP256k1))

	address := C.TWHDWalletGetAddressForCoin(wallet, C.TWCoinTypeBitcoin)
	fmt.Println("<== address: ", TWStringGoString(address))
	address2 := C.TWCoinTypeDeriveAddress(cointype, key2)
	fmt.Println("<== address2: ", TWStringGoString(address2))
	address3 := C.TWCoinTypeDeriveAddress(cointype, key3)
	fmt.Println("<== address3: ", TWStringGoString(address3))

	CreateBIP44("bitcoin", wallet, 0, 0, 0)

}
