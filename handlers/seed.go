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
	h "../helpers"
	m "../models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
)

func CreateSeed(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	masterPassword := req.Context().Value("masterPassword")
	user := req.Context().Value("user").(*m.User)
	log.Println(user)
	if masterPassword == nil {
		ReturnErrorWithStatusString(w, Error(NotInitialized), 400, "Not initialized")
		return
	}

	var r interface{}
	err := json.NewDecoder(req.Body).Decode(&r)

	if err != nil {
		ReturnError(w, Error(BadRequest))
		return
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

	pub := m.ConfigRecord{Name: "PublicKey"}.Get(db)
	if pub.ID == 0 || len(pub.Value) == 0 {
		ReturnErrorWithStatusString(w, Error(NotInitialized), 400, "Master public key is not set")
		return
	}

	encryptedSeed, _ := h.EncryptWithRSA([]byte(pub.Value), seed)

	seedHex := hex.EncodeToString(encryptedSeed)

	seedName := "New seed"
	seedParam, exists := rm["name"]
	if exists {
		seedName = seedParam.(string)
	}

	dbSeed := m.Seed{
		Name:  seedName,
		Seed:  seedHex,
		Owner: user,
	}
	db.Create(&dbSeed)

	ReturnResult(w, dbSeed.ID)
}

func GetSeeds(db *gorm.DB, w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*m.User)

	var seeds m.Seeds

	db.Where("owner_id = ? ", user.ID).Find(&seeds)

	res, err := json.Marshal(seeds)
	if err != nil {
		ReturnResult(w, seeds)
		return
	}

	fmt.Println((string)(res))

	ReturnResult(w, seeds)
}
