package main

import "C"
import "unsafe"

// #cgo CFLAGS: -I../wallet-core/include -I../wallet-core/src -I/usr/local/include/boost
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWAnySigner.h>
// #include <TrustWalletCore/TWCurve.h>
// #include <TrustWalletCore/TWCoinType.h>
// #include <TrustWalletCore/TWEthereumAbiEncoder.h>
// #include <TrustWalletCore/TWEthereumAbiFunction.h>
import "C"

import (
	"fmt"
	"math/big"
)
import "encoding/hex"
import "../proto/Ethereum"

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

func main() {
	fmt.Println("==> calling wallet core from go")
	str := TWStringCreateWithGoString("confirm bleak useless tail chalk destroy horn step bulb genuine attract split")
	empty := TWStringCreateWithGoString("")
	defer C.TWStringDelete(str)
	defer C.TWStringDelete(empty)

	fmt.Println("<== mnemonic is valid: ", C.TWHDWalletIsValid(str))

	wallet := C.TWHDWalletCreateWithMnemonic(str, empty)

	defer C.TWHDWalletDelete(wallet)

	walletSeed := C.TWHDWalletSeed(wallet)
	walletSeedHex := hex.EncodeToString(TWDataGoBytes(walletSeed))
	fmt.Println("wallet seed: ", walletSeedHex)

	walletMnemonic := C.TWHDWalletMnemonic(wallet)
	fmt.Println("wallet mnemonic  : ", TWStringGoString(walletMnemonic))

	key := C.TWHDWalletGetKeyForCoin(wallet, C.TWCoinTypeEthereum)
	keyData := C.TWPrivateKeyData(key)
	keyHex := hex.EncodeToString(TWDataGoBytes(keyData))
	//keyHex := TWStringGoString( C.TWStringCreateWithHexData(keyData) )

	fmt.Println("<== ethereum private key: ", keyHex)

	address := C.TWHDWalletGetAddressForCoin(wallet, C.TWCoinTypeEthereum)
	fmt.Println("<== ethereum address: ", TWStringGoString(address))

	abiFx := C.TWEthereumAbiEncoderBuildFunction(TWStringCreateWithGoString("transfer"))
	addressInt, _ := new(big.Int).SetString("C37054b3b48C3317082E7ba872d7753D13da4986", 16)
	valueInt, _ := new(big.Int).SetString("10000000000", 10)

	to := TWDataCreateWithGoBytes(addressInt.Bytes())
	value := TWDataCreateWithGoBytes(valueInt.Bytes())
	p1 := C.TWEthereumAbiFunctionAddParamAddress(abiFx, to, false)
	p2 := C.TWEthereumAbiFunctionAddParamUInt256(abiFx, value, false)

	encoded := C.TWEthereumAbiEncoderEncode(abiFx)

	fmt.Println("ABI", p1, p2, TWDataGoBytes(encoded))

	is_valid_priv := C.TWPrivateKeyIsValid(keyData, C.TWCurveSECP256k1)
	fmt.Println("Valid", is_valid_priv)

	input := TW_Ethereum_Proto.SigningInput{}
	input.Amount = new(big.Int).SetInt64(96000000000).Bytes()
	input.PrivateKey = TWDataGoBytes(keyData)
	input.ChainId = new(big.Int).SetInt64(1).Bytes()
	input.GasLimit = new(big.Int).SetInt64(21000).Bytes()
	input.GasPrice = new(big.Int).SetInt64(4).Bytes()
	input.Nonce = new(big.Int).SetInt64(1).Bytes()
	input.ToAddress = "0xC37054b3b48C3317082E7ba872d7753D13da4986"
	input.Payload = TWDataGoBytes(encoded)

	fmt.Println(input.String())

	input_c := TWDataCreateWithGoBytes(([]byte)(input.String()))

	fmt.Println(([]byte)(input.String()))
	fmt.Println(input.GetPrivateKey())
	fmt.Println(input_c)
	out, err := input.XXX_Marshal(nil, true)
	fmt.Println((string)(out), err)

	ethout := C.TWAnySignerSign(TWDataCreateWithGoBytes(out), C.TWCoinTypeEthereum)

	fmt.Println((string)(TWDataGoBytes(unsafe.Pointer(ethout))))

	output := TW_Ethereum_Proto.SigningOutput{}
	output.XXX_Unmarshal(TWDataGoBytes(unsafe.Pointer(ethout)))
	fmt.Println(output.String())

	fmt.Println(hex.EncodeToString(output.Encoded))

}
