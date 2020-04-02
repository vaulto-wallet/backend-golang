package vaulto

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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

func (a *VaultoAPI) Register(username string, password string) (bool, error) {
	resp, err := a.Request("POST", "/users/register", LoginRequest{username, password})
	if err != nil {
		return false, err
	}
	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}

	return response.Result, nil
}

func (a *VaultoAPI) CreateAsset(name string, symbol string, index int, decimals int, rounding int) (bool, error) {
	resp, err := a.Request("POST", "/assets", AssetRequest{
		Name:      name,
		CoinIndex: index,
		Symbol:    symbol,
		Type:      1,
		Decimals:  decimals,
		Rounding:  rounding,
	})
	if err != nil {
		return false, err
	}

	log.Print(string(resp))

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
}

func (a *VaultoAPI) GetAssets() (AssetsResponse, error) {
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
	var ar AssetsResponse
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&ar)
	return ar, nil
}

func (a *VaultoAPI) CreateSeed(name string, seed string) (bool, error) {
	resp, err := a.Request("POST", "/seeds", SeedRequest{name, ""})
	if err != nil {
		return false, err
	}
	log.Print(string(resp))

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
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

func (a *VaultoAPI) CreateWallet(name string, seedId uint, assetId uint) (bool, error) {
	resp, err := a.Request("POST", "/wallets", WalletRequest{
		Name:    name,
		AssetId: assetId,
		SeedId:  seedId,
	})
	if err != nil {
		return false, err
	}

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
}

func (a *VaultoAPI) GetWallets() (WalletsResponse, error) {
	resp, err := a.Request("GET", "/wallets", nil)

	var response ResponseInterface

	if err != nil {
		return nil, err
	}

	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var wr WalletsResponse
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&wr)
	return wr, nil
}

func (a *VaultoAPI) GetWalletsForAsset(asset string) (WalletsResponse, error) {
	resp, err := a.Request("GET", "/wallets/"+asset, nil)

	if err != nil {
		return nil, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var wr WalletsResponse
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&wr)
	return wr, nil
}

func (a *VaultoAPI) CreateAddress(name string, walletId int) (bool, error) {
	resp, err := a.Request("POST", "/address", AddressRequest{
		Name:     name,
		WalletID: walletId,
	})
	if err != nil {
		return false, err
	}

	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil
}

func (a *VaultoAPI) CreateOrder(asset string, address_to string, amount float64, comment string) (bool, error) {
	resp, err := a.Request("POST", "/orders", OrderRequest{
		Symbol:    asset,
		AddressTo: address_to,
		Amount:    amount,
		Comment:   comment,
	})
	if err != nil {
		return false, err
	}
	var response ResponseBool
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return false, errors.New(response.ErrorText)
	}
	return true, nil

}

func (a *VaultoAPI) GetOrders() (OrdersResponse, error) {
	resp, err := a.Request("GET", "/orders", nil)

	if err != nil {
		return nil, err
	}

	var response ResponseInterface
	json.NewDecoder(strings.NewReader(string(resp))).Decode(&response)

	if len(response.Error) > 0 {
		return nil, errors.New(response.ErrorText)
	}

	var orders OrdersResponse
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(response.Result)
	json.NewDecoder(buf).Decode(&orders)
	return orders, nil
}
