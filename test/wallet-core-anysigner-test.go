package main

// #cgo CFLAGS: -I../wallet-core/include
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWEthereumProto.h>
// #include <TrustWalletCore/TWEthereumSigner.h>
// #include <TrustWalletCore/TWAnyProto.h>
// #include <TrustWalletCore/TWAnySigner.h>
// #include <TrustWalletCore/TWCurve.h>
import "C"

import (
	"encoding/base64"
	"fmt"
)
import "unsafe"
import "encoding/hex"
import "../proto/Any"

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

	//message := "AQ=="
	//message := "Mw=="
	message := base64.StdEncoding.EncodeToString([]byte{3})

	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	base64.StdEncoding.Decode(base64Text, []byte(message))
	fmt.Println(message, base64Text)

	tx := "{" +
		"\"chainId\":\"" + "Mw==" +
		"\",\"gasPrice\":\"" + "1pOkAA==" +
		"\",\"gasLimit\":\"" + "Ugg=" +
		"\",\"toAddress\":\"" + "0xC37054b3b48C3317082E7ba872d7753D13da4986" +
		"\",\"amount\":\"" + "A0i8paFgAA==" + "\"}"

	is_valid_priv := C.TWPrivateKeyIsValid(keyData, C.TWCurveSECP256k1)
	fmt.Println("Valid", is_valid_priv)

	input := TW_Any_Proto.SigningInput{}
	input.CoinType = C.TWCoinTypeEthereum
	input.PrivateKey = keyHex
	input.Transaction = tx

	fmt.Println(input.String())
	fmt.Println(tx)

	input_c := TWDataCreateWithGoBytes(([]byte)(input.String()))

	fmt.Println(([]byte)(input.String()))
	fmt.Println(input.GetPrivateKey())
	fmt.Println(input_c)
	out, err := input.XXX_Marshal(nil, true)
	fmt.Println((string)(out), err)

	ethout := C.TWAnySignerSign((C.TW_Any_Proto_SigningInput)(TWDataCreateWithGoBytes(out)))

	fmt.Println((string)(TWDataGoBytes(unsafe.Pointer(ethout))))

	output := TW_Any_Proto.SigningOutput{}
	output.XXX_Unmarshal(TWDataGoBytes(unsafe.Pointer(ethout)))
	fmt.Println(output.String())

	fmt.Println(output.Output)

}

// f86a8084d693a40082520894c37054b3b48c3317082e7ba872d7753d13da4986870348bca5a160008026a09f9b20dcb30f7e4a3990d24a283dd1f31f7f3dfcebc9d952cd5699cf0c391e0aa06920439229e3c9ed3d1e8e3ff6dee5f8005691171ef4f960d0b696943c599171
// f86b8084d693a40082520894c37054b3b48c3317082e7ba872d7753d13da4986870348bca5a1600080818aa0417a00a54bc2ea383802a8102b981e901589cf606de95facce55c463107db4c3a020509e38245f6908898ac03d598bf4b9ca61f683bab2b1aca69334984cfd8735
