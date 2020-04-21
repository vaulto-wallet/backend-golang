package builder

// #cgo CFLAGS: -I../../../wallet-core/include -I../../../wallet-core/src -I/usr/local/include/boost
// #cgo LDFLAGS: -L../../../wallet-core/build -L../../../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
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
	h "../../../trusthelpers"
	"log"
	"math/big"
	"unsafe"
)
import "encoding/hex"
import "../../../proto/Ethereum"

func BuildEthereum(private_key []byte, to string, value big.Int, gasLimit big.Int, gasPrice big.Int, nonce big.Int, payload []byte) (tx string) {
	//wallet := C.TWHDWalletCreateWithData(h.TWDataCreateWithGoBytes(keyData), "")

	//keyData := C.TWPrivateKeyCreateWithData(h.TWDataCreateWithGoBytes(private_key))
	privateKeyHex := make([]byte, 32)
	hex.Decode(privateKeyHex, private_key)

	input := new(TW_Ethereum_Proto.SigningInput)
	input.Amount = value.Bytes()
	input.PrivateKey = privateKeyHex
	input.ChainId = new(big.Int).SetInt64(3).Bytes()
	input.GasLimit = gasLimit.Bytes()
	input.GasPrice = gasPrice.Bytes()
	input.Nonce = nonce.Bytes()
	input.ToAddress = to
	input.Payload = payload

	out, _ := input.XXX_Marshal(nil, true)

	//input_c := h.TWDataCreateWithGoBytes(([]byte)(input.String()))
	ethout := C.TWAnySignerSign(h.TWDataCreateWithGoBytes(out), C.TWCoinTypeEthereum)

	output := new(TW_Ethereum_Proto.SigningOutput)
	output.XXX_Unmarshal(h.TWDataGoBytes(unsafe.Pointer(ethout)))
	log.Println("buildEtehreum tx :", hex.EncodeToString(output.Encoded))
	return hex.EncodeToString(output.Encoded)
}

func BuildERC20Transfer(to string, value big.Int) []byte {
	toHex := to
	if toHex[0:2] == "0x" {
		toHex = toHex[2:]
	}
	addressInt, _ := new(big.Int).SetString(toHex, 16)

	abiFx := C.TWEthereumAbiEncoderBuildFunction(h.TWStringCreateWithGoString("transfer"))
	valueData := h.TWDataCreateWithGoBytes(value.Bytes())
	toData := h.TWDataCreateWithGoBytes(addressInt.Bytes())
	C.TWEthereumAbiFunctionAddParamAddress(abiFx, toData, false)
	C.TWEthereumAbiFunctionAddParamUInt256(abiFx, valueData, false)

	encoded := C.TWEthereumAbiEncoderEncode(abiFx)

	return h.TWDataGoBytes(encoded)

}
