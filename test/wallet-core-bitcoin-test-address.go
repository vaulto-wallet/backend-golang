package main

// #cgo CFLAGS: -I../wallet-core/include
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWCoinType.h>
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWBitcoinSigHashType.h>
// #include <TrustWalletCore/TWBitcoinScript.h>
// #include <TrustWalletCore/TWAnySigner.h>
import "C"

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/big"
)
import "unsafe"
import "encoding/hex"
import "math/rand"
import "../proto/Bitcoin"

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

func BytesReversed(a *big.Int) (ret []byte) {
	retbuf := a.Bytes()
	ret = make([]byte, len(retbuf))
	for i, a := range retbuf {
		ret[len(ret)-1-i] = a
	}
	return ret
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

	key := C.TWHDWalletGetKeyForCoin(wallet, C.TWCoinTypeBitcoinTest)
	keyData := C.TWPrivateKeyData(key)
	keyHex := hex.EncodeToString(TWDataGoBytes(keyData))
	fmt.Println("<== private key: ", keyHex)

	cointype := (uint32)(C.TWCoinTypeBitcoinTest)

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

	// address := C.TWHDWalletGetAddressForCoin(wallet, C.TWCoinTypeBitcoinTest)
	//fmt.Println("<== address: ", TWStringGoString(address))
	address1 := C.TWCoinTypeDeriveAddress((uint32)(C.TWCoinTypeBitcoinTest), key)
	fmt.Println("<== address btct: ", TWStringGoString(address1))

	address2 := C.TWCoinTypeDeriveAddress(cointype, key2)
	fmt.Println("<== address btct: ", TWStringGoString(address2))
	address3 := C.TWCoinTypeDeriveAddress((uint32)(C.TWCoinTypeBitcoin), key2)
	fmt.Println("<== address btc: ", TWStringGoString(address3))
	address4 := C.TWCoinTypeDeriveAddress((uint32)(C.TWCoinTypeLitecoin), key2)
	fmt.Println("<== address ltc: ", TWStringGoString(address4))

	CreateBIP44("bitcoin", wallet, 0, 0, 0)

	// tb1qwmlqfd2cp0mupn942u9xk7as8hsfmtjumu8gkl c3355792f9d43ec52b0967330c1119f790fd8c817bffd63444269c286a635bae
	// tb1q30ry2eqcwjh3hx2gqsgn2katk0t7fu246rla4x 567701342537848c39f0c26014f39ee18b890626ecf7fcf664d72ffc6c0ca012
	privKeyInt, _ := new(big.Int).SetString("c3355792f9d43ec52b0967330c1119f790fd8c817bffd63444269c286a635bae", 16)
	//privkey,_ := hex.DecodeString("c3355792f9d43ec52b0967330c1119f790fd8c817bffd63444269c286a635bae")
	//privkey := BytesReversed(privKeyInt)
	privkey := privKeyInt.Bytes()

	script := C.TWBitcoinScriptBuildForAddress(TWStringCreateWithGoString("tb1qwmlqfd2cp0mupn942u9xk7as8hsfmtjumu8gkl"), cointype)
	//script := C.TWBitcoinScriptBuildPayToWitnessScriptHash(TWDataCreateWithGoBytes(privkey))
	//isType := C.TWBitcoinScriptIsPayToScriptHash(script)
	isType := C.TWBitcoinScriptIsPayToWitnessScriptHash(script)
	fmt.Println("<== script : ", TWDataGoBytes(C.TWBitcoinScriptData(script)), isType)
	fmt.Println("<== scriptHash : ", TWDataGoBytes(C.TWBitcoinScriptScriptHash(script)))
	defer C.TWBitcoinScriptDelete(script)

	//tx :=  new(TW_Bitcoin_Proto.UnspentTransaction)
	prevTxHash0, _ := new(big.Int).SetString("bbfd3792e7112a2041de88d59c16bfe025c6d53da6b6b3d300da38a77c3a01d3", 16)
	fmt.Println("<== txHash : ", prevTxHash0.Bytes())

	input0 := TW_Bitcoin_Proto.UnspentTransaction{
		OutPoint: &TW_Bitcoin_Proto.OutPoint{
			Hash:     BytesReversed(prevTxHash0),
			Index:    0,
			Sequence: uint32(0x0),
		},
		Amount: 1588868,
		Script: TWDataGoBytes(C.TWBitcoinScriptData(script)),
	}
	/*
		input := TW_Bitcoin_Proto.SigningInput{
			HashType: 1,
			Amount: 0x2222,
			ByteFee: 10,
			ToAddress: "tb1q30ry2eqcwjh3hx2gqsgn2katk0t7fu246rla4x",
			ChangeAddress: "tb1qwmlqfd2cp0mupn942u9xk7as8hsfmtjumu8gkl",
			PrivateKey: [][]byte{},
			Utxo: []*TW_Bitcoin_Proto.UnspentTransaction{},
			CoinType: uint32(C.TWCoinTypeBitcoinTest),
			//Scripts: map[string][]byte{},
		}*/

	input := TW_Bitcoin_Proto.SigningInput{
		HashType:   1,
		Amount:     0x2222,
		ByteFee:    10,
		PrivateKey: [][]byte{},
		Utxo:       []*TW_Bitcoin_Proto.UnspentTransaction{},
		CoinType:   uint32(C.TWCoinTypeBitcoinTest),
		Scripts:    map[string][]byte{},
	}

	input.Scripts[hex.EncodeToString(TWDataGoBytes(C.TWBitcoinScriptScriptHash(script)))] = TWDataGoBytes(C.TWBitcoinScriptData(script))
	input.PrivateKey = append(input.PrivateKey, privkey)
	input.Utxo = append(input.Utxo, &input0)

	out, err := proto.Marshal(&input)
	fmt.Println((string)(out), err)
	outData := TWDataCreateWithGoBytes(out)
	goData := TWDataGoBytes(outData)
	fmt.Println((string)(goData), err)

	btct_plan := C.TWAnySignerPlan(outData, C.TWCoinTypeBitcoin)
	plan := new(TW_Bitcoin_Proto.TransactionPlan)
	proto.Unmarshal(TWDataGoBytes(unsafe.Pointer(btct_plan)), plan)

	input.Plan = plan

	fmt.Println("Plan amount", plan.Amount)
	fmt.Println("Plan Utxos", plan.Utxos)

	btct_out := C.TWAnySignerSign(TWDataCreateWithGoBytes(out), C.TWCoinTypeBitcoinTest)
	output := new(TW_Bitcoin_Proto.SigningOutput)
	proto.Unmarshal(TWDataGoBytes(unsafe.Pointer(btct_out)), output)
	fmt.Println(output.String())
	fmt.Println(output.GetEncoded())
	fmt.Println(hex.EncodeToString(output.GetEncoded()))
	/*
		prevTxHash0,_ := hex.DecodeString("fff7f7881a8099afa6940d42d1e7f6362bec38171ea3edf433541db4e4ad969f")
		script, _ := hex.DecodeString("2103c9f4836b9a4f77fc0d81f7bcb01b7f1b35916864b9476c241ce9fc198bd25432ac")
		input0 := TW_Bitcoin_Proto.UnspentTransaction{
			OutPoint:       &TW_Bitcoin_Proto.OutPoint{
				Hash:                 prevTxHash0,
				Index:                1,
				Sequence:             uint32(2),
			},
			Amount: 0x77000000,
			Script: script,
		}
		privkey,_ := hex.DecodeString("bbc27228ddcb9209d7fd6f36b02f7dfa6252af40bb2f1cbc7a557da8027ff866")

		input := TW_Bitcoin_Proto.SigningInput{
			HashType: 1,
			Amount: 0x2222,
			ByteFee: 1,
			ToAddress: "1Bp9U1ogV3A14FMvKbRJms7ctyso4Z4Tcx",
			ChangeAddress: "1FQc5LdgGHMHEN9nwkjmz6tWkxhPpxBvBU",
			PrivateKey: [][]byte{},
			Utxo: []*TW_Bitcoin_Proto.UnspentTransaction{&input0},
			CoinType: uint32(C.TWCoinTypeBitcoin),
			//Scripts: map[string][]byte{},
		}

		//input.Utxo = append(input.Utxo, input0)
		ib,_ := proto.Marshal(&input)
		print(ib)
		//input.Utxo[0] = input0

		input.PrivateKey = append(input.PrivateKey, privkey)

		//input.ProtoMessage()

		out, err := proto.Marshal(&input)
		fmt.Print(out)

		input2 := TW_Bitcoin_Proto.SigningInput{}
		proto.Unmarshal(out, &input2)
		fmt.Print(input2)

		fmt.Println((string)(out), err)
		outData := TWDataCreateWithGoBytes(out)
		goData := TWDataGoBytes(outData)
		fmt.Println((string)(goData), err)

		btct_plan := C.TWAnySignerPlan(outData, C.TWCoinTypeBitcoin)
		plan := new(TW_Bitcoin_Proto.TransactionPlan)
		proto.Unmarshal(TWDataGoBytes(unsafe.Pointer(btct_plan)), plan)

		input.Plan = plan

		fmt.Println("Plan amount", plan.Amount)
		fmt.Println( "Plan Utxos", plan.Utxos )

		btct_out := C.TWAnySignerSign(TWDataCreateWithGoBytes(out), C.TWCoinTypeBitcoin)
		output := new(TW_Bitcoin_Proto.SigningOutput)
		proto.Unmarshal(TWDataGoBytes(unsafe.Pointer(btct_out)), output)
		fmt.Println(output.String())
		fmt.Println(output.GetEncoded())
		fmt.Println(hex.EncodeToString(output.GetEncoded()))
	*/

}
