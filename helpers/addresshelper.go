package helpers

// #cgo CFLAGS: -I../wallet-core/include
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWCoinType.h>
import "C"
import "encoding/hex"

func GenerateAddress(asset string, seed string, change uint32, n uint32) ([]byte, string) {
	coins := map[string]uint32{
		"BTC":  C.TWCoinTypeBitcoin,
		"BTCT": C.TWCoinTypeBitcoinTest,
		"ETH":  C.TWCoinTypeEthereum,
	}
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return nil, ""
	}

	empty := TWStringCreateWithGoString("")
	defer C.TWStringDelete(empty)
	seedPointer := TWDataCreateWithGoBytes(seedBytes)
	defer C.TWDataDelete(seedPointer)
	wallet := C.TWHDWalletCreateWithData(seedPointer, empty)
	defer C.TWHDWalletDelete(wallet)
	assetId, exists := coins[asset]
	if !exists {
		return nil, ""
	}

	key := C.TWHDWalletGetKeyBIP44(wallet, assetId, 0, (C.uint)(change), (C.uint)(n))
	keyData := C.TWPrivateKeyData(key)
	//defer C.TWDataDelete(keyData)
	address := C.TWCoinTypeDeriveAddress(assetId, key)
	//defer C.TWDataDelete(address)

	//stringKey := hex.EncodeToString(TWDataGoBytes(keyData))
	stringKey := TWDataGoBytes(keyData)
	stringAddress := TWStringGoString(address)
	return stringKey, stringAddress
}
