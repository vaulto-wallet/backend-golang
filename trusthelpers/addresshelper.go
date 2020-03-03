package trusthelpers

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

func GenerateAddress(asset string, seed string, change uint32, n uint32) (string, string) {
	coins := map[string]uint32{
		"BTC": C.TWCoinTypeBitcoin,
		"ETH": C.TWCoinTypeEthereum,
	}
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return "", ""
	}

	empty := TWStringCreateWithGoString("")
	defer C.TWStringDelete(empty)
	wallet := C.TWHDWalletCreateWithData(TWDataCreateWithGoBytes(seedBytes), empty)
	defer C.TWHDWalletDelete(wallet)
	assetId, exists := coins[asset]
	if !exists {
		return "", ""
	}

	key := C.TWHDWalletGetKeyBIP44(wallet, assetId, 0, (C.uint)(change), (C.uint)(n))
	keyData := C.TWPrivateKeyData(key)
	//defer C.TWDataDelete(keyData)
	address := C.TWCoinTypeDeriveAddress(assetId, key)
	//defer C.TWDataDelete(address)
	stringKey := hex.EncodeToString(TWDataGoBytes(keyData))
	stringAddress := TWStringGoString(address)
	return stringKey, stringAddress
}
