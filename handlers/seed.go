package handlers

// #cgo CFLAGS: -I../wallet-core/include
// #cgo LDFLAGS: -L../wallet-core/build -L../wallet-core/build/trezor-crypto -lTrustWalletCore -lprotobuf -lTrezorCrypto -lc++ -lm
// #include <TrustWalletCore/TWHDWallet.h>
// #include <TrustWalletCore/TWString.h>
// #include <TrustWalletCore/TWData.h>
// #include <TrustWalletCore/TWPrivateKey.h>
// #include <TrustWalletCore/TWPublicKey.h>
// #include <TrustWalletCore/TWCoinType.h>
import "C"

import (
	m "../models"
	h "../trusthelpers"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
)

func CreateSeed(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	if dbUser.ID == 0 {
		ReturnError(w, Error(NoUser))
		return
	}
	var r interface{}
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
	}
	rm := r.(map[string]interface{})

	var seed []byte

	mnemonicParam, mnemonicExists := rm["mnemonic"]
	if mnemonicExists {
		mnemonic := h.TWStringCreateWithGoString(mnemonicParam.(string))
		empty := h.TWStringCreateWithGoString("")

		defer C.TWStringDelete(mnemonic)
		defer C.TWStringDelete(empty)

		wallet := C.TWHDWalletCreateWithMnemonic(mnemonic, empty)
		defer C.TWHDWalletDelete(wallet)

		walletSeed := C.TWHDWalletSeed(wallet)
		defer C.TWDataDelete(walletSeed)
		seed = h.TWDataGoBytes(walletSeed)[0:32]
	} else {
		seed = make([]byte, 32)
		rand.Read(seed)
	}

	seedHex := hex.EncodeToString(seed)

	seedName := "New seed"
	seedParam, exists := rm["name"]
	if exists {
		seedName = seedParam.(string)
	}

	dbSeed := m.Seed{
		Name:  seedName,
		Seed:  seedHex,
		Owner: dbUser,
	}
	db.Create(&dbSeed)

	ReturnResult(w, dbSeed.ID)
}

func GetSeeds(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	username := req.Context().Value("user")
	dbUser := m.User{}
	db.First(&dbUser, "Username = ?", username)

	var seeds m.Seeds

	db.Where("owner_id = ? ", dbUser.ID).Find(&seeds)

	res, err := json.Marshal(seeds)
	if err != nil {
		ReturnResult(w, seeds)
		return
	}

	fmt.Println((string)(res))

	ReturnResult(w, seeds)
}
