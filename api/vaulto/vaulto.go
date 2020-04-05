package vaulto

import (
	m "../../models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type VaultoAPI struct {
	token  string
	url    string
	client *http.Client
}

func (a *VaultoAPI) Init(url string) {
	a.url = url
	a.client = &http.Client{}
}

func (a *VaultoAPI) Request(method string, endpoint string, data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	req, e := http.NewRequest(method, a.url+endpoint, buf)
	if e != nil {
		return []byte{}, e
	}
	if len(a.token) > 0 {
		req.Header.Add("Authorization", "Bearer "+a.token)
	}

	resp, e := a.client.Do(req)
	if e != nil {
		return []byte{}, e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return []byte{}, e
	}
	return body, nil
}

func (a *VaultoAPI) Clear() (bool, error) {
	resp, err := a.Request("GET", "/clear", nil)
	if err != nil {
		return false, err
	}
	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	return response.Result, nil
}

func (a *VaultoAPI) Login(username string, password string) (bool, error) {

	resp, err := a.Request("POST", "/users/login", LoginRequest{username, password})
	if err != nil {
		return false, err
	}

	log.Print(string(resp))

	var response ResponseString
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	a.token = response.Result
	return true, nil
}

func (a *VaultoAPI) Register(username string, password string) (int64, error) {
	resp, err := a.Request("POST", "/users/register", LoginRequest{username, password})
	if err != nil {
		return -1, err
	}
	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}

	return int64(response.Result.(float64)), nil
}

func (a *VaultoAPI) CreateAsset(name string, symbol string, index int, decimals int, rounding int) (int64, error) {
	resp, err := a.Request("POST", "/assets", m.Asset{
		Name:      name,
		CoinIndex: index,
		Symbol:    symbol,
		Type:      1,
		Decimals:  decimals,
		Rounding:  rounding,
	})
	if err != nil {
		return -1, err
	}

	log.Print(string(resp))

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil
}

func (a *VaultoAPI) GetAssets() (m.Assets, error) {
	resp, err := a.Request("GET", "/assets", nil)
	if err != nil {
		return nil, err
	}

	log.Print(string(resp))

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}
	var ar m.Assets
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&ar)
	return ar, nil
}

func (a *VaultoAPI) CreateSeed(name string, seed string) (int64, error) {
	resp, err := a.Request("POST", "/seeds", SeedRequest{name, ""})
	if err != nil {
		return -1, err
	}
	log.Print(string(resp))

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil
}

func (a *VaultoAPI) GetSeeds() (interface{}, error) {
	resp, err := a.Request("GET", "/seeds", nil)
	if err != nil {
		return false, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return response.Result, nil
}

func (a *VaultoAPI) CreateWallet(name string, seedId uint, assetId uint) (int64, error) {
	resp, err := a.Request("POST", "/wallets", m.Wallet{
		Name:    name,
		AssetId: assetId,
		SeedId:  seedId,
	})
	if err != nil {
		return -1, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil
}

func (a *VaultoAPI) GetWallets() (m.Wallets, error) {
	resp, err := a.Request("GET", "/wallets", nil)

	var response ResponseInterface

	if err != nil {
		return nil, err
	}

	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var wr m.Wallets
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&wr)
	return wr, nil
}

func (a *VaultoAPI) GetWalletsForAsset(asset string) (m.Wallets, error) {
	resp, err := a.Request("GET", "/wallets/"+asset, nil)

	if err != nil {
		return nil, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var wr m.Wallets
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&wr)
	return wr, nil
}

func (a *VaultoAPI) CreateAddress(name string, walletId int) (int64, error) {
	resp, err := a.Request("POST", "/address", m.Address{
		Name:     name,
		WalletID: walletId,
	})
	if err != nil {
		return -1, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil
}

func (a *VaultoAPI) GetAddressesForWallet(wallet uint) ([]m.Address, error) {
	resp, err := a.Request("GET", "/address/"+strconv.Itoa(int(wallet)), nil)

	//	var addresses []m.Address
	addresses := new([]m.Address)

	if err != nil {
		return *addresses, err
	}

	//var response ResponseInterface
	var response struct {
		Error     string
		ErrorText string
		Result    []m.Address
	}

	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return *addresses, errors.New(response.ErrorText)
	}

	//var wr m.Addresses
	/*
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(response.Result)
		for _, v := range response.Result.([]interface{}) {
			address := new(m.Address)
			mapstructure.Decode(v, address )
			*addresses = append(*addresses, *address)
		}
		return *addresses, nil*/
	return response.Result, nil
}

func (a *VaultoAPI) CreateOrder(asset string, address_to string, amount float64, comment string) (int64, error) {
	resp, err := a.Request("POST", "/orders", m.OrderData{
		Symbol:    asset,
		AddressTo: address_to,
		Amount:    amount,
		Comment:   comment,
	})
	if err != nil {
		return -1, err
	}
	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil

}

func (a *VaultoAPI) UpdateOrder(orderId uint, status m.OrderStatus) (bool, error) {
	resp, err := a.Request("PUT", "/orders", m.Order{
		Model:     gorm.Model{ID: orderId},
		OrderData: m.OrderData{Status: status},
	})
	if err != nil {
		return false, err
	}
	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return response.Result.(bool), nil
}

func (a *VaultoAPI) GetOrders() (m.Orders, error) {
	resp, err := a.Request("GET", "/orders", nil)

	if err != nil {
		return nil, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var orders m.Orders
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&orders)
	return orders, nil
}

func (a *VaultoAPI) CreateTransaction(assetId uint, walletId uint, orderId uint, addressIds []uint,
	txHash string, tx string, txData string) (int64, error) {
	resp, err := a.Request("POST", "/transactions", m.Transaction{
		Name:      "",
		AssetId:   assetId,
		WalletId:  walletId,
		OrderId:   orderId,
		AddressId: addressIds,
		TxHash:    txHash,
		Tx:        tx,
		TxData:    txData,
		Status:    m.TransactionStatusNew,
	})
	if err != nil {
		return -1, err
	}
	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return -1, errors.New(response.ErrorText)
	}
	return int64(response.Result.(float64)), nil

}

func (a *VaultoAPI) GetTransactions() ([]m.Transaction, error) {
	resp, err := a.Request("GET", "/transactions", nil)

	if err != nil {
		return nil, err
	}

	var response struct {
		ResponseInterface
		Result []m.Transaction
	}
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	return response.Result, nil
}
